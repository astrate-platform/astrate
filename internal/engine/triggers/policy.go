package triggers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

// Strategy is what a policy prescribes for one failed delivery attempt.
type Strategy int

const (
	// StrategyDiscard discards the event after a failed attempt.
	StrategyDiscard Strategy = iota
	// StrategyRetry re-delivers the event after a failed attempt.
	StrategyRetry
)

// StatusTransport is the pseudo-status passed to Decide when the request
// never produced a response at all.
const StatusTransport = 0

// handlerOnKind classifies the "on" clause of an error handler.
type handlerOnKind int

const (
	onKeyword  handlerOnKind = iota // "any_error", "client_error", "server_error"
	onExplicit                      // explicit status code list
)

// compiledHandler is one parsed error-handler clause.
type compiledHandler struct {
	onKind   handlerOnKind
	keyword  string // "any_error", "client_error", "server_error" when onKind == onKeyword
	codes    []int  // explicit codes when onKind == onExplicit
	strategy Strategy
}

// claims reports whether this handler applies to the given HTTP status code.
// For transport failures (status == 0) only keyword handlers match.
func (h *compiledHandler) claims(status int) bool {
	if status == StatusTransport {
		if h.onKind != onKeyword {
			return false
		}
		return h.keyword == "any_error" || h.keyword == "server_error"
	}
	switch h.onKind {
	case onKeyword:
		switch h.keyword {
		case "any_error":
			return status >= 400 && status <= 599
		case "client_error":
			return status >= 400 && status <= 499
		case "server_error":
			return status >= 500 && status <= 599
		}
		return false
	case onExplicit:
		for _, c := range h.codes {
			if c == status {
				return true
			}
		}
		return false
	}
	return false
}

// onLabel returns a human-readable label for the "on" clause, used in
// decision reasons.
func (h *compiledHandler) onLabel() string {
	if h.onKind == onKeyword {
		return h.keyword
	}
	if len(h.codes) == 1 {
		return fmt.Sprintf("status %d", h.codes[0])
	}
	return fmt.Sprintf("status codes %v", h.codes)
}

// Policy is a compiled trigger-delivery policy that decides whether to retry
// or discard a failed delivery attempt.
type Policy struct {
	Name            string
	MaximumCapacity int
	RetryTimes      int           // 0 when no handler retries
	EventTTL        time.Duration // 0 when unset
	handlers        []compiledHandler
}

var policyNameRE = regexp.MustCompile(`^[a-zA-Z0-9_.~-]{1,128}$`)

// CompilePolicy parses and validates one stored policy definition.
func CompilePolicy(def []byte) (*Policy, error) {
	var doc struct {
		Name          string `json:"name"`
		ErrorHandlers []struct {
			On       json.RawMessage `json:"on"`
			Strategy string          `json:"strategy"`
		} `json:"error_handlers"`
		MaximumCapacity int  `json:"maximum_capacity"`
		RetryTimes      *int `json:"retry_times"`
		EventTTL        *int `json:"event_ttl"`
	}
	if err := json.Unmarshal(def, &doc); err != nil {
		return nil, fmt.Errorf("triggers: policy does not parse as JSON: %w", err)
	}
	if !policyNameRE.MatchString(doc.Name) {
		return nil, fmt.Errorf("triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]")
	}
	if len(doc.ErrorHandlers) == 0 {
		return nil, fmt.Errorf("triggers: policy must have at least one error handler")
	}

	handlers := make([]compiledHandler, 0, len(doc.ErrorHandlers))
	hasRetry := false
	for i, h := range doc.ErrorHandlers {
		var strat Strategy
		switch h.Strategy {
		case "discard":
			strat = StrategyDiscard
		case "retry":
			strat = StrategyRetry
			hasRetry = true
		default:
			return nil, fmt.Errorf("triggers: error handler strategy must be discard or retry")
		}
		ch, err := compileHandlerOn(h.On)
		if err != nil {
			return nil, fmt.Errorf("triggers: error handler %d: %w", i+1, err)
		}
		ch.strategy = strat
		handlers = append(handlers, ch)
	}

	if hasRetry {
		if doc.RetryTimes == nil || *doc.RetryTimes < 1 || *doc.RetryTimes > 100 {
			return nil, fmt.Errorf("triggers: retry_times must be 1-100 when any handler retries")
		}
	} else if doc.RetryTimes != nil {
		return nil, fmt.Errorf("triggers: retry_times requires a retry handler")
	}
	if doc.MaximumCapacity < 1 {
		return nil, fmt.Errorf("triggers: maximum_capacity must be a positive integer")
	}
	if doc.EventTTL != nil && *doc.EventTTL < 0 {
		return nil, fmt.Errorf("triggers: event_ttl must be non-negative")
	}

	retryTimes := 0
	if doc.RetryTimes != nil {
		retryTimes = *doc.RetryTimes
	}
	eventTTL := time.Duration(0)
	if doc.EventTTL != nil {
		eventTTL = time.Duration(*doc.EventTTL) * time.Second
	}

	return &Policy{
		Name:            doc.Name,
		MaximumCapacity: doc.MaximumCapacity,
		RetryTimes:      retryTimes,
		EventTTL:        eventTTL,
		handlers:        handlers,
	}, nil
}

// compileHandlerOn parses the "on" member of an error handler, accepting the
// three keyword forms, a bare array of ints, or the object form with
// "custom_status_codes".
func compileHandlerOn(raw json.RawMessage) (compiledHandler, error) {
	var keyword string
	if err := json.Unmarshal(raw, &keyword); err == nil {
		switch keyword {
		case "any_error", "client_error", "server_error":
			return compiledHandler{onKind: onKeyword, keyword: keyword}, nil
		default:
			return compiledHandler{}, fmt.Errorf("on-keyword must be any_error, client_error, or server_error")
		}
	}

	var codes []int
	if err := json.Unmarshal(raw, &codes); err != nil {
		var custom struct {
			CustomStatusCodes []int `json:"custom_status_codes"`
		}
		if err := json.Unmarshal(raw, &custom); err != nil || custom.CustomStatusCodes == nil {
			return compiledHandler{}, fmt.Errorf("on must be a keyword or a status-code array")
		}
		codes = custom.CustomStatusCodes
	}
	if len(codes) == 0 {
		return compiledHandler{}, fmt.Errorf("status-code list must be non-empty")
	}
	for _, c := range codes {
		if c < 400 || c > 599 {
			return compiledHandler{}, fmt.Errorf("status codes must be in 400..599")
		}
	}
	return compiledHandler{onKind: onExplicit, codes: codes}, nil
}

// Decision is what a policy prescribes, plus why — the reason is logged by
// the executor so an operator can see which rule governed a delivery.
type Decision struct {
	Strategy Strategy
	Reason   string
}

// Decide returns the decision for a finished attempt. status is the HTTP
// status code, or StatusTransport when the request never produced a response.
func (p *Policy) Decide(status int) Decision {
	if p == nil {
		return Decision{Strategy: StrategyDiscard, Reason: "no policy applies"}
	}
	if status >= 200 && status <= 399 {
		return Decision{Strategy: StrategyDiscard, Reason: "not a failure status code"}
	}
	for i := range p.handlers {
		h := &p.handlers[i]
		if !h.claims(status) {
			continue
		}
		strat := "discard"
		if h.strategy == StrategyRetry {
			strat = "retry"
		}
		if status == StatusTransport {
			return Decision{
				Strategy: h.strategy,
				Reason:   fmt.Sprintf("transport failure, treated as a server error; handler %d (%s): %s", i+1, h.onLabel(), strat),
			}
		}
		return Decision{
			Strategy: h.strategy,
			Reason:   fmt.Sprintf("handler %d (%s): %s", i+1, h.onLabel(), strat),
		}
	}
	if status == StatusTransport {
		return Decision{Strategy: StrategyDiscard, Reason: "transport failure, treated as a server error; no handler claims it: discard"}
	}
	return Decision{Strategy: StrategyDiscard, Reason: fmt.Sprintf("no handler claims status %d: discard", status)}
}
