package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "todos.json")
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	return buf.String()
}

func TestDispatchAddCreatesTodo(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"buy", "milk"}, path); err != nil {
		t.Fatalf("dispatch add returned error: %v", err)
	}
	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("got %d todos, want 1", len(todos))
	}
	if todos[0].Text != "buy milk" {
		t.Errorf("text = %q, want %q", todos[0].Text, "buy milk")
	}
	if todos[0].EffectivePriority() != PriorityMedium {
		t.Errorf("priority = %q, want %q", todos[0].Priority, PriorityMedium)
	}
}

func TestDispatchAddRequiresText(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", nil, path); err == nil {
		t.Fatal("dispatch add with no text should return an error")
	}
}

func TestDispatchEditChangesText(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"old", "text"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("edit", []string{"1", "new", "text"}, path); err != nil {
		t.Fatalf("dispatch edit returned error: %v", err)
	}
	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if todos[0].Text != "new text" {
		t.Errorf("text = %q, want %q", todos[0].Text, "new text")
	}
}

func TestDispatchEditUnknownID(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("edit", []string{"99", "text"}, path); err == nil {
		t.Fatal("dispatch edit with unknown id should return an error")
	}
}

func TestDispatchDoneMarksComplete(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("done", []string{"1"}, path); err != nil {
		t.Fatalf("dispatch done returned error: %v", err)
	}
	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if !todos[0].Done {
		t.Errorf("expected todo 1 to be done")
	}
}

func TestDispatchRmRemovesTodo(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("rm", []string{"1"}, path); err != nil {
		t.Fatalf("dispatch rm returned error: %v", err)
	}
	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("got %d todos after rm, want 0", len(todos))
	}
}

func TestDispatchListTagFilter(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"--tag", "work", "task one"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("add", []string{"--tag", "home", "task two"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := dispatch("list", []string{"--tag", "work"}, path); err != nil {
			t.Fatalf("dispatch list returned error: %v", err)
		}
	})

	if !strings.Contains(out, "task one") {
		t.Errorf("output %q should contain %q", out, "task one")
	}
	if strings.Contains(out, "task two") {
		t.Errorf("output %q should not contain %q", out, "task two")
	}
}

func TestDispatchUnknownCommand(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("bogus", nil, path); err == nil {
		t.Fatal("dispatch with unknown command should return an error")
	}
}
