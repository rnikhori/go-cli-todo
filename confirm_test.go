package main

import (
	"strings"
	"testing"
)

func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()
	old := stdin
	stdin = strings.NewReader(input)
	defer func() { stdin = old }()
	fn()
}

func TestConfirm(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"y\n", true},
		{"Y\n", true},
		{"yes\n", true},
		{"YES\n", true},
		{"n\n", false},
		{"no\n", false},
		{"\n", false},
		{"", false},
		{"maybe\n", false},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			var got bool
			withStdin(t, c.input, func() {
				out := captureStdout(t, func() {
					var err error
					got, err = confirm("proceed?")
					if err != nil {
						t.Fatalf("confirm returned error: %v", err)
					}
				})
				if !strings.Contains(out, "proceed?") {
					t.Errorf("expected prompt to be printed, got %q", out)
				}
			})
			if got != c.want {
				t.Errorf("confirm with input %q = %v, want %v", c.input, got, c.want)
			}
		})
	}
}
