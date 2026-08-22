package triggers

import (
	"testing"
)

func TestCompilePolicy(t *testing.T) {
	tests := []struct {
		name              string
		def               string
		wantErr           string
		wantPrefetchCount int // 0 means don't check
	}{
		{
			name:    "valid discard policy",
			def:     `{"name":"test","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":10}`,
			wantErr: "",
		},
		{
			name:    "valid retry policy",
			def:     `{"name":"test","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":10,"retry_times":3}`,
			wantErr: "",
		},
		{
			name:    "valid retry with event_ttl",
			def:     `{"name":"test","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":5,"retry_times":1,"event_ttl":60}`,
			wantErr: "",
		},
		{
			name:    "on keyword any_error",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "",
		},
		{
			name:    "on keyword client_error",
			def:     `{"name":"k","error_handlers":[{"on":"client_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "",
		},
		{
			name:    "on keyword server_error",
			def:     `{"name":"k","error_handlers":[{"on":"server_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "",
		},
		{
			name:    "on bare array",
			def:     `{"name":"k","error_handlers":[{"on":[404,503],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "",
		},
		{
			name:    "on custom_status_codes object",
			def:     `{"name":"k","error_handlers":[{"on":{"custom_status_codes":[429]},"strategy":"retry"}],"maximum_capacity":1,"retry_times":5}`,
			wantErr: "",
		},
		{
			name:    "retry_times required when retry handler",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1}`,
			wantErr: "retry_times must be 1-100",
		},
		{
			name:    "retry_times forbidden when no retry handler",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"retry_times":3}`,
			wantErr: "retry_times requires a retry handler",
		},
		{
			name:    "retry_times out of range low",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":0}`,
			wantErr: "retry_times must be 1-100",
		},
		{
			name:    "retry_times out of range high",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":101}`,
			wantErr: "retry_times must be 1-100",
		},
		{
			name:    "maximum_capacity too low",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":0}`,
			wantErr: "maximum_capacity must be a positive integer",
		},
		{
			name:    "negative event_ttl",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"event_ttl":-1}`,
			wantErr: "event_ttl must be non-negative",
		},
		{
			name:    "name starts with @",
			def:     `{"name":"@bad","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "policy name",
		},
		{
			name:    "no handlers",
			def:     `{"name":"k","error_handlers":[],"maximum_capacity":1}`,
			wantErr: "at least one error handler",
		},
		{
			name:    "invalid JSON",
			def:     `{bad json`,
			wantErr: "does not parse",
		},
		{
			name:    "empty on keyword",
			def:     `{"name":"k","error_handlers":[{"on":"bad","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "on-keyword must be",
		},
		{
			name:    "on empty array",
			def:     `{"name":"k","error_handlers":[{"on":[],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "status-code list must be non-empty",
		},
		{
			name:    "on code out of range",
			def:     `{"name":"k","error_handlers":[{"on":[200],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "status codes must be in 400..599",
		},
		// --- pairwise-disjoint error handlers (#65, probed upstream v1.2.0) ---
		{
			name:    "overlapping explicit lists rejected",
			def:     `{"name":"k","error_handlers":[{"on":[400,401],"strategy":"discard"},{"on":[401,402],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "any_error twice rejected",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"},{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "any_error plus explicit code rejected",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"},{"on":[450],"strategy":"retry"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "server_error plus any_error rejected",
			def:     `{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"},{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "client_error plus explicit 404 rejected cleanly",
			def:     `{"name":"k","error_handlers":[{"on":"client_error","strategy":"retry"},{"on":[404],"strategy":"discard"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "same code discarded twice rejected regardless of strategy",
			def:     `{"name":"k","error_handlers":[{"on":[500],"strategy":"discard"},{"on":[500],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "chained three-way overlap caught by pairwise checks",
			def:     `{"name":"k","error_handlers":[{"on":[400,401],"strategy":"discard"},{"on":[401,402],"strategy":"discard"},{"on":[402,403],"strategy":"discard"}],"maximum_capacity":1}`,
			wantErr: "must all handle distinct errors",
		},
		{
			name:    "disjoint explicit codes accepted",
			def:     `{"name":"k","error_handlers":[{"on":[400],"strategy":"discard"},{"on":[401],"strategy":"retry"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "",
		},
		{
			name:    "server_error plus client-range code accepted",
			def:     `{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"},{"on":[401],"strategy":"discard"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "",
		},
		{
			name:              "prefetch_count round-trips",
			def:               `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"prefetch_count":5}`,
			wantErr:           "",
			wantPrefetchCount: 5,
		},
		{
			name:              "prefetch_count defaults to 1 when omitted",
			def:               `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantErr:           "",
			wantPrefetchCount: 1,
		},
		{
			name:              "prefetch_count 1 accepted (lower bound)",
			def:               `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"prefetch_count":1}`,
			wantErr:           "",
			wantPrefetchCount: 1,
		},
		{
			name:              "prefetch_count 300 accepted (upper bound)",
			def:               `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"prefetch_count":300}`,
			wantErr:           "",
			wantPrefetchCount: 300,
		},
		{
			name:    "prefetch_count 0 rejected",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"prefetch_count":0}`,
			wantErr: "prefetch_count must be between 1 and 300",
		},
		{
			name:    "prefetch_count 301 rejected",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"prefetch_count":301}`,
			wantErr: "prefetch_count must be between 1 and 300",
		},
		{
			name:    "retry_times 1 accepted (lower bound)",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":1}`,
			wantErr: "",
		},
		{
			name:    "retry_times 100 accepted (upper bound)",
			def:     `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":100}`,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := CompilePolicy([]byte(tt.def))
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if p == nil {
					t.Fatal("expected non-nil policy")
				}
				if tt.wantPrefetchCount != 0 && p.PrefetchCount != tt.wantPrefetchCount {
					t.Errorf("PrefetchCount = %d, want %d", p.PrefetchCount, tt.wantPrefetchCount)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestDecide(t *testing.T) {
	retryServer := `{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
	retryAny := `{"name":"k","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
	discardAll := `{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`
	retryClient := `{"name":"k","error_handlers":[{"on":"client_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
	codes503 := `{"name":"k","error_handlers":[{"on":[503],"strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
	multiHandler := `{"name":"k","error_handlers":[{"on":"client_error","strategy":"discard"},{"on":"server_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`

	tests := []struct {
		name       string
		def        string
		status     int
		wantStrat  Strategy
		wantReason string
	}{
		{
			name:      "non-failure 200",
			def:       retryServer,
			status:    200,
			wantStrat: StrategyDiscard,
		},
		{
			name:      "non-failure 399",
			def:       retryServer,
			status:    399,
			wantStrat: StrategyDiscard,
		},
		{
			name:       "transport matched by server_error",
			def:        retryServer,
			status:     StatusTransport,
			wantStrat:  StrategyRetry,
			wantReason: "transport failure",
		},
		{
			name:       "transport matched by any_error",
			def:        retryAny,
			status:     StatusTransport,
			wantStrat:  StrategyRetry,
			wantReason: "transport failure",
		},
		{
			name:       "transport NOT matched by explicit [500]",
			def:        codes503,
			status:     StatusTransport,
			wantStrat:  StrategyDiscard,
			wantReason: "no handler claims",
		},
		{
			name:       "transport NOT matched by client_error",
			def:        retryClient,
			status:     StatusTransport,
			wantStrat:  StrategyDiscard,
			wantReason: "no handler claims",
		},
		{
			name:       "explicit 503 matched",
			def:        codes503,
			status:     503,
			wantStrat:  StrategyRetry,
			wantReason: "status 503",
		},
		{
			name:      "500 not matched by [503]",
			def:       codes503,
			status:    500,
			wantStrat: StrategyDiscard,
		},
		{
			name:      "418 unclaimed",
			def:       discardAll,
			status:    418,
			wantStrat: StrategyDiscard,
		},
		{
			name:       "handler precedence: first wins",
			def:        multiHandler,
			status:     404,
			wantStrat:  StrategyDiscard,
			wantReason: "handler 1",
		},
		{
			name:       "handler precedence: second wins for 500",
			def:        multiHandler,
			status:     500,
			wantStrat:  StrategyRetry,
			wantReason: "handler 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := CompilePolicy([]byte(tt.def))
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			d := p.Decide(tt.status)
			if d.Strategy != tt.wantStrat {
				t.Errorf("strategy: got %v, want %v", d.Strategy, tt.wantStrat)
			}
			if d.Reason == "" {
				t.Error("reason must be non-empty")
			}
			if tt.wantReason != "" && !contains(d.Reason, tt.wantReason) {
				t.Errorf("reason %q does not contain %q", d.Reason, tt.wantReason)
			}
		})
	}
}

func TestDecideNilPolicy(t *testing.T) {
	var p *Policy
	d := p.Decide(500)
	if d.Strategy != StrategyDiscard {
		t.Errorf("expected StrategyDiscard, got %v", d.Strategy)
	}
	if d.Reason == "" {
		t.Error("reason must be non-empty")
	}
	if !contains(d.Reason, "no policy") {
		t.Errorf("reason %q should mention no policy", d.Reason)
	}
}

func TestDecideTransportReasonDistinct(t *testing.T) {
	def := `{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
	p, err := CompilePolicy([]byte(def))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	dTransport := p.Decide(StatusTransport)
	dCode := p.Decide(500)
	if dTransport.Reason == dCode.Reason {
		t.Errorf("transport reason and status-code reason must be distinct, both got %q", dTransport.Reason)
	}
	if !contains(dTransport.Reason, "transport") {
		t.Errorf("transport reason should mention transport, got %q", dTransport.Reason)
	}
	if contains(dCode.Reason, "transport") {
		t.Errorf("status-code reason should not mention transport, got %q", dCode.Reason)
	}
}

func TestDecideEveryDecisionHasReason(t *testing.T) {
	defs := []string{
		`{"name":"k","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
		`{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`,
		`{"name":"k","error_handlers":[{"on":[404],"strategy":"discard"}],"maximum_capacity":1}`,
	}
	statuses := []int{StatusTransport, 400, 418, 500, 503}
	for _, def := range defs {
		p, err := CompilePolicy([]byte(def))
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		for _, s := range statuses {
			d := p.Decide(s)
			if d.Reason == "" {
				t.Errorf("empty reason for status %d on %s", s, def)
			}
		}
	}
	// nil policy
	var np *Policy
	d := np.Decide(500)
	if d.Reason == "" {
		t.Error("nil policy: empty reason")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
