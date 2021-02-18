package main

import "testing"

func TestParsePriority(t *testing.T) {
	cases := []struct {
		in      string
		want    Priority
		wantErr bool
	}{
		{"low", PriorityLow, false},
		{"LOW", PriorityLow, false},
		{"l", PriorityLow, false},
		{"medium", PriorityMedium, false},
		{"med", PriorityMedium, false},
		{"", PriorityMedium, false},
		{"high", PriorityHigh, false},
		{"H", PriorityHigh, false},
		{"urgent", "", true},
		{"123", "", true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := ParsePriority(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("ParsePriority(%q) = %v, want error", c.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePriority(%q) returned unexpected error: %v", c.in, err)
			}
			if got != c.want {
				t.Errorf("ParsePriority(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestPriorityRankOrdering(t *testing.T) {
	if PriorityLow.rank() >= PriorityMedium.rank() {
		t.Errorf("expected low rank < medium rank")
	}
	if PriorityMedium.rank() >= PriorityHigh.rank() {
		t.Errorf("expected medium rank < high rank")
	}
}

func TestEffectivePriorityDefaultsToMedium(t *testing.T) {
	todo := Todo{ID: 1, Text: "legacy todo"}
	if got := todo.EffectivePriority(); got != PriorityMedium {
		t.Errorf("EffectivePriority() = %q, want %q", got, PriorityMedium)
	}

	todo.Priority = PriorityHigh
	if got := todo.EffectivePriority(); got != PriorityHigh {
		t.Errorf("EffectivePriority() = %q, want %q", got, PriorityHigh)
	}
}
