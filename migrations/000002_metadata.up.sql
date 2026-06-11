-- 000002: relational metadata schema, transcribed verbatim from
-- docs/DESIGN.md §2.2 (docs/ROADMAP.md §3.1 file 2.2).

-- Realms (Housekeeping domain)
CREATE TABLE realms (
    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
    ca_private_key   bytea NOT NULL,                 -- encrypted at rest (AES-256-GCM,
                                                     -- key from config/KMS env var)
    device_registration_limit integer,
    created_at       timestamptz NOT NULL DEFAULT now()
);

-- Interfaces (Realm Management domain). The raw JSON is the source of truth;
-- generated columns lift the routing-critical fields out for indexing.
CREATE TABLE interfaces (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    definition    jsonb NOT NULL,
    name          text     GENERATED ALWAYS AS (definition->>'interface_name') STORED,
    major_version integer  GENERATED ALWAYS AS ((definition->>'version_major')::int) STORED,
    minor_version integer  GENERATED ALWAYS AS ((definition->>'version_minor')::int) STORED,
    type          text     GENERATED ALWAYS AS (definition->>'type') STORED,          -- datastream|properties
    ownership     text     GENERATED ALWAYS AS (definition->>'ownership') STORED,     -- device|server
    aggregation   text     GENERATED ALWAYS AS (coalesce(definition->>'aggregation','individual')) STORED,
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (realm_id, name, major_version)
);

-- Mappings, normalized for endpoint-id stability (mirrors Astarte's endpoint UUIDs).
CREATE TABLE endpoints (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    interface_id  bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint      text NOT NULL,            -- e.g. '/%{sensor_id}/value'
    value_type    text NOT NULL,            -- double|integer|boolean|longinteger|string|
                                            -- binaryblob|datetime|<type>array
    reliability   text NOT NULL DEFAULT 'unreliable',  -- → QoS 0|1|2
    retention     text NOT NULL DEFAULT 'discard',
    expiry        integer NOT NULL DEFAULT 0,
    database_retention_policy text NOT NULL DEFAULT 'no_ttl',
    database_retention_ttl    integer,
    explicit_timestamp boolean NOT NULL DEFAULT false,
    allow_unset   boolean NOT NULL DEFAULT false,
    UNIQUE (interface_id, endpoint)
);

-- Devices (Pairing + AppEngine domains)
CREATE TABLE devices (
    id                  uuid NOT NULL,          -- the 128-bit Astarte device ID
    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    credentials_secret_hash text NOT NULL,      -- bcrypt
    status              text NOT NULL DEFAULT 'registered',  -- registered|confirmed|inhibited
    introspection       jsonb NOT NULL DEFAULT '{}',  -- {"iface.Name": {"major":1,"minor":2}, ...}
    old_introspection   jsonb NOT NULL DEFAULT '{}',
    aliases             jsonb NOT NULL DEFAULT '{}',
    attributes          jsonb NOT NULL DEFAULT '{}',
    cert_serial         text,                   -- serial of last issued client cert
    cert_aki            text,                   -- authority key identifier
    first_registration  timestamptz NOT NULL DEFAULT now(),
    first_credentials_request timestamptz,
    last_credentials_request_ip inet,
    last_connection     timestamptz,
    last_disconnection  timestamptz,
    last_seen_ip        inet,
    connected           boolean NOT NULL DEFAULT false,
    total_received_msgs  bigint NOT NULL DEFAULT 0,
    total_received_bytes bigint NOT NULL DEFAULT 0,
    payload_format_hint  text NOT NULL DEFAULT 'bson',  -- bson|json, see §3.5.4
    PRIMARY KEY (realm_id, id)
);
CREATE INDEX devices_aliases_gin ON devices USING gin (aliases);

-- Device groups (AppEngine)
CREATE TABLE groups (
    id        bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name      text NOT NULL,
    UNIQUE (realm_id, name)
);
CREATE TABLE group_devices (
    group_id  bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    realm_id  smallint NOT NULL,
    device_id uuid NOT NULL,
    PRIMARY KEY (group_id, device_id),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
);

-- Triggers (Realm Management domain; executed by the engine)
CREATE TABLE triggers (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name       text NOT NULL,
    definition jsonb NOT NULL,        -- Astarte trigger JSON (simple_triggers + action)
    UNIQUE (realm_id, name)
);
