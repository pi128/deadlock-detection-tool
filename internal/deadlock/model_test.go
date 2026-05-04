package deadlock

import (
	"strings"
	"testing"
)

func TestDetectsCycleDeadlock(t *testing.T) {
	input := `{
		"resources": {"R1": 1, "R2": 1},
		"allocations": {"P1:R1": 1, "P2:R2": 1},
		"requests": {"P1:R2": 1, "P2:R1": 1}
	}`
	system, err := Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	result := system.Detect()
	if !result.Deadlocked {
		t.Fatal("expected deadlock")
	}
	if len(result.Cycles) != 1 {
		t.Fatalf("expected 1 cycle, got %d", len(result.Cycles))
	}
	if got := formatCycle(result.Cycles[0]); got != "P1 -> P2 -> P1" {
		t.Fatalf("unexpected cycle %q", got)
	}
	if result.Safe {
		t.Fatal("deadlocked state should not be safe")
	}
}

func TestSafeScenarioHasSequence(t *testing.T) {
	input := `{
		"resources": {"R1": 2, "R2": 1},
		"allocations": {"P1:R1": 1, "P2:R2": 1},
		"requests": {"P1:R2": 1, "P2:R1": 1}
	}`
	system, err := Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	result := system.Detect()
	if result.Deadlocked {
		t.Fatalf("did not expect deadlock, cycles: %#v", result.Cycles)
	}
	if !result.Safe {
		t.Fatalf("expected safe state, blocked: %#v", result.Blocked)
	}
	if got := join(result.SafeSequence); got != "P2, P1" {
		t.Fatalf("unexpected safe sequence %q", got)
	}
}

func TestRejectsOverAllocation(t *testing.T) {
	input := `{
		"resources": {"R1": 1},
		"allocations": {"P1:R1": 2},
		"requests": {}
	}`
	_, err := Load(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected validation error")
	}
}
