package main

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

// TestDrainStageBudgetIsPerStage pins that every graceful-drain stage gets its
// own shutdownTimeout budget: a stage that burns all of it (a slow in-flight
// HTTP request has no per-request deadline) must not hand the stages after it
// an already-expired context, which would make StopFlow return at
// f.router.Drain and MarkRunningFlowsStopped skip every flow.
func TestDrainStageBudgetIsPerStage(t *testing.T) {
	restore := shutdownTimeout
	shutdownTimeout = 20 * time.Millisecond
	t.Cleanup(func() { shutdownTimeout = restore })

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	var later []error
	live := func(context.Context) error {
		// Burn the whole budget, the way a stuck in-flight request does.
		time.Sleep(2 * shutdownTimeout)
		return context.DeadlineExceeded
	}
	// The stages shutdown runs after srv.Shutdown, standing in for
	// flow.Manager.Shutdown, MarkRunningFlowsStopped and engine.Drain.
	for i := 0; i < 3; i++ {
		stage("first", log, live)
		stage("after", log, func(ctx context.Context) error {
			if _, ok := ctx.Deadline(); !ok {
				t.Error("stage context carries no deadline")
			}
			later = append(later, ctx.Err())
			return nil
		})
	}
	for i, err := range later {
		if err != nil {
			t.Errorf("stage %d after an exhausted stage got %v, want a live context", i, err)
		}
	}
}

// TestDrainStageBoundsItsOwnStage is the other half of the contract: the
// budget is a deadline, not just a fresh context, so a stage that never
// returns on its own is still cut off at shutdownTimeout.
func TestDrainStageBoundsItsOwnStage(t *testing.T) {
	restore := shutdownTimeout
	shutdownTimeout = 20 * time.Millisecond
	t.Cleanup(func() { shutdownTimeout = restore })

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	deadline := time.Time{}
	stage("hangs", log, func(ctx context.Context) error {
		<-ctx.Done()
		deadline, _ = ctx.Deadline()
		return ctx.Err()
	})
	if !deadline.After(time.Now().Add(-shutdownTimeout)) {
		t.Errorf("stage deadline %v is not within the last %v of the budget", deadline, shutdownTimeout)
	}
}
