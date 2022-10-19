#
# This file is part of Astarte.
#
# Copyright 2017-2018 Ispirata Srl
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

defmodule Astarte.TriggerEngine.AMQPEventsConsumer do
  require Logger
  use GenServer

  alias AMQP.Basic
  alias AMQP.Channel
  alias AMQP.Connection
  alias AMQP.Exchange
  alias AMQP.Queue
  alias Astarte.TriggerEngine.Config
  alias Astarte.TriggerEngine.Policy
  alias Astarte.TriggerEngine.Policy.PolicySupervisor

  @connection_backoff 10000
  # API

  def start_link(args \\ []) do
    GenServer.start_link(__MODULE__, args, name: __MODULE__)
  end

  def ack(delivery_tag) do
    GenServer.call(__MODULE__, {:ack, delivery_tag})
  end

  # Server callbacks

  def init(_args) do
    send(self(), :try_to_connect)
    {:ok, %{channel: nil}}
  end

  def terminate(_reason, state) do
    if state.channel do
      conn = state.channel.conn
      Channel.close(state.channel)
      Connection.close(conn)
    end
  end

  def handle_call({:ack, _delivery_tag}, _from, %{channel: nil} = state) do
    {:reply, {:error, :disconnected}, state}
  end

  def handle_call({:ack, delivery_tag}, _from, %{channel: chan} = state) do
    res = Basic.ack(chan, delivery_tag)
    {:reply, res, state}
  end

  # Confirmation sent by the broker after registering this process as a consumer
  def handle_info({:basic_consume_ok, %{consumer_tag: _consumer_tag}}, state) do
    {:noreply, state}
  end

  # Sent by the broker when the consumer is unexpectedly cancelled (such as after a queue deletion)
  def handle_info({:basic_cancel, %{consumer_tag: _consumer_tag}}, state) do
    {:noreply, state}
  end

  # Confirmation sent by the broker to the consumer process after a Basic.cancel
  def handle_info({:basic_cancel_ok, %{consumer_tag: _consumer_tag}}, state) do
    {:noreply, state}
  end

  # Message consumed
  def handle_info({:basic_deliver, payload, meta}, state) do
    {headers, other_meta} = Map.pop(meta, :headers, [])
    headers_map = amqp_headers_to_map(headers)

    Logger.debug(
      "got event, payload: #{inspect(payload)}, headers: #{inspect(headers_map)}, meta: #{
        inspect(other_meta)
      }"
    )

    {realm_name, policy_name} = get_headers_map_trigger_info(headers_map)
    policy_process = get_policy_process(realm_name, policy_name)
    Policy.handle_event(policy_process, payload, meta, state.channel)

    {:noreply, state}
  end

  def handle_info(:try_to_connect, _state) do
    {:ok, new_state} = connect()
    {:noreply, new_state}
  end

  def handle_info({:DOWN, _, :process, _pid, reason}, _state) do
    Logger.warn("RabbitMQ connection lost: #{inspect(reason)}. Trying to reconnect...")
    {:ok, new_state} = connect()
    {:noreply, new_state}
  end

  defp connect() do
    with amqp_consumer_options = Config.amqp_consumer_options!(),
         {:ok, conn} <- Connection.open(amqp_consumer_options),
         {:ok, chan} <- Channel.open(conn),
         :ok <- Basic.qos(chan, prefetch_count: Config.amqp_consumer_prefetch_count!()),
         events_exchange_name = Config.events_exchange_name!(),
         events_queue_name = Config.events_queue_name!(),
         events_routing_key = Config.events_routing_key!(),
         :ok <- Exchange.declare(chan, events_exchange_name, :direct, durable: true),
         {:ok, _queue} <- Queue.declare(chan, events_queue_name, durable: true),
         :ok <-
           Queue.bind(
             chan,
             events_queue_name,
             events_exchange_name,
             routing_key: events_routing_key
           ),
         {:ok, _consumer_tag} <- Basic.consume(chan, events_queue_name),
         # Get notifications when the chan or conn go down
         Process.monitor(chan.pid) do
      # TODO add policy to state
      {:ok, %{channel: chan}}
    else
      {:error, reason} ->
        Logger.warn("RabbitMQ Connection error: #{inspect(reason)}")
        retry_after(@connection_backoff)
        {:ok, %{channel: nil}}

      _ ->
        Logger.warn("Unknown RabbitMQ connection error")
        retry_after(@connection_backoff)
        {:ok, %{channel: nil}}
    end
  end

  defp retry_after(backoff) when is_integer(backoff) do
    Logger.warn("Retrying connection in #{backoff} ms")
    Process.send_after(self(), :try_to_connect, backoff)
  end

  defp amqp_headers_to_map(headers) do
    Enum.reduce(headers, %{}, fn {key, _type, value}, acc ->
      Map.put(acc, key, value)
    end)
  end

  defp get_headers_map_trigger_info(headers_map) do
    with {:ok, realm} <- Map.fetch(headers_map, "x_astarte_realm"),
         {:ok, policy} <- Map.fetch(headers_map, "x_astarte_trigger_policy") do
      {realm, policy}
    end
  end

  defp get_policy_process(realm, policy_name) do
    case Registry.lookup(Registry.PolicyRegistry, {realm, policy_name}) do
      [] ->
        child = {Policy, [realm_name: realm, policy_name: policy_name]}
        {:ok, pid} = PolicySupervisor.start_child(child)
        pid

      [{pid, nil}] ->
        pid
    end
  end
end
