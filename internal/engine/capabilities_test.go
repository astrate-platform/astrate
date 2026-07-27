package engine

import (
	"testing"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/astrate-platform/astrate/internal/store"
)

func TestDecodeCapabilities(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "empty payload",
			payload: nil,
			wantErr: true,
		},
		{
			name:    "invalid BSON",
			payload: []byte("garbage"),
			wantErr: true,
		},
		{
			name:    "zlib compression capability",
			payload: mustBSON(t, bson.M{"purge_properties_compression_format": "zlib"}),
			want:    map[string]string{"purge_properties_compression_format": "zlib"},
		},
		{
			name:    "plaintext compression capability",
			payload: mustBSON(t, bson.M{"purge_properties_compression_format": "plaintext"}),
			want:    map[string]string{"purge_properties_compression_format": "plaintext"},
		},
		{
			name:    "unknown capability accepted",
			payload: mustBSON(t, bson.M{"future_cap": "value"}),
			want:    map[string]string{"future_cap": "value"},
		},
		{
			name:    "non-string values ignored",
			payload: mustBSON(t, bson.M{"purge_properties_compression_format": "zlib", "numeric": 42}),
			want:    map[string]string{"purge_properties_compression_format": "zlib"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCapabilities(tc.payload)
			if (err != nil) != tc.wantErr {
				t.Fatalf("decodeCapabilities() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("decodeCapabilities()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestValidateCapability(t *testing.T) {
	cases := []struct {
		key     string
		value   string
		wantErr bool
	}{
		{capPurgePropertiesCompressionFormat, "zlib", false},
		{capPurgePropertiesCompressionFormat, "plaintext", false},
		{capPurgePropertiesCompressionFormat, "gzip", true},
		{capPurgePropertiesCompressionFormat, "", true},
		{"unknown_key", "anything", false},
	}
	for _, tc := range cases {
		t.Run(tc.key+"/"+tc.value, func(t *testing.T) {
			err := validateCapability(tc.key, tc.value)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateCapability(%q, %q) error = %v, wantErr %v", tc.key, tc.value, err, tc.wantErr)
			}
		})
	}
}

func TestHandleCapabilities(t *testing.T) {
	rig, _, _ := newWiredRig(t, Config{})

	// Valid zlib capability.
	ack := &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
		"purge_properties_compression_format": "zlib",
	}), ack))
	if !ack.acked() {
		t.Error("zlib capability not acknowledged")
	}
	dev := rig.e.devices.peek(realmAlpha, devAlpha)
	if f := dev.purgeCompressionFormat(); f != "zlib" {
		t.Errorf("purgeCompression after zlib = %q, want zlib", f)
	}

	// Valid plaintext capability overrides.
	ack = &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
		"purge_properties_compression_format": "plaintext",
	}), ack))
	if !ack.acked() {
		t.Error("plaintext capability not acknowledged")
	}
	if f := dev.purgeCompressionFormat(); f != "plaintext" {
		t.Errorf("purgeCompression after plaintext = %q, want plaintext", f)
	}
}

func TestHandleCapabilitiesRejects(t *testing.T) {
	rig, _, _ := newWiredRig(t, Config{})

	// Empty payload.
	ack := &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, nil, ack))
	if !ack.acked() {
		t.Error("empty capabilities not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonCapabilitiesInvalid)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1", reasonCapabilitiesInvalid, got)
	}

	// Invalid BSON.
	ack = &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, []byte("garbage"), ack))
	if !ack.acked() {
		t.Error("invalid BSON capabilities not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonCapabilitiesInvalid)); got != 2 {
		t.Errorf("rejects[%s] = %v, want 2", reasonCapabilitiesInvalid, got)
	}

	// Invalid value.
	ack = &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
		"purge_properties_compression_format": "gzip",
	}), ack))
	if !ack.acked() {
		t.Error("invalid value capabilities not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonCapabilitiesInvalid)); got != 3 {
		t.Errorf("rejects[%s] = %v, want 3", reasonCapabilitiesInvalid, got)
	}
}

func TestPlaintextConsumerProperties(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	const serverIface = "com.astrate.test.ServerProperties"

	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path: "/limits/maxConnections", Value: []byte("42"),
	})

	// Set plaintext capability.
	ack := &ackCounter{}
	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
		"purge_properties_compression_format": "plaintext",
	}), ack))
	if !ack.acked() {
		t.Fatal("capabilities not acknowledged")
	}

	// emptyCache triggers consumer/properties send.
	ack = &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Fatal("emptyCache not acknowledged")
	}

	base := realmAlpha + "/" + devAlpha.String()
	purges := port.publishedTo(base + "/control/consumer/properties")
	if len(purges) != 1 {
		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
	}

	// The payload must be raw plaintext, not a zlib frame.
	payload := string(purges[0].payload)
	want := serverIface + "/limits/maxConnections"
	if payload != want {
		t.Errorf("plaintext purge = %q, want %q", payload, want)
	}

	// Verify it is NOT a valid zlib frame (the 4-byte header + zlib stream).
	if _, err := inflateProperties(purges[0].payload); err == nil {
		t.Error("plaintext purge is a valid zlib frame; should be raw text")
	}
}

func TestZlibConsumerPropertiesDefault(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	const serverIface = "com.astrate.test.ServerProperties"

	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path: "/limits/maxConnections", Value: []byte("42"),
	})

	// emptyCache without any capability set → zlib (default).
	ack := &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Fatal("emptyCache not acknowledged")
	}

	base := realmAlpha + "/" + devAlpha.String()
	purges := port.publishedTo(base + "/control/consumer/properties")
	if len(purges) != 1 {
		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
	}

	// Must be a valid zlib frame.
	entries, err := inflateProperties(purges[0].payload)
	if err != nil {
		t.Fatalf("default purge is not zlib: %v", err)
	}
	want := serverIface + "/limits/maxConnections"
	if len(entries) != 1 || entries[0] != want {
		t.Errorf("purge entries = %v, want [%s]", entries, want)
	}
}

// TestPurgeCompressionFor tests the fallback logic.
func TestPurgeCompressionFor(t *testing.T) {
	if got := purgeCompressionFor(nil); got != compressionZlib {
		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)
	}
	dev := &deviceState{}
	if got := purgeCompressionFor(dev); got != compressionZlib {
		t.Errorf("purgeCompressionFor(empty) = %q, want zlib", got)
	}
	dev.setPurgeCompression("plaintext")
	if got := purgeCompressionFor(dev); got != "plaintext" {
		t.Errorf("purgeCompressionFor(plaintext) = %q, want plaintext", got)
	}
}

// mustBSON marshals a bson.M, failing the test on error.
func mustBSON(t *testing.T, v bson.M) []byte {
	t.Helper()
	b, err := bson.Marshal(v)
	if err != nil {
		t.Fatalf("bson.Marshal: %v", err)
	}
	return b
}
