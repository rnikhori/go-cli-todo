package main

import (
	"testing"
	"time"
)

func TestParseDue(t *testing.T) {
	t.Run("empty returns nil", func(t *testing.T) {
		got, err := parseDue("")
		if err != nil {
			t.Fatalf("parseDue(\"\") returned error: %v", err)
		}
		if got != nil {
			t.Errorf("parseDue(\"\") = %v, want nil", got)
		}
	})

	t.Run("valid date", func(t *testing.T) {
		got, err := parseDue("2026-01-15")
		if err != nil {
			t.Fatalf("parseDue returned error: %v", err)
		}
		want := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		if got == nil || !got.Equal(want) {
			t.Errorf("parseDue = %v, want %v", got, want)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		if _, err := parseDue("01/15/2026"); err == nil {
			t.Fatal("expected error for invalid date format")
		}
	})

	t.Run("garbage input", func(t *testing.T) {
		if _, err := parseDue("not-a-date"); err == nil {
			t.Fatal("expected error for garbage input")
		}
	})
}

func TestTodoIsOverdue(t *testing.T) {
	now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	past := now.AddDate(0, 0, -1)
	future := now.AddDate(0, 0, 1)

	cases := []struct {
		name string
		todo Todo
		want bool
	}{
		{"no due date", Todo{}, false},
		{"due in future", Todo{Due: &future}, false},
		{"due in past, not done", Todo{Due: &past}, true},
		{"due in past, done", Todo{Due: &past, Done: true}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.todo.IsOverdue(now); got != c.want {
				t.Errorf("IsOverdue() = %v, want %v", got, c.want)
			}
		})
	}
}
