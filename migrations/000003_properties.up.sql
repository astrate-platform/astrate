-- 000003: properties storage, transcribed verbatim from docs/DESIGN.md §2.3
-- (docs/ROADMAP.md §3.1 file 2.3).
CREATE TABLE properties (
    realm_id     smallint NOT NULL,
    device_id    uuid NOT NULL,
    interface_id bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint_id  bigint NOT NULL REFERENCES endpoints(id),
    path         text NOT NULL,               -- concrete path, e.g. '/lcdCmd' or '/4/enable'
    value        jsonb NOT NULL,
    value_type   text NOT NULL,
    set_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (realm_id, device_id, interface_id, path),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
);
