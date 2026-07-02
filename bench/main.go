// Command bench is a wire-protocol benchmark harness that drives an Astarte
// or Astrate deployment through identical MQTT + REST traffic, so the two
// platforms can be compared apples-to-apples (see README.md for the
// methodology and the fairness caveats, especially around PUBACK semantics).
//
// Subcommands:
//
//	provision  create a realm, install the bench interfaces, register devices
//	ingest     N devices publish at a fixed rate; e2e/PUBACK latency + loss
//	connstorm  time a mass mTLS connect of pre-registered devices
//	query      AppEngine read mix against previously ingested data
//
// All state produced by provision (realm, keys, device secrets) lands in a
// JSON state file consumed by the other subcommands, so one provision run
// can back many measurement runs.
package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed interfaces/*.json
var interfaceFS embed.FS

const (
	ifaceIndividual = "org.astrate.bench.Individual"
	ifaceObject     = "org.astrate.bench.Object"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "provision":
		err = cmdProvision(os.Args[2:])
	case "ingest":
		err = cmdIngest(os.Args[2:])
	case "connstorm":
		err = cmdConnstorm(os.Args[2:])
	case "query":
		err = cmdQuery(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "bench: unknown subcommand %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "bench %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: bench <subcommand> [flags]

subcommands:
  provision   create realm + interfaces + devices, write the state file
  ingest      publish load and measure e2e latency, PUBACK latency, loss
  connstorm   mass-connect devices over mTLS and time it
  query       run an AppEngine read mix and report latency percentiles

Run "bench <subcommand> -h" for the flags of each subcommand.
`)
}
