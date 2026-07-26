-- 000006: trigger delivery policies (upstream 1.1 Realm Management surface,
-- M10 dashboard compatibility). Policies are stored and served over the API;
-- the trigger executor does not consult them yet — see the
-- TODO(policies-semantics) block in internal/engine/triggers/actions.go and
-- docs/COMPATIBILITY.md for the sanctioned follow-up scope.

CREATE TABLE trigger_policies (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name       text NOT NULL,
    definition jsonb NOT NULL,        -- upstream delivery-policy JSON (error_handlers, retry_times, ...)
    UNIQUE (realm_id, name)
);
