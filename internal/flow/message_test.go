package flow

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlowMessage_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name: "integer",
			msg: Message{
				Key:       "sensor/temp",
				Type:      TypeInteger,
				Data:      int64(42),
				Timestamp: 1551884045074181,
				Metadata:  map[string]string{"source": "device-1"},
			},
		},
		{
			name: "real",
			msg: Message{
				Key:       "sensor/pressure",
				Type:      TypeReal,
				Data:      3.14159,
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "boolean",
			msg: Message{
				Key:       "sensor/active",
				Type:      TypeBoolean,
				Data:      true,
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "string",
			msg: Message{
				Key:       "sensor/name",
				Type:      TypeString,
				Data:      "hello world",
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "binary",
			msg: Message{
				Key:       "sensor/blob",
				Type:      TypeBinary,
				Subtype:   "application/octet-stream",
				Data:      []byte{0x00, 0x01, 0xFF},
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "datetime",
			msg: Message{
				Key:       "sensor/time",
				Type:      TypeDatetime,
				Data:      time.Date(2019, 3, 5, 10, 0, 0, 0, time.UTC),
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "map with mixed types",
			msg: Message{
				Key:       "maps_stream",
				Type:      TypeMap,
				Timestamp: 1551884045074181,
				Metadata:  map[string]string{"hello": "world"},
				FieldTypes: map[string]DataType{
					"a": TypeReal,
					"b": TypeBinary,
				},
				FieldSubtypes: map[string]string{
					"b": "text/plain",
				},
				Data: map[string]any{
					"a": -1.0,
					"b": []byte("Hello"),
				},
			},
		},
		{
			name: "empty metadata",
			msg: Message{
				Key:       "stream/empty",
				Type:      TypeString,
				Data:      "test",
				Timestamp: 1000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(&tt.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal: err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			var got Message
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			if got.Key != tt.msg.Key {
				t.Errorf("Key = %q, want %q", got.Key, tt.msg.Key)
			}
			if got.Type != tt.msg.Type {
				t.Errorf("Type = %v, want %v", got.Type, tt.msg.Type)
			}
			if got.Timestamp != tt.msg.Timestamp {
				t.Errorf("Timestamp = %d, want %d", got.Timestamp, tt.msg.Timestamp)
			}
			if got.Subtype != tt.msg.Subtype {
				t.Errorf("Subtype = %q, want %q", got.Subtype, tt.msg.Subtype)
			}

			// Compare metadata.
			if len(got.Metadata) != len(tt.msg.Metadata) {
				t.Errorf("Metadata len = %d, want %d", len(got.Metadata), len(tt.msg.Metadata))
			} else {
				for k, v := range tt.msg.Metadata {
					if got.Metadata[k] != v {
						t.Errorf("Metadata[%q] = %q, want %q", k, got.Metadata[k], v)
					}
				}
			}

			// Compare data based on type.
			switch tt.msg.Type {
			case TypeInteger:
				if got.Data.(int64) != tt.msg.Data.(int64) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeReal:
				if got.Data.(float64) != tt.msg.Data.(float64) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeBoolean:
				if got.Data.(bool) != tt.msg.Data.(bool) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeString:
				if got.Data.(string) != tt.msg.Data.(string) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeBinary:
				gotBs := got.Data.([]byte)
				wantBs := tt.msg.Data.([]byte)
				if len(gotBs) != len(wantBs) {
					t.Errorf("Data len = %d, want %d", len(gotBs), len(wantBs))
				}
				for i := range gotBs {
					if gotBs[i] != wantBs[i] {
						t.Errorf("Data[%d] = %x, want %x", i, gotBs[i], wantBs[i])
					}
				}
			case TypeDatetime:
				gotT := got.Data.(time.Time)
				wantT := tt.msg.Data.(time.Time)
				if !gotT.Equal(wantT) {
					t.Errorf("Data = %v, want %v", gotT, wantT)
				}
			case TypeMap:
				gotMap := got.Data.(map[string]any)
				wantMap := tt.msg.Data.(map[string]any)
				if len(gotMap) != len(wantMap) {
					t.Errorf("Data map len = %d, want %d", len(gotMap), len(wantMap))
				}
				// Check field types.
				for k, v := range tt.msg.FieldTypes {
					if got.FieldTypes[k] != v {
						t.Errorf("FieldTypes[%q] = %v, want %v", k, got.FieldTypes[k], v)
					}
				}
				// Check binary in map is round-tripped.
				gotBin := gotMap["b"].([]byte)
				wantBin := wantMap["b"].([]byte)
				if len(gotBin) != len(wantBin) || gotBin[0] != wantBin[0] {
					t.Errorf("Data[b] = %x, want %x", gotBin, wantBin)
				}
			}
		})
	}
}

func TestFlowMessage_WireFormatSchema(t *testing.T) {
	msg := Message{
		Key:       "test",
		Type:      TypeString,
		Data:      "hi",
		Timestamp: 12345,
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}

	var schema string
	if err := json.Unmarshal(raw["schema"], &schema); err != nil {
		t.Fatal(err)
	}
	if schema != WireSchema {
		t.Errorf("schema = %q, want %q", schema, WireSchema)
	}
}

func TestFlowMessage_UnmarshalRejectsUnknownSchema(t *testing.T) {
	raw := `{"schema":"unknown/v1","key":"k","type":"string","data":"x","timestamp_us":0}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		t.Fatal("expected error for unknown schema")
	}
}

func TestFlowMessage_MapSubtypesOmittedWhenEmpty(t *testing.T) {
	msg := Message{
		Key:       "m",
		Type:      TypeMap,
		Timestamp: 1,
		FieldTypes: map[string]DataType{
			"x": TypeString,
		},
		Data: map[string]any{"x": "hello"},
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["subtype"]; ok {
		t.Error("subtype should be omitted when FieldSubtypes is nil")
	}
}

func TestFlowMessage_ScalarSubtypeOmittedWhenEmpty(t *testing.T) {
	msg := Message{
		Key:       "s",
		Type:      TypeString,
		Data:      "x",
		Timestamp: 1,
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["subtype"]; ok {
		t.Error("subtype should be omitted when Subtype is empty")
	}
}

func TestFlowMessage_BinaryBase64Encoding(t *testing.T) {
	msg := Message{
		Key:       "bin",
		Type:      TypeBinary,
		Subtype:   "application/octet-stream",
		Data:      []byte{0x00, 0x01, 0xFF},
		Timestamp: 100,
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	// Verify data is base64-encoded on wire.
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	s, ok := raw["data"].(string)
	if !ok {
		t.Fatalf("data is %T, want string (base64)", raw["data"])
	}
	if s != "AAH/" {
		t.Errorf("data = %q, want %q", s, "AAH/")
	}
}

func TestFlowMessage_MapBinaryBase64Encoding(t *testing.T) {
	msg := Message{
		Key:       "map_bin",
		Type:      TypeMap,
		Timestamp: 100,
		FieldTypes: map[string]DataType{
			"blob": TypeBinary,
		},
		FieldSubtypes: map[string]string{
			"blob": "image/png",
		},
		Data: map[string]any{
			"blob": []byte{0xDE, 0xAD},
		},
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	d := raw["data"].(map[string]any)
	s, ok := d["blob"].(string)
	if !ok {
		t.Fatalf("data.blob is %T, want string", d["blob"])
	}
	if s != "3q0=" {
		t.Errorf("data.blob = %q, want %q", s, "3q0=")
	}
}

func TestDataTypeStringRoundTrip(t *testing.T) {
	types := []DataType{TypeInteger, TypeReal, TypeBoolean, TypeDatetime, TypeBinary, TypeString}
	for _, dt := range types {
		s := dataTypeString(dt)
		got, err := parseDataType(s)
		if err != nil {
			t.Errorf("parseDataType(%q): %v", s, err)
			continue
		}
		if got != dt {
			t.Errorf("round-trip %v → %q → %v", dt, s, got)
		}
	}
}
