// Minimal Astrate Flow container contract: POST /v1/message with a Message
// JSON body and echo it back (or drop on empty array / 204).
//
// See README.md and docs/handoff/flow-design-b-container-block-2026-07-29.md.
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Optional opaque config from the container block (ASTRATE_FLOW_CONFIG).
	if cfg := os.Getenv("ASTRATE_FLOW_CONFIG"); cfg != "" {
		log.Printf("ASTRATE_FLOW_CONFIG=%s", strconv.Quote(sanitizeLog(cfg)))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /v1/message", handleMessage)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("flow-container-echo listening on %s", strconv.Quote(sanitizeLog(addr)))
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

// sanitizeLog drops control characters so hostile input cannot forge log lines.
func sanitizeLog(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Validate JSON object (Message); echo the bytes as-is on success.
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	// Optional drop: metadata "echo_drop" == "1" → 204 (filter semantics).
	if md, ok := probe["metadata"].(map[string]any); ok {
		if v, _ := md["echo_drop"].(string); v == "1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
