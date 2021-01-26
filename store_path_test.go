package main

import (
	"os"
	"testing"
)

func TestResolveStorePath(t *testing.T) {
	t.Run("flag takes precedence", func(t *testing.T) {
		t.Setenv(storeEnvVar, "/env/path.json")
		got, err := resolveStorePath("/flag/path.json")
		if err != nil {
			t.Fatalf("resolveStorePath returned error: %v", err)
		}
		if got != "/flag/path.json" {
			t.Errorf("got %q, want /flag/path.json", got)
		}
	})

	t.Run("env var used when no flag", func(t *testing.T) {
		t.Setenv(storeEnvVar, "/env/path.json")
		got, err := resolveStorePath("")
		if err != nil {
			t.Fatalf("resolveStorePath returned error: %v", err)
		}
		if got != "/env/path.json" {
			t.Errorf("got %q, want /env/path.json", got)
		}
	})

	t.Run("falls back to default when neither set", func(t *testing.T) {
		os.Unsetenv(storeEnvVar)
		got, err := resolveStorePath("")
		if err != nil {
			t.Fatalf("resolveStorePath returned error: %v", err)
		}
		want, err := defaultStorePath()
		if err != nil {
			t.Fatalf("defaultStorePath returned error: %v", err)
		}
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
