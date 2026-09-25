slug: docs-sync-rm-delete-device-async-operation-param
verdict: blocked
at:  9f20120
ran: 2026-09-25T20:45:21Z on DietPi in 177s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/realm/async_operation_test.go

$ sed -n '95,130p' /root/astrate-mule/internal/realm/http.go
if err := a.svc.DeletePolicy(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteDevice synchronously removes a device and its data (the dashboard's
// Delete button; docs/COMPATIBILITY.md notes the deviation from upstream's
// async deletion).
func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
	err := a.svc.DeleteDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if errors.Is(err, store.ErrNotFound) {
		_ = astarteapi.WriteDeviceNotFound(w)
		return
	}
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// getVersion reports the emulated upstream API level (the dashboard gates
// feature UI on it — see APICompatVersion).
func (a *API) getVersion(w http.ResponseWriter, _ *http.Request) {
	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
}

func (a *API) getRegistrationLimit(w http.ResponseWriter, r *http.Request) {
	limit, err := a.svc.GetDeviceRegistrationLimit(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, limit)

$ cd /root/astrate-mule && sed -n '650,730p' docs/api/astarte_realm_management_api.yaml
operationId: getVersion
      summary: Get emulated API version
      description: |
        Returns the emulated upstream API compatibility version. The Astarte
        Dashboard gates feature UI on this value.
      tags: [Realm Management - Misc]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: API compatibility version.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: string
                    description: Version string.
              example:
                data: "1.1.0"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Device lifecycle ────────────────────────────────────────────────

  /realmmanagement/v1/{realm}/devices/{device}:
    delete:
      operationId: deleteDevice
      summary: Delete a device
      description: |
        Synchronously deletes a device and all its data. This is a deliberate
        deviation from upstream's async deletion (see COMPATIBILITY.md).
        Returns 204 on success.
      tags: [Realm Management - Devices]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/AsyncOperation"
      responses:
        "204":
          description: Device deleted. No body returned.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Trigger delivery policies ───────────────────────────────────────

  /realmmanagement/v1/{realm}/policies:
    get:
      operationId: listPolicies
      summary: List trigger delivery policies
      description: Returns the names of all trigger delivery policies in the realm.
      tags: [Realm Management - Policies]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: A list of policy names.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeStringArray"
              example:
                data:

$ cd /root/astrate-mule && git status --porcelain && echo "---LOG---" && git log --oneline -5
---LOG---
9f20120 mule: log docs-sync-rm-delete-device-async-operation-param
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7d78bd3 mule: log docs-sync-rm-async-operation-param
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
17bfa74 mule: log docs-sync-hk-retention-zero-is-unset

$ cd /root/astrate-mule && git log --oneline -3 -- docs/api/astarte_realm_management_api.yaml internal/realm/async_operation_test.go docs/COMPATIBILITY.md
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
9c05370 mule: docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && git log -S'AsyncOperation' --oneline -- docs/api/astarte_realm_management_api.yaml
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && git show --stat f965f39 && echo "=== 3f0eca8 ===" && git show --stat 3f0eca8
commit f965f39a02b555c6c11e21a94294d98126fe354b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 25 22:19:10 2026 +0200

    mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

 .mule/todo.md                              |  1 +
 docs/COMPATIBILITY.md                      | 17 ++++++++--------
 docs/api/astarte_realm_management_api.yaml | 12 +++++++-----
 internal/realm/async_operation_test.go     | 31 +++++++++++++++++++++++++-----
 internal/swagger/swagger_test.go           | 19 ++++++++++--------
 5 files changed, 54 insertions(+), 26 deletions(-)
=== 3f0eca8 ===
commit 3f0eca89014a4a6d993f076621b2dba499628d58
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 25 22:00:01 2026 +0200

    mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

 .mule/todo.md                              |  1 +
 docs/api/astarte_realm_management_api.yaml | 19 +++++++++++++
 internal/swagger/swagger_test.go           | 43 ++++++++++++++++++++++++++++++
 3 files changed, 63 insertions(+)

$ cd /root/astrate-mule && git show f965f39 -- .mule/todo.md && echo "=== todo grep ===" && grep -n 'delete-device-async' .mule/todo.md
commit f965f39a02b555c6c11e21a94294d98126fe354b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 25 22:19:10 2026 +0200

    mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

diff --git a/.mule/todo.md b/.mule/todo.md
index 871c64e..aaa88fd 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -226,4 +226,5 @@ worth having. Before any of it reaches `main`:
 - [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
 - [ ] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== todo grep ===
229:- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
230:- [ ] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && git show f965f39 -- docs/COMPATIBILITY.md docs/api/astarte_realm_management_api.yaml internal/swagger/swagger_test.go internal/realm/async_operation_test.go
commit f965f39a02b555c6c11e21a94294d98126fe354b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 25 22:19:10 2026 +0200

    mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

diff --git a/docs/COMPATIBILITY.md b/docs/COMPATIBILITY.md
index e34dd45..879ee9b 100644
--- a/docs/COMPATIBILITY.md
+++ b/docs/COMPATIBILITY.md
@@ -334,14 +334,15 @@ All additive or strictly-safer; none affect unmodified device SDKs.
 
 17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
     upstream 1.4 runs realm create/delete (Housekeeping), interface
-    install/update/delete (Realm Management) and trigger-delivery-policy
-    delete in the background, answering immediately and letting the caller opt
-    into synchronous execution with `?async_operation=false`. Astrate performs
-    every one of these synchronously and answers only once the work is done —
-    a strictly stronger guarantee, so a client written against upstream never
-    observes a regression. The parameter is accepted and ignored on either
-    value, on every one of those endpoints, so upstream clients that send it
-    keep working unchanged; pinned by
+    install/update/delete and device deletion (Realm Management) and
+    trigger-delivery-policy delete in the background, answering immediately
+    and letting the caller opt into synchronous execution with
+    `?async_operation=false`. Astrate performs every one of these
+    synchronously and answers only once the work is done — a strictly
+    stronger guarantee, so a client written against upstream never observes a
+    regression. The parameter is accepted and ignored on either value, on
+    every one of those endpoints, so upstream clients that send it keep
+    working unchanged; pinned by
     `TestRealmManagementAsyncOperationParam` and
     `TestHousekeepingAsyncOperationParam`. Tracked in
     `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 5a229f6..d6886f4 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -694,6 +694,7 @@ paths:
       parameters:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/DeviceID"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Device deleted. No body returned.
@@ -905,11 +906,12 @@ components:
       required: false
       description: |
         Accepted and ignored on either value. Upstream 1.4 runs interface
-        install/update/delete and trigger-delivery-policy delete in the
-        background and lets the caller opt into synchronous execution with
-        `async_operation=false`; Astrate always works synchronously and answers
-        only once the work is done, so the value changes nothing and the
-        response is identical either way (deviation 17 in COMPATIBILITY.md).
+        install/update/delete, device deletion and trigger-delivery-policy
+        delete in the background and lets the caller opt into synchronous
+        execution with `async_operation=false`; Astrate always works
+        synchronously and answers only once the work is done, so the value
+        changes nothing and the response is identical either way
+        (deviation 17 in COMPATIBILITY.md).
       schema:
         type: boolean
         default: false
diff --git a/internal/realm/async_operation_test.go b/internal/realm/async_operation_test.go
index 1560203..a957bc6 100644
--- a/internal/realm/async_operation_test.go
+++ b/internal/realm/async_operation_test.go
@@ -3,19 +3,23 @@
 package realm
 
 import (
+	"context"
 	"fmt"
 	"net/http"
 	"testing"
+
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 
 // TestRealmManagementAsyncOperationParam pins deviation 17
-// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete
-// and policy delete asynchronously unless the caller opts out with
-// `?async_operation=false`. Astrate is always synchronous, so the parameter is
-// accepted and ignored on either value — sending it must never draw a 4xx, and
-// the effect must be visible on the very next read.
+// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete,
+// device deletion and policy delete asynchronously unless the caller opts out
+// with `?async_operation=false`. Astrate is always synchronous, so the
+// parameter is accepted and ignored on either value — sending it must never
+// draw a 4xx, and the effect must be visible on the very next read.
 func TestRealmManagementAsyncOperationParam(t *testing.T) {
 	r := newRig(t)
+	ctx := context.Background()
 
 	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {
 		name := fmt.Sprintf("com.ex.M7a.Async%d", i)
@@ -64,5 +68,22 @@ func TestRealmManagementAsyncOperationParam(t *testing.T) {
 		if rec := r.req(t, http.MethodGet, "/policies/"+policyName, "", r.rmaToken); rec.Code != http.StatusNotFound {
 			t.Errorf("get policy after delete %s: got %d, want 404", q, rec.Code)
 		}
+
+		// Device delete carries the parameter upstream too: the dashboard's
+		// Delete button, which upstream answers before the data is gone.
+		dev, err := deviceid.Random()
+		if err != nil {
+			t.Fatalf("random device id: %v", err)
+		}
+		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
+			t.Fatalf("register device %s: %v", dev, err)
+		}
+		if rec := r.req(t, http.MethodDelete, "/devices/"+dev.String()+q, "", r.rmaToken); rec.Code != http.StatusNoContent {
+			t.Fatalf("delete device %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
+		}
+		// Synchronous means gone by the time the 204 is written, on both values.
+		if _, err := r.st.GetDevice(ctx, r.realmID, dev); err == nil {
+			t.Errorf("device row survived the delete %s", q)
+		}
 	}
 }
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index a2d357a..8f1e81b 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -211,13 +211,13 @@ func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
 
 // TestRealmManagementAsyncOperationParamDocumented guards that the realm
 // management spec tells clients the `?async_operation` parameter exists on the
-// four operations upstream 1.4 runs in the background: interface
-// install/update/delete and trigger-delivery-policy delete. Astrate accepts and
-// ignores it there (deviation 17, docs/COMPATIBILITY.md) — none of the four
-// handlers reads the query string, which is the whole of the acceptance — so an
-// upstream client keeps working, but a client generated from the spec has to
-// learn the parameter from the spec. This is the documentation half of the
-// behaviour pinned by TestRealmManagementAsyncOperationParam in
+// five operations upstream 1.4 runs in the background: interface
+// install/update/delete, device deletion and trigger-delivery-policy delete.
+// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
+// none of the five handlers reads the query string, which is the whole of the
+// acceptance — so an upstream client keeps working, but a client generated from
+// the spec has to learn the parameter from the spec. This is the documentation
+// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
 // internal/realm.
 func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
 	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
@@ -227,7 +227,10 @@ func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
 	lines := strings.Split(string(b), "\n")
 
 	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
-	for _, op := range []string{"installInterface", "updateInterface", "deleteInterface", "deletePolicy"} {
+	for _, op := range []string{
+		"installInterface", "updateInterface", "deleteInterface",
+		"deleteDevice", "deletePolicy",
+	} {
 		block := operationBlock(t, lines, op)
 		if !containsLine(block, ref) {
 			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)

$ cd /root/astrate-mule && cat docs/Makefile 2>/dev/null | head -40
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

The work for this task appears to already exist in the tree. Let me verify carefully before deciding.
$ cd /root/astrate-mule && grep -n 'AsyncOperation' -A 20 docs/api/astarte_realm_management_api.yaml | sed -n '/parameters:/,$p' | head -5; echo "=== component ==="; sed -n '900,925p' docs/api/astarte_realm_management_api.yaml
=== component ===
      schema:
        type: string

    AsyncOperation:
      name: async_operation
      in: query
      required: false
      description: |
        Accepted and ignored on either value. Upstream 1.4 runs interface
        install/update/delete, device deletion and trigger-delivery-policy
        delete in the background and lets the caller opt into synchronous
        execution with `async_operation=false`; Astrate always works
        synchronously and answers only once the work is done, so the value
        changes nothing and the response is identical either way
        (deviation 17 in COMPATIBILITY.md).
      schema:
        type: boolean
        default: false

  schemas:
    # ── Config types ──

    AuthConfig:
      type: object
      required: [jwt_public_key_pem]
      properties:

All three pieces of the task are already in the tree. Now let me run the gate to confirm the state is actually green.
$ cd /root/astrate-mule && go test ./internal/swagger/ -run 'AsyncOperation' -v 2>&1 | tail -20
=== RUN   TestHousekeepingAsyncOperationParamDocumented
--- PASS: TestHousekeepingAsyncOperationParamDocumented (0.00s)
=== RUN   TestRealmManagementAsyncOperationParamDocumented
--- PASS: TestRealmManagementAsyncOperationParamDocumented (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.059s

$ cd /root/astrate-mule && gofmt -l . 2>&1; echo "--- gofmt exit: $? ---"; go vet ./... 2>&1 | tail -5; echo "--- vet exit: $? ---"
--- gofmt exit: 0 ---
--- vet exit: 0 ---

$ cd /root/astrate-mule && sed -n '205,250p' internal/swagger/swagger_test.go
}
	}
	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
	}
}

// TestRealmManagementAsyncOperationParamDocumented guards that the realm
// management spec tells clients the `?async_operation` parameter exists on the
// five operations upstream 1.4 runs in the background: interface
// install/update/delete, device deletion and trigger-delivery-policy delete.
// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
// none of the five handlers reads the query string, which is the whole of the
// acceptance — so an upstream client keeps working, but a client generated from
// the spec has to learn the parameter from the spec. This is the documentation
// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
// internal/realm.
func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
	for _, op := range []string{
		"installInterface", "updateInterface", "deleteInterface",
		"deleteDevice", "deletePolicy",
	} {
		block := operationBlock(t, lines, op)
		if !containsLine(block, ref) {
			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
		}
	}

	comp := componentBlock(t, lines, "    AsyncOperation:")
	for _, want := range []string{
		"      name: async_operation",
		"      in: query",
		"      required: false",
		"      schema:",
		"        type: boolean",
		"        default: false",
	} {
		if !containsLine(comp, want) {
			t.Errorf("components parameter AsyncOperation is missing line %q", want)

$ cd /root/astrate-mule && go vet -tags integration ./internal/realm/ 2>&1 | tail -5; echo "--- integration vet exit: $? ---"
--- integration vet exit: 0 ---

$ cd /root/astrate-mule && (command -v mkdocs && echo MKDOCS_OK) ; (python3 -c 'import yaml; print("pyyaml ok")' 2>&1 | tail -1)
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && go test ./... 2>&1 | grep -v '^ok\|no test files' | head -20; echo "--- full test done ---"
--- full test done ---

$ cd /root/astrate-mule && ls internal/swagger/ && grep -rn 'func Test' internal/swagger/*_test.go | head -30
swagger.go
swagger_test.go
18:func TestMount(t *testing.T) {
77:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
90:func TestSpecs(t *testing.T) {
123:func TestRealmManagement403(t *testing.T) {
179:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
222:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
297:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {

$ cd /root/astrate-mule && grep -rn 'yaml.Unmarshal\|yaml.v3\|kin-openapi\|openapi3' --include=*.go internal/ pkg/ 2>/dev/null | head -20
(no output)

$ cd /root/astrate-mule && sed -n '90,123p' internal/swagger/swagger_test.go
func TestSpecs(t *testing.T) {
	got := Specs()

	if len(got) == 0 {
		t.Fatal("Specs() returned no filenames")
	}

	for _, name := range got {
		if strings.HasPrefix(name, "api/") || strings.Contains(name, "/") {
			t.Errorf("Specs() entry %q should have no path prefix or dirs", name)
		}
		if !strings.HasSuffix(name, ".yaml") {
			t.Errorf("Specs() entry %q does not end in .yaml", name)
		}
		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
		}
	}

	want, err := embeddedYAMLFilenames()
	if err != nil {
		t.Fatalf("enumerating embedded basenames: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Specs() = %v, want %v", got, want)
	}
}

// TestRealmManagement403 guards the contract that every realm-management
// operation guarded by a_rma documents a 403 Forbidden response, not just
// 401: RequireRealm answers 403 for a verified realm JWT whose a_rma grants
// do not authorize the method+path (internal/auth/middleware.go). Matches the
// upstream-parity the housekeeping and pairing specs already document.
func TestRealmManagement403(t *testing.T) {

$ ls -d $(go env GOMODCACHE)/gopkg.in/yaml* 2>/dev/null; echo "---"; grep -c 'yaml' /root/astrate-mule/go.sum 2>/dev/null
/root/go/pkg/mod/gopkg.in/yaml.v2@v2.2.8
/root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1
---
5

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && cat > go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
cat > main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		fmt.Println("PARSE FAIL:", err)
		os.Exit(1)
	}
	paths := doc["paths"].(map[string]any)
	params := doc["components"].(map[string]any)["parameters"].(map[string]any)
	fmt.Println("parsed OK, paths:", len(paths), "component params:", len(params))
	if _, ok := params["AsyncOperation"]; !ok {
		fmt.Println("MISSING components/parameters/AsyncOperation")
		os.Exit(1)
	}
	var refs []string
	for p, v := range paths {
		for m, op := range v.(map[string]any) {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			for _, pv := range o["parameters"].([]any) {
				pm, _ := pv.(map[string]any)
				if r, ok := pm["$ref"].(string); ok && r == "#/components/parameters/AsyncOperation" {
					refs = append(refs, m+" "+p)
				}
			}
		}
	}
	sort.Strings(refs)
	fmt.Println("ops $ref-ing AsyncOperation:")
	for _, r := range refs {
		fmt.Println("  ", r)
	}
}
EOF
GOFLAGS=-mod=mod GOPROXY=off go mod tidy 2>&1 | tail -3; GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
parsed OK, paths: 12 component params: 7
ops $ref-ing AsyncOperation:
   delete /realmmanagement/v1/{realm}/devices/{device}
   delete /realmmanagement/v1/{realm}/interfaces/{name}/{major}
   delete /realmmanagement/v1/{realm}/policies/{name}
   post /realmmanagement/v1/{realm}/interfaces
   put /realmmanagement/v1/{realm}/interfaces/{name}/{major}

$ cd /root/astrate-mule && git check-ignore -v docs/site/DESIGN.md docs/site/api 2>&1; echo "--- ignore exit: $? ---"
.gitignore:32:/docs/site/DESIGN.md	docs/site/DESIGN.md
.gitignore:34:/docs/site/api/	docs/site/api
--- ignore exit: 0 ---

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -8; echo "--- exit: ${PIPESTATUS[0]} ---"; git status --porcelain
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
/bin/sh: 1: Bad substitution

$ cd /root/astrate-mule && git status --porcelain; echo "--- clean check done ---"; grep -rn 'async_operation' docs/*.md docs/api/*.yaml | grep -v 'name: async_operation' | head -30
--- clean check done ---
docs/COMPATIBILITY.md:340:    `?async_operation=false`. Astrate performs every one of these
docs/UPSTREAM-EXPERIMENTAL.md:25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
docs/api/astarte_housekeeping_api.yaml:246:        `async_operation=false`; Astrate always works synchronously and answers
docs/api/astarte_realm_management_api.yaml:911:        execution with `async_operation=false`; Astrate always works

$ cd /root/astrate-mule && sed -n '333,350p' docs/COMPATIBILITY.md
log-and-drop).

17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
    upstream 1.4 runs realm create/delete (Housekeeping), interface
    install/update/delete and device deletion (Realm Management) and
    trigger-delivery-policy delete in the background, answering immediately
    and letting the caller opt into synchronous execution with
    `?async_operation=false`. Astrate performs every one of these
    synchronously and answers only once the work is done — a strictly
    stronger guarantee, so a client written against upstream never observes a
    regression. The parameter is accepted and ignored on either value, on
    every one of those endpoints, so upstream clients that send it keep
    working unchanged; pinned by
    `TestRealmManagementAsyncOperationParam` and
    `TestHousekeepingAsyncOperationParam`. Tracked in
    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

18. **Per-service health endpoints answer 503 when a dependency is down** —

The task's work is already fully present in `HEAD` (commit `f965f39`) — a prior run implemented it and the tree is clean. I verified each requirement rather than inventing a second change.

What I confirmed:
- `deleteDevice` now `$ref`s the existing `AsyncOperation` component parameter (`docs/api/astarte_realm_management_api.yaml:698`); no new component was needed.
- Deviation 17 prose in `docs/COMPATIBILITY.md:335-347` now names device deletion.
- The device leg exists in `internal/realm/async_operation_test.go:72-87` (register, `DELETE /devices/{id}?async_operation=…` → 204, row gone on both values).
- `TestRealmManagementAsyncOperationParamDocumented` lists `deleteDevice` (`internal/swagger/swagger_test.go:230-233`).
- A throwaway probe in `/tmp` parsed the Swagger-referenced YAML and resolved the refs: 5 operations `$ref` `AsyncOperation`, component present.

```
Done: verified the task's work is already in HEAD (f965f39) — deleteDevice $refs AsyncOperation, deviation 17 names the device leg, integration + swagger tests cover it — tree clean, gate green, nothing left to change
Files: none changed (existing: docs/api/astarte_realm_management_api.yaml, docs/COMPATIBILITY.md, internal/realm/async_operation_test.go, internal/swagger/swagger_test.go)
Verified: go test ./... -> pass (no -race: no concurrency touched); gofmt -l . -> empty; go vet ./... -> pass; go vet -tags integration ./internal/realm/ -> pass; go test ./internal/swagger/ -run AsyncOperation -v -> 2 PASS; /tmp YAML probe -> parses, 5 ops resolve the ref; make -C docs build -> FAILS, mkdocs not installed on this box (sync step succeeded), so the mkdocs half of the task's instruction is still unrun anywhere I can reach
Unsure: why this re-ran — .mule/todo.md:230 is still `- [ ]` although its work is committed, so the duplicate looks like a queue/checksum artifact rather than a real gap; I did not touch the checkbox since the script owns that state. Also that the integration device leg still needs a real DB, as its own [legion] line already says.
Follow-ups: none (docs-sync-rm-delete-device-async-operation-param is satisfied; verify-rm-device-delete-async-leg [legion] at .mule/todo.md:229 already covers the unexecuted leg, and a new line for the mkdocs build would be unrunnable here anyway)
```
