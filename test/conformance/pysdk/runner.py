#!/usr/bin/env python3
"""CP-D Python SDK conformance runner (docs/ROADMAP.md §10 file 9.3).

Driven by pysdk_test.go: it composes a live Astrate instance, registers a
device, and invokes this script with the pairing URL, realm, device id,
credentials secret, and an interface definition. The script drives the
*unmodified* official astarte-device-sdk-python through the device loop —
connect, send an individual datastream, set a property — and the Go harness
cross-checks the persisted rows. A non-zero exit fails the checkpoint.

This runs where the pinned SDK is installable (the Linux CI/nightly job). The
Go harness skips it when the SDK cannot be imported.
"""
import argparse
import json
import os
import sys
import time
from pathlib import Path

from astarte.device import DeviceMqtt


def parse_args():
    p = argparse.ArgumentParser()
    p.add_argument("--pairing-url", required=True)
    p.add_argument("--realm", required=True)
    p.add_argument("--device-id", required=True)
    p.add_argument("--secret", required=True)
    p.add_argument("--persistency-dir", required=True)
    p.add_argument("--datastream-interface", required=True, help="path to the datastream interface JSON")
    p.add_argument("--properties-interface", required=True, help="path to the properties interface JSON")
    return p.parse_args()


def main():
    args = parse_args()

    # The SDK requires the persistency directory to already exist
    # (PersistencyDirectoryNotFoundError otherwise).
    os.makedirs(args.persistency_dir, exist_ok=True)

    device = DeviceMqtt(
        device_id=args.device_id,
        realm=args.realm,
        credentials_secret=args.secret,
        pairing_base_url=args.pairing_url,
        persistency_dir=args.persistency_dir,
        ignore_ssl_errors=True,  # the test broker presents a self-signed server cert
    )
    # add_interface_from_file calls Path methods on its argument.
    device.add_interface_from_file(Path(args.datastream_interface))
    device.add_interface_from_file(Path(args.properties_interface))

    with open(args.datastream_interface, encoding="utf-8") as f:
        ds_name = json.load(f)["interface_name"]
    with open(args.properties_interface, encoding="utf-8") as f:
        prop_name = json.load(f)["interface_name"]

    device.connect()
    for _ in range(100):
        if device.is_connected():
            break
        time.sleep(0.1)
    else:
        print("device did not connect", file=sys.stderr)
        return 1
    # The SDK flips to connected BEFORE it publishes the introspection from the
    # network thread; give that queue a moment so our data can't overtake it.
    time.sleep(0.5)

    # Individual datastream and a property set (the values the harness asserts
    # against the database). send() routes by interface type.
    device.send(ds_name, "/value", 42.5)
    device.send(prop_name, "/mode", "eco")
    time.sleep(1.5)

    device.disconnect()
    return 0


if __name__ == "__main__":
    sys.exit(main())
