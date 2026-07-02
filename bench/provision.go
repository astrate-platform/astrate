package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

// cmdProvision creates the realm, installs the bench interfaces, registers
// probe + workload devices, and writes the state file everything else reads.
func cmdProvision(args []string) error {
	fs := flag.NewFlagSet("provision", flag.ExitOnError)
	baseURL := fs.String("base-url", "", "API base URL (e.g. http://api.astarte.localhost or http://127.0.0.1:8080); service paths are derived astartectl-style")
	hkURL := fs.String("housekeeping-url", "", "override the derived housekeeping base")
	rmURL := fs.String("realm-url", "", "override the derived realm-management base")
	pairURL := fs.String("pairing-url", "", "override the derived pairing base")
	aeaURL := fs.String("appengine-url", "", "override the derived appengine base")
	hkKeyPath := fs.String("housekeeping-key", "", "PEM private key authorized for the housekeeping API (required)")
	realm := fs.String("realm", "bench", "realm name to create (or reuse)")
	devices := fs.Int("devices", 200, "workload devices to register")
	probes := fs.Int("probes", 3, "extra devices reserved as e2e latency probes")
	concurrency := fs.Int("concurrency", 8, "parallel device registrations")
	statePath := fs.String("state", "bench-state.json", "state file to write")
	skipVerify := fs.Bool("insecure-tls", true, "skip broker/API TLS verification (dev composes use self-signed certs)")
	brokerOverride := fs.String("broker-url", "", "override the broker URL advertised by pairing (host networking quirks)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *baseURL == "" {
		return fmt.Errorf("-base-url is required")
	}
	if *hkKeyPath == "" {
		return fmt.Errorf("-housekeeping-key is required (Astarte: compose/housekeeping_private.pem from generate-compose-files.sh; Astrate: bench/keys/housekeeping.pem from up-astrate.sh)")
	}

	ep := deriveEndpoints(*baseURL)
	for _, override := range []struct {
		val *string
		dst *string
	}{{hkURL, &ep.Housekeeping}, {rmURL, &ep.RealmMgmt}, {pairURL, &ep.Pairing}, {aeaURL, &ep.AppEngine}} {
		if *override.val != "" {
			*override.dst = *override.val
		}
	}

	hkPEM, err := os.ReadFile(*hkKeyPath)
	if err != nil {
		return err
	}
	hkKey, err := parsePrivateKeyPEM(hkPEM)
	if err != nil {
		return fmt.Errorf("housekeeping key: %w", err)
	}

	// Realm signing key: generated fresh, registered with the realm, kept in
	// the state file so measurement runs can mint a_pa/a_aea tokens.
	realmKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	realmKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(realmKey)})
	pubDER, err := x509.MarshalPKIXPublicKey(&realmKey.PublicKey)
	if err != nil {
		return err
	}
	realmPubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	c := newClient(30 * time.Second)

	fmt.Printf("creating realm %q…\n", *realm)
	if err := c.createRealm(ep, hkKey, *realm, string(realmPubPEM)); err != nil {
		return fmt.Errorf("creating realm: %w", err)
	}

	for _, name := range []string{ifaceIndividual, ifaceObject} {
		raw, err := interfaceFS.ReadFile("interfaces/" + name + ".json")
		if err != nil {
			return err
		}
		fmt.Printf("installing interface %s…\n", name)
		if err := c.installInterface(ep, realmKey, *realm, raw); err != nil {
			return fmt.Errorf("installing %s: %w", name, err)
		}
	}

	total := *devices + *probes
	fmt.Printf("registering %d devices (%d workload + %d probes)…\n", total, *devices, *probes)
	regs := make([]Device, total)
	var (
		wg     sync.WaitGroup
		sem    = make(chan struct{}, *concurrency)
		mu     sync.Mutex
		regErr error
	)
	for i := range regs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			id, err := randomDeviceID()
			if err == nil {
				var secret string
				secret, err = c.registerDevice(ep, realmKey, *realm, id)
				regs[i] = Device{ID: id, Secret: secret}
			}
			if err != nil {
				mu.Lock()
				if regErr == nil {
					regErr = err
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if regErr != nil {
		return fmt.Errorf("registering devices: %w", regErr)
	}

	broker := *brokerOverride
	if broker == "" {
		broker, err = c.brokerURL(ep, *realm, regs[0])
		if err != nil {
			return fmt.Errorf("resolving broker URL: %w", err)
		}
	}

	st := &State{
		BaseURL:       *baseURL,
		Endpoints:     ep,
		Realm:         *realm,
		RealmKeyPEM:   string(realmKeyPEM),
		BrokerURL:     broker,
		TLSSkipVerify: *skipVerify,
		Probes:        *probes,
		Devices:       regs,
	}
	if err := st.save(*statePath); err != nil {
		return err
	}
	fmt.Printf("provisioned: realm=%s devices=%d broker=%s → %s\n", *realm, total, broker, *statePath)
	return nil
}
