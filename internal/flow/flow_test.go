package flow

import (
	"testing"
)

func TestFlowStatus_String(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{FlowStatusCreating, "creating"},
		{FlowStatusRunning, "running"},
		{FlowStatusStopped, "stopped"},
		{FlowStatusFailed, "failed"},
		{Status(99), "unknown(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFlowStatus_Transitions(t *testing.T) {
	tests := []struct {
		name       string
		initialise func(t *testing.T) *Flow
		wantFinal  Status
	}{
		{
			name: "creating -> running",
			initialise: func(t *testing.T) *Flow {
				mgr := NewManager()
				blk := &passthroughBlock{}
				f, err := mgr.StartFlow(t.Context(), Config{
					PipelineID: "trans-1",
					Blocks:     []Block{blk},
				})
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = mgr.StopFlow(t.Context(), "trans-1") })
				return f
			},
			wantFinal: FlowStatusRunning,
		},
		{
			name: "creating -> running -> stopped",
			initialise: func(t *testing.T) *Flow {
				mgr := NewManager()
				blk := &passthroughBlock{}
				f, err := mgr.StartFlow(t.Context(), Config{
					PipelineID: "trans-2",
					Blocks:     []Block{blk},
				})
				if err != nil {
					t.Fatal(err)
				}
				if err := mgr.StopFlow(t.Context(), "trans-2"); err != nil {
					t.Fatal(err)
				}
				return f
			},
			wantFinal: FlowStatusStopped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.initialise(t)
			if f == nil {
				t.Fatal("expected non-nil flow")
			}
			if got := f.Status(); got != tt.wantFinal {
				t.Errorf("final status = %v, want %v", got, tt.wantFinal)
			}
		})
	}

	t.Run("construction failure leaves no flow", func(t *testing.T) {
		mgr := NewManager()
		f, err := mgr.StartFlow(t.Context(), Config{
			PipelineID: "trans-3",
			Blocks:     []Block{nil},
		})
		if err == nil {
			t.Fatal("expected error for nil block")
		}
		if f != nil {
			t.Fatal("expected nil flow on construction failure")
		}
	})
}

func TestFlow_TimestampsSet(t *testing.T) {
	mgr := NewManager()
	blk := &passthroughBlock{}
	f, err := mgr.StartFlow(t.Context(), Config{
		PipelineID: "ts-test",
		Blocks:     []Block{blk},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(t.Context(), "ts-test") })

	if f.CreatedAt().IsZero() {
		t.Error("CreatedAt is zero")
	}
	if f.StartedAt().IsZero() {
		t.Error("StartedAt is zero")
	}
	if !f.StoppedAt().IsZero() {
		t.Error("StoppedAt should be zero while still running")
	}
}

func TestFlow_StatusAfterStop(t *testing.T) {
	mgr := NewManager()
	blk := &passthroughBlock{}
	f, err := mgr.StartFlow(t.Context(), Config{
		PipelineID: "stop-test",
		Blocks:     []Block{blk},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := f.Status(); got != FlowStatusRunning {
		t.Fatalf("pre-stop status = %v, want running", got)
	}

	if err := mgr.StopFlow(t.Context(), "stop-test"); err != nil {
		t.Fatal(err)
	}

	if got := f.Status(); got != FlowStatusStopped {
		t.Fatalf("post-stop status = %v, want stopped", got)
	}
	if f.StoppedAt().IsZero() {
		t.Error("StoppedAt is zero after stop")
	}
}

func TestFlow_ConstructionFailureNoMapEntry(t *testing.T) {
	mgr := NewManager()
	f, err := mgr.StartFlow(t.Context(), Config{
		PipelineID: "fail-ts",
		Blocks:     []Block{nil},
	})
	if err == nil {
		t.Fatal("expected construction error")
	}
	if f != nil {
		t.Fatal("expected nil flow on construction failure")
	}
	if _, err := mgr.GetFlowStatus("fail-ts"); err == nil {
		t.Fatal("expected ErrFlowNotFound for unregistered instance")
	}
}

func TestFlow_PipelineID(t *testing.T) {
	mgr := NewManager()
	blk := &passthroughBlock{}
	f, err := mgr.StartFlow(t.Context(), Config{
		PipelineID: "pid-test",
		Blocks:     []Block{blk},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(t.Context(), "pid-test") })

	if got := f.PipelineID(); got != "pid-test" {
		t.Errorf("PipelineID = %q, want %q", got, "pid-test")
	}
}

func TestFlow_IDNotEmpty(t *testing.T) {
	mgr := NewManager()
	blk := &passthroughBlock{}
	f, err := mgr.StartFlow(t.Context(), Config{
		PipelineID: "id-test",
		Blocks:     []Block{blk},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(t.Context(), "id-test") })

	if got := f.ID(); got == "" {
		t.Error("ID is empty")
	}
}
