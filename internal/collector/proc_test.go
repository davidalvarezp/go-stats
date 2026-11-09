package collector

import (
	"testing"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	if c == nil {
		t.Fatal("Expected non-nil Collector instance")
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{-10.5, 0.0},
		{0.0, 0.0},
		{45.2, 45.2},
		{100.0, 100.0},
		{150.8, 100.0},
	}

	for _, tt := range tests {
		got := clamp(tt.input)
		if got != tt.expected {
			t.Errorf("clamp(%f) = %f; want %f", tt.input, got, tt.expected)
		}
	}
}

func TestParseCPULine(t *testing.T) {
	line := "cpu  2255 34 2290 22625563 6290 127 456 0 0 0"
	stat, err := parseCPULine(line)
	if err != nil {
		t.Fatalf("parseCPULine failed: %v", err)
	}

	if stat.User != 2255 {
		t.Errorf("stat.User = %d; want 2255", stat.User)
	}
	if stat.Idle != 22625563 {
		t.Errorf("stat.Idle = %d; want 22625563", stat.Idle)
	}
	if stat.Total <= stat.Busy {
		t.Errorf("stat.Total (%d) must be greater than stat.Busy (%d)", stat.Total, stat.Busy)
	}
}

func TestCollectSmoke(t *testing.T) {
	c := NewCollector()
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if snap == nil {
		t.Fatal("Expected non-nil SystemSnapshot")
	}
	if snap.Timestamp <= 0 {
		t.Errorf("Timestamp = %d; want > 0", snap.Timestamp)
	}
}
