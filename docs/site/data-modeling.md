# Data Modeling

Astrate uses PostgreSQL 16 + TimescaleDB for all metadata, properties, and time-series data. This replaces Astarte's Cassandra/ScyllaDB dependency.

## Tenancy layout

**Decision: shared tables + `realm_id`.** Row-level tenancy with composite keys is simpler and faster at the 1-5 realm scale Astrate targets. Realm deletion is a transactional cascade.

## Relational metadata

### Realms

```sql
CREATE TABLE realms (
    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',
    ca_certificate   text NOT NULL,
    ca_private_key   bytea NOT NULL,  -- AES-256-GCM encrypted
    device_registration_limit integer,
    created_at       timestamptz NOT NULL DEFAULT now()
);
```

### Interfaces

Interface definitions stored as JSONB with generated columns for routing-critical fields:

```sql
CREATE TABLE interfaces (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    definition    jsonb NOT NULL,
    name          text     GENERATED ALWAYS AS (definition->>'interface_name') STORED,
    major_version integer  GENERATED ALWAYS AS ((definition->>'version_major')::int) STORED,
    minor_version integer  GENERATED ALWAYS AS ((definition->>'version_minor')::int) STORED,
    type          text     GENERATED ALWAYS AS (definition->>'type') STORED,
    ownership     text     GENERATED ALWAYS AS (definition->>'ownership') STORED,
    aggregation   text     GENERATED ALWAYS AS (coalesce(definition->>'aggregation','individual')) STORED,
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (realm_id, name, major_version)
);
```

### Endpoints

Normalized endpoint mappings with typed columns for each mapping attribute:

```sql
CREATE TABLE endpoints (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    interface_id  bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint      text NOT NULL,
    value_type    text NOT NULL,
    reliability   text NOT NULL DEFAULT 'unreliable',
    retention     text NOT NULL DEFAULT 'discard',
    expiry        integer NOT NULL DEFAULT 0,
    database_retention_policy text NOT NULL DEFAULT 'no_ttl',
    database_retention_ttl    integer,
    explicit_timestamp boolean NOT NULL DEFAULT false,
    allow_unset   boolean NOT NULL DEFAULT false,
    UNIQUE (interface_id, endpoint)
);
```

### Devices

```sql
CREATE TABLE devices (
    id                  uuid NOT NULL,
    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    credentials_secret_hash text NOT NULL,
    status              text NOT NULL DEFAULT 'registered',
    introspection       jsonb NOT NULL DEFAULT '{}',
    old_introspection   jsonb NOT NULL DEFAULT '{}',
    aliases             jsonb NOT NULL DEFAULT '{}',
    attributes          jsonb NOT NULL DEFAULT '{}',
    cert_serial         text,
    cert_aki            text,
    first_registration  timestamptz NOT NULL DEFAULT now(),
    first_credentials_request timestamptz,
    last_credentials_request_ip inet,
    last_connection     timestamptz,
    last_disconnection  timestamptz,
    last_seen_ip        inet,
    connected           boolean NOT NULL DEFAULT false,
    total_received_msgs  bigint NOT NULL DEFAULT 0,
    total_received_bytes bigint NOT NULL DEFAULT 0,
    payload_format_hint  text NOT NULL DEFAULT 'bson',
    PRIMARY KEY (realm_id, id)
);
```

Device status values: `registered` (awaiting first credentials request), `confirmed` (has connected at least once), `inhibited` (blocked from new credentials and connections).

## Properties storage

Properties are last-value-wins key/value state -- a plain relational table with upsert semantics.

```sql
CREATE TABLE properties (
    realm_id     smallint NOT NULL,
    device_id    uuid NOT NULL,
    interface_id bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint_id  bigint NOT NULL REFERENCES endpoints(id),
    path         text NOT NULL,
    value        jsonb NOT NULL,
    value_type   text NOT NULL,
    set_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (realm_id, device_id, interface_id, path),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
);
```

- **Device-owned properties:** device publishes, server stores; empty payload = unset (`allow_unset`).
- **Server-owned properties:** AppEngine writes, retained in broker for device delivery.
- **Purge:** `producer/properties` control message lists device-owned paths; server deletes any not in the list.

## Datastream storage (TimescaleDB hypertables)

### Individual datastreams

One wide hypertable with sparse typed columns -- the same shape Astarte uses on Cassandra:

```sql
CREATE TABLE individual_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    endpoint_id   bigint NOT NULL,
    path          text NOT NULL,
    ts            timestamptz NOT NULL,
    reception_ts  timestamptz NOT NULL,
    value_double       double precision,
    value_integer      integer,
    value_longinteger  bigint,
    value_boolean      boolean,
    value_string       text,
    value_binaryblob   bytea,
    value_datetime     timestamptz,
    value_array        jsonb
);
```

Hypertable with 7-day chunks, compressed after 7 days (segment by `realm_id, device_id, interface_id, path`).

### Object datastreams

```sql
CREATE TABLE object_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    path          text NOT NULL,
    ts            timestamptz NOT NULL,
    reception_ts  timestamptz NOT NULL,
    value         jsonb NOT NULL
);
```

Same hypertable/compression shape as individual datastreams.

## Retention and expiry

- **`database_retention_policy: use_ttl`** (per endpoint): a background job runs chunk-batched `DELETE` for aged rows.
- **Global hard-cap retention** (`retention.max_days`): uses Timescale `add_retention_policy` for cheap whole-chunk eviction.
- **Datastream `expiry`** (message validity for offline devices): honoured on the server-to-device path via MQTT message expiry.
- **Downsample queries:** map onto Timescale `time_bucket()` + optionally `lttb()` from the toolkit.
