package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandleMessageContract pins the HTTP message contract every container
// author copies from this example (see README.md and
// docs/handoff/flow-design-b-container-block-2026-07-29.md).
func TestHandleMessageContract(t *testing.T) {
	cases := []struct {
		name     string
		body     []byte
		wantCode int
		wantBody []byte
	}{
		{
			name:     "echoes valid JSON object verbatim",
			body:     []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","type":"string","data":"hello","timestamp_us":0}`),
			wantCode: http.StatusOK,
			wantBody: []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","type":"string","data":"hello","timestamp_us":0}`),
		},
		{
			name:     "echoes pretty-printed JSON byte-for-byte",
			body:     []byte("{\n  \"schema\": \"astarte_flow/message/v0.1\",\n  \"key\": \"demo\",\n  \"data\": 1\n}\n"),
			wantCode: http.StatusOK,
			wantBody: []byte("{\n  \"schema\": \"astarte_flow/message/v0.1\",\n  \"key\": \"demo\",\n  \"data\": 1\n}\n"),
		},
		{
			name:     "400 on invalid JSON",
			body:     []byte(`{"schema": `),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "400 on non-object JSON",
			body:     []byte(`["not","an","object"]`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "204 on empty body",
			body:     []byte{},
			wantCode: http.StatusNoContent,
		},
		{
			name:     "204 when metadata echo_drop is 1",
			body:     []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":"1"}}`),
			wantCode: http.StatusNoContent,
		},
		{
			name:     "echoes when echo_drop is not the string 1",
			body:     []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":"0"}}`),
			wantCode: http.StatusOK,
			wantBody: []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":"0"}}`),
		},
		{
			name:     "echoes when echo_drop is a non-string value",
			body:     []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":true}}`),
			wantCode: http.StatusOK,
			wantBody: []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":true}}`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/message", bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			handleMessage(rec, req)
			resp := rec.Result()
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != tc.wantCode {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tc.wantCode)
			}
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if tc.wantBody != nil && !bytes.Equal(got, tc.wantBody) {
				t.Fatalf("body = %q, want %q", got, tc.wantBody)
			}
			if tc.wantCode == http.StatusOK {
				if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ct)
				}
			}
		})
	}
}

// TestHandleMessageReadCap pins the 1 MiB read cap: an oversized body is
// truncated, so the truncated JSON is invalid and the request is rejected
// rather than echoed back unwound.
func TestHandleMessageReadCap(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/message",
		strings.NewReader(`{"schema":"astarte_flow/message/v0.1","data":"`+strings.Repeat("a", 2<<20)+`"}`))
	rec := httptest.NewRecorder()
	handleMessage(rec, req)
	resp := rec.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		got, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 400 (read cap truncates oversized JSON); body = %.80q", resp.StatusCode, got)
	}
}

// TestSanitizeLog pins that hostile input cannot forge log lines: every C0
// control character and DEL is stripped, printable output is kept.
func TestSanitizeLog(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"hello", "hello"},
		{"a\nb", "ab"},
		{"a\tb", "ab"},
		{"a\x00b", "ab"},
		{"a\x1fb", "ab"},
		{"a\x7fb", "ab"},
		{"a\r\nb", "ab"},
		{string([]rune{0xe4, 0x1a, 0xb8}), "ä¸"}, // printable kept, C0 stripped
	}
	for _, tc := range cases {
		if got := sanitizeLog(tc.in); got != tc.want {
			t.Errorf("sanitizeLog(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
