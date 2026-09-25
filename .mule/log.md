# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
| 2026-09-05 | docs-sync-pairing-health-path | done | 186s | d73e225 |
| 2026-09-05 | docs-sync-pairing-register-404 | done | 106s | 44a7cae |
| 2026-09-05 | swagger-httptest-coverage | done | 267s | ba79b34 |
| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
| 2026-09-05 | flow-validate-source-sink | blocked | 360s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | flow-validate-dead-source-sink-recompute | done | 200s | 7445176 |
| 2026-09-06 | docs-sync-rm-datastream-retention-endpoint | done | 232s | 0574fd9 |
| 2026-09-06 | docs-sync-rm-interfaces-detailed-param | done | 352s | 4094145 |
| 2026-09-06 | housekeeping-tests | done | 455s | 7583a7d |
| 2026-09-06 | probe-props-resend-error-triggers | blocked | 331s | wrote nothing |
| 2026-09-06 | server-data-trigger-bus | done | 827s | e46ce11 |
| 2026-09-06 | docs-sync-appengine-by-alias-endpoints | done | 144s | 35359df |
| 2026-09-06 | docs-sync-appengine-group-endpoints | done | 365s | 6ba74b8 |
| 2026-09-06 | docs-sync-appengine-get-group-device | done | 114s | 1aa98d6 |
| 2026-09-07 | docs-sync-appengine-query-params-status | done | 547s | bff7e1b |
| 2026-09-07 | docs-sync-appengine-group-patch-status | done | 75s | 81ced8f |
| 2026-09-07 | docs-sync-appengine-data-422-interface-level | done | 153s | f53d1ee |
| 2026-09-07 | appengine-snapshot-ignores-query-params | blocked | 327s | wrote nothing |
| 2026-09-07 | appengine-group-token-roundtrip-test | blocked | 306s | lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive) |
| 2026-09-07 | docs-sync-appengine-data-output-params | done | 215s | 262fdd5 |
| 2026-09-08 | probe-required-mapping-flag | done | 426s | db9d06e |
| 2026-09-08 | compat-note-v14-rc | blocked | 326s | wrote nothing |
| 2026-09-08 | detailed-listing-required-encrypted-flags | done | 159s | 63b629e |
| 2026-09-08 | store-alias-values-taken-test | done | 85s | a07bd6a |
| 2026-09-08 | store-latest-individual-test | done | 134s | 741794c |
| 2026-09-08 | docs-sync-appengine-device-status-schema | done | 255s | 00da133 |
| 2026-09-08 | docs-sync-appengine-data-set-422 | blocked | 251s | wrote nothing |
| 2026-09-08 | docs-sync-appengine-downsample-min | done | 86s | db506f9 |
| 2026-09-09 | pairing-bearer-secret-test | done | 90s | 32544c0 |
| 2026-09-09 | pairing-remoteip-fallback-test | done | 469s | 3d6e790 |
| 2026-09-09 | docs-sync-hk-patch-endpoint | done | 403s | 9fc9f33 |
| 2026-09-09 | docs-sync-hk-retention-field | done | 142s | d6d7c7f |
| 2026-09-09 | realm-pure-helper-tests | done | 403s | 39eb70a |
| 2026-09-09 | docs-sync-native-compat-health-503 | done | 202s | 656e21e |
| 2026-09-09 | docs-sync-native-compat-version-endpoints | done | 480s | d148f16 |
| 2026-09-10 | flowapi-autorestart-terminal-failure | done | 251s | 8ce04d2 |
| 2026-09-10 | flowapi-validationdetail-test | done | 121s | 832d8fb |
| 2026-09-10 | flowapi-autorestart-default-test | done | 148s | e272cfa |
| 2026-09-10 | docs-sync-pairing-status-enum | done | 89s | 2f6176c |
| 2026-09-10 | docs-sync-pairing-version-endpoint | done | 95s | c867b2c |
| 2026-09-11 | docs-sync-rm-delete-interface-status | done | 123s | b37d11f |
| 2026-09-11 | docs-sync-rm-mapping-required-encrypted | done | 137s | 5f8cb25 |
| 2026-09-11 | docs-sync-rm-put-auth-422 | done | 145s | 8554720 |
| 2026-09-11 | docs-sync-rm-version-example | blocked | 97s | opencode exited 1 |
| 2026-09-11 | interfaceschema-compat-attrs-flip-test | done | 329s | 4fe1f4f |
| 2026-09-11 | interfaceschema-object-attrs-uniformity | done | 344s | c6c8362 |
| 2026-09-11 | interfaceschema-malformed-placeholder-fixtures | done | 130s | f6a44a2 |
| 2026-09-11 | interfaceschema-enum-roundtrip-coverage | done | 185s | f955958 |
| 2026-09-11 | interfaceschema-trie-interior-literal-param | done | 121s | 705795f |
| 2026-09-12 | docs-sync-ae-write-405 | done | 400s | 9ed08c3 |
| 2026-09-12 | docs-sync-ae-write-value-422 | done | 338s | 17f6f70 |
| 2026-09-12 | docs-sync-ae-post-groups-409 | done | 128s | d6865a1 |
| 2026-09-12 | payload-longinteger-fraction-quantize | done | 356s | d7babc6 |
| 2026-09-12 | payload-json-malformed-t-not-tolerated | done | 99s | 5b56557 |
| 2026-09-12 | docs-sync-ae-post-group-devices-409 | done | 179s | cbf81af |
| 2026-09-12 | docs-sync-ae-patch-merge-patch-media-type | done | 137s | 64361d9 |
| 2026-09-13 | store-todo-lttb-stale | done | 66s | 829f408 |
| 2026-09-13 | config-fail-loud-engine-shards | done | 150s | 83482e4 |
| 2026-09-13 | config-fail-loud-dev-mode | done | 227s | 27d8165 |
| 2026-09-13 | config-key-resolvers-test | done | 143s | 9e68a01 |
| 2026-09-13 | docs-sync-hk-delete-gating-responses | done | 331s | d07b6c7 |
| 2026-09-13 | docs-sync-hk-validation-example | done | 68s | d762f00 |
| 2026-09-14 | auth-iat-not-required-test | done | 257s | 6524c7f |
| 2026-09-14 | auth-empty-claims-403-test | done | 348s | b91247c |
| 2026-09-14 | auth-cache-default-size-test | done | 157s | 874797b |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 904s | 015a706 |
| 2026-09-15 | hk-zero-reglimit-create-asymmetry | blocked | 37s | opencode exited 1 |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 567s | 203af1c |
| 2026-09-15 | docs-sync-native-socket-query-token-auth | done | 98s | 8148aae |
| 2026-09-15 | astarteapi-metadata-envelope-test | done | 238s | 1365cb8 |
| 2026-09-16 | astarteapi-multikey-sorted-golden | done | 215s | bfffbed |
| 2026-09-16 | docs-sync-pairing-deviceendpoints-dead-404-403 | done | 379s | f0d6a4d |
| 2026-09-16 | docs-sync-pairing-version-404-unreachable | done | 200s | bc73f65 |
| 2026-09-16 | observability-readiness-wedge-test | done | 178s | 4ef2e23 |
| 2026-09-16 | observability-dbpool-gauge-coverage | done | 62s | 19cd4f1 |
| 2026-09-16 | observability-health-content-type-test | done | 250s | f6c27f5 |
| 2026-09-17 | docs-sync-pairing-version-value | done | 1059s | c69a164 |
| 2026-09-17 | docs-sync-native-version-value | done | 503s | 08048b3 |
| 2026-09-17 | httpx-cors-vary-origin-passthrough | done | 143s | 986c100 |
| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
| 2026-09-18 | channels-group-watch-membership-scope | done | 677s | aeaeb47 |
| 2026-09-18 | channels-rejoin-joinref-tagging-test | done | 208s | d48e17f |
| 2026-09-19 | channels-rejoin-authz-mismatch | blocked | 179s | wrote nothing |
| 2026-09-20 | errorname-missing-required-test | done | 135s | 739df67 |
| 2026-09-20 | triggers-custom-action-policy-nodecide | done | 459s | 3c8b89a |
| 2026-09-20 | compat-note-custom-action-policy-boundary | done | 97s | 8cd17a8 |
| 2026-09-20 | docs-sync-native-metrics-example-fake-series | done | 93s | 43c4a8b |
| 2026-09-20 | docs-sync-native-socket-missing-403-500 | done | 116s | f9d8b27 |
| 2026-09-21 | store-pipelines-empty-name-zero-blocks-test | done | 91s | 82fc0a4 |
| 2026-09-21 | broker-acl-coldstart-fallback-flood | done | 218s | 0e5078f |
| 2026-09-21 | broker-offlineacl-entry-eviction | done | 406s | 72b2686 |
| 2026-09-22 | docs-sync-ae-read-query-params | blocked | 159s | wrote nothing |
| 2026-09-23 | examples-echo-container-contract-test | done | 251s | f35af21 |
| 2026-09-23 | flow-mqtt-source-reconnect-recovery | blocked | 84s | wrote nothing |
| 2026-09-23 | flow-msg-json-integer-precision | done | 612s | 4df3f6b |
| 2026-09-23 | flow-sort-bounded-buffer | done | 634s | 9bf0ab7 |
| 2026-09-23 | flow-randomsource-span-overflow | blocked | 288s | lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec) |
| 2026-09-23 | flow-filter-key-contains-test | done | 107s | 3875ba6 |
| 2026-09-23 | docs-sync-rm-put-interface-409 | done | 356s | 7256af9 |
| 2026-09-24 | docs-sync-rm-interface-422-shapes | done | 490s | 93e784e |
| 2026-09-24 | docs-sync-rm-validation-example-prefix | done | 238s | 02e269b |
| 2026-09-24 | docs-sync-rm-auth-403 | done | 668s | 9c05370 |
| 2026-09-24 | swagger-sub-failfast | done | 371s | 1d321ea |
| 2026-09-25 | drain-per-stage-budget | done | 265s | 84a050e |
| 2026-09-25 | cmd-devcert-fields-test | done | 357s | f1434eb |
| 2026-09-25 | cmd-healthcheck-contract-test | blocked | 199s | tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s) |
