package realm

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// renderDetailedInterface materialises every mapping default the 1.4
// detailed-listing contract promises (reliability unreliable, retention
// discard, database_retention_policy no_ttl, ...) over a stored definition.
// def is parsed with interfaceschema.ParseInterface — the stored truth — so
// the emitted defaults are exactly the ones parsing applies. Field order is
// fixed and empty description/doc are omitted; the output is deterministic
// for a given stored definition.
func renderDetailedInterface(def []byte) (json.RawMessage, error) {
	iface, err := interfaceschema.ParseInterface(def)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.WriteString(`{"interface_name":`)
	b.Write(quoteJSON(iface.Name))
	fmt.Fprintf(&b, `,"version_major":%d,"version_minor":%d`, iface.Major, iface.Minor)
	b.WriteString(`,"type":`)
	b.Write(quoteJSON(iface.Type.String()))
	b.WriteString(`,"ownership":`)
	b.Write(quoteJSON(iface.Ownership.String()))
	b.WriteString(`,"aggregation":`)
	b.Write(quoteJSON(iface.Aggregation.String()))
	if iface.Description != "" {
		b.WriteString(`,"description":`)
		b.Write(quoteJSON(iface.Description))
	}
	if iface.Doc != "" {
		b.WriteString(`,"doc":`)
		b.Write(quoteJSON(iface.Doc))
	}
	b.WriteString(`,"mappings":[`)
	for i := range iface.Mappings {
		if i > 0 {
			b.WriteByte(',')
		}
		writeDetailedMapping(&b, &iface.Mappings[i], iface.Type == interfaceschema.Datastream)
	}
	b.WriteString(`]}`)
	return json.RawMessage(b.Bytes()), nil
}

// writeDetailedMapping emits one mapping document. Datastream mappings carry
// every delivery default in contract order, plus database_retention_ttl only
// when a TTL is in force; properties mappings carry allow_unset.
func writeDetailedMapping(b *bytes.Buffer, m *interfaceschema.Mapping, datastream bool) {
	b.WriteString(`{"endpoint":`)
	b.Write(quoteJSON(m.Endpoint))
	b.WriteString(`,"type":`)
	b.Write(quoteJSON(m.Type.String()))
	if !datastream {
		fmt.Fprintf(b, `,"allow_unset":%t}`, m.AllowUnset)
		return
	}
	b.WriteString(`,"reliability":`)
	b.Write(quoteJSON(m.Reliability.String()))
	b.WriteString(`,"retention":`)
	b.Write(quoteJSON(m.Retention.String()))
	fmt.Fprintf(b, `,"expiry":%d,"explicit_timestamp":%t,"database_retention_policy":`,
		m.Expiry, m.ExplicitTimestamp)
	b.Write(quoteJSON(m.DatabaseRetentionPolicy.String()))
	if m.DatabaseRetentionPolicy == interfaceschema.UseTTL {
		fmt.Fprintf(b, `,"database_retention_ttl":%d`, m.DatabaseRetentionTTL)
	}
	b.WriteByte('}')
}
