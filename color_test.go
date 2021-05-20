package main

import "testing"

func TestColorEnabledNoColorFlag(t *testing.T) {
	if colorEnabled(true) {
		t.Error("colorEnabled(true) should always return false")
	}
}

func TestColorEnabledRespectsNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if colorEnabled(false) {
		t.Error("colorEnabled should return false when NO_COLOR is set")
	}
}

func TestColorize(t *testing.T) {
	if got := colorize(false, ansiRed, "text"); got != "text" {
		t.Errorf("colorize(false, ...) = %q, want %q", got, "text")
	}
	want := ansiRed + "text" + ansiReset
	if got := colorize(true, ansiRed, "text"); got != want {
		t.Errorf("colorize(true, ...) = %q, want %q", got, want)
	}
}

func TestPriorityColor(t *testing.T) {
	cases := map[Priority]string{
		PriorityHigh:   ansiRed,
		PriorityLow:    ansiGreen,
		PriorityMedium: ansiYellow,
	}
	for p, want := range cases {
		if got := priorityColor(p); got != want {
			t.Errorf("priorityColor(%q) = %q, want %q", p, got, want)
		}
	}
}
