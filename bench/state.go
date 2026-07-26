package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// State is everything provision produces and the measurement subcommands
// consume. It contains private keys and credentials secrets — treat the file
// as throwaway benchmark material, not production secrets.
type State struct {
	// BaseURL is the API entry the deployment was provisioned against, e.g.
	// http://api.astarte.localhost (Astarte compose) or http://127.0.0.1:8080
	// (Astrate compose). Service paths are derived from it; individual
	// endpoint overrides used at provision time are stored resolved.
	BaseURL   string    `json:"base_url"`
	Endpoints Endpoints `json:"endpoints"`

	Realm         string `json:"realm"`
	RealmKeyPEM   string `json:"realm_key_pem"`   // RS256 signer for a_rma/a_aea/a_pa JWTs
	BrokerURL     string `json:"broker_url"`      // from pairing /info, may be overridden per run
	TLSSkipVerify bool   `json:"tls_skip_verify"` // both dev composes use self-signed broker certs

	// Devices[:Probes] are reserved as e2e latency probes by ingest; the
	// workload devices come after them.
	Probes  int      `json:"probes"`
	Devices []Device `json:"devices"`
}

// Device is one registered device identity.
type Device struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func loadState(path string) (*State, error) {
	raw, err := os.ReadFile(path) //nolint:gosec // user-supplied state file path is the CLI contract
	if err != nil {
		return nil, fmt.Errorf("reading state file (run `bench provision` first): %w", err)
	}
	var st State
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, fmt.Errorf("parsing state file %s: %w", path, err)
	}
	if len(st.Devices) == 0 {
		return nil, fmt.Errorf("state file %s has no devices", path)
	}
	return &st, nil
}

func (s *State) save(path string) error {
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o600)
}

// workload returns the probe and workload device slices for a run of n
// workload devices.
func (s *State) workload(n int) (probes, workers []Device, err error) {
	if s.Probes+n > len(s.Devices) {
		return nil, nil, fmt.Errorf("state has %d devices (%d probes + %d workload available), need %d workload — re-provision with more",
			len(s.Devices), s.Probes, len(s.Devices)-s.Probes, n)
	}
	return s.Devices[:s.Probes], s.Devices[s.Probes : s.Probes+n], nil
}
