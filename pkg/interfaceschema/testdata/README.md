# interfaceschema test fixtures

## valid/

The `org.astarte-platform.*` files are vendored verbatim from the upstream
Astarte standard interfaces, for SDK wire-compatibility parity:

- Source: https://github.com/astarte-platform/astarte, directory
  `standard-interfaces/`
- Pinned commit: `d6a1f5e578f84b90c06409ed8e95dc973917133a` (2025-12-11)
- Files: `org.astarte-platform.genericcommands.ServerCommands.json`,
  `org.astarte-platform.genericevents.DeviceEvents.json`,
  `org.astarte-platform.genericsensors.AvailableSensors.json`,
  `org.astarte-platform.genericsensors.Geolocation.json`,
  `org.astarte-platform.genericsensors.SamplingRate.json`,
  `org.astarte-platform.genericsensors.Values.json`

Do not edit the vendored files; re-vendor from upstream and update the pinned
commit here instead.

The `com.astrate.test.*` files are Astrate-authored fixtures that extend
coverage to every value type, object aggregation (including the depth-1
empty-prefix form), array-typed properties, and deep parametric endpoints.

## invalid/

Each `NN-*.json` file violates exactly one validation rule. `manifest.json`
maps every fixture to the substring its `ParseInterface` error must contain;
`parse_test.go` fails if a fixture is missing from the manifest or vice versa.
