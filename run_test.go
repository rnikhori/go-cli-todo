package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withArgs(t *testing.T, args []string, fn func()) {
	t.Helper()
	old := os.Args
	os.Args = append([]string{"todo"}, args...)
	defer func() { os.Args = old }()
	fn()
}

func TestRunVersionFlag(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		withArgs(t, []string{"--version"}, func() {
			err = run()
		})
	})
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
	if strings.TrimSpace(out) != version {
		t.Errorf("output = %q, want %q", out, version)
	}
}

func TestRunNoCommand(t *testing.T) {
	var err error
	captureStdout(t, func() {
		withArgs(t, nil, func() {
			err = run()
		})
	})
	if exitCodeFor(err) != exitUsage {
		t.Errorf("exitCodeFor(run()) = %d, want %d (err: %v)", exitCodeFor(err), exitUsage, err)
	}
}

func TestRunEndToEndWithStoreFlag(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")

	captureStdout(t, func() {
		withArgs(t, []string{"--store", path, "add", "buy milk"}, func() {
			if err := run(); err != nil {
				t.Fatalf("run() add returned error: %v", err)
			}
		})
	})

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 || todos[0].Text != "buy milk" {
		t.Errorf("unexpected todos after run(): %+v", todos)
	}
}

func TestRunUnknownCommandExitUsage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")
	var err error
	captureStdout(t, func() {
		withArgs(t, []string{"--store", path, "bogus"}, func() {
			err = run()
		})
	})
	if exitCodeFor(err) != exitUsage {
		t.Errorf("exitCodeFor(run()) = %d, want %d (err: %v)", exitCodeFor(err), exitUsage, err)
	}
}
