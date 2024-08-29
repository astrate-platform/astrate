import Config

config :logger, :console,
  format: {PrettyLog.LogfmtFormatter, :format},
  metadata: [:realm, :datacenter, :replication_factor, :module, :function, :tag]

# Configure your database
config :astarte_dev_tool, Astarte.DataAccess.Repo,
  # Overrides the default type `bigserial` used for version attribute in schema migration
  migration_primary_key: [name: :id, type: :binary_id],
  contact_points: [System.get_env("CASSANDRA_DB_HOST") || "localhost"],
  keyspace: "astarte",
  port: System.get_env("CASSANDRA_DB_PORT") || 9042,
  # Waiting time in milliseconds for the database connection
  sync_connect: 5000,
  log: :info,
  stacktrace: true,
  show_sensitive_data_on_connection_error: true,
  pool_size: 10

config :astarte_dev_tool, ecto_repos: [Astarte.DataAccess.Repo]
