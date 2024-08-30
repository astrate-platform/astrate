#
# This file is part of Astarte.
#
# Copyright 2024 SECO Mind Srl
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

defmodule AstarteDevTool.Utilities.Auth do
  alias Astarte.Client.Credentials
  alias X509.PrivateKey
  alias X509.PublicKey

  def gen_auth_token(private_key) when is_bitstring(private_key) do
    Credentials.dashboard_credentials() |> Credentials.to_jwt(private_key)
  end

  def pem_key(), do: {:ok, PrivateKey.new_ec(:secp256r1) |> PrivateKey.to_pem()}

  def pem_key(priv) do
    case PrivateKey.from_pem(priv) do
      {:ok, result} -> {:ok, result |> PublicKey.derive() |> PublicKey.to_pem()}
      error -> error
    end
  end

  def pem_keys() do
    with {:ok, priv} <- pem_key(),
         {:ok, pub} <- pem_key(priv) do
      {:ok, {priv, pub}}
    end
  end
end
