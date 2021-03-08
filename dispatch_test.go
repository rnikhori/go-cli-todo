package main

import (
	"path/filepath"
	"testing"
)

func tempStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "todos.json")
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

func TestDispatchUnknownCommand(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("bogus", nil, path); err == nil {
		t.Fatal("dispatch with unknown command should return an error")
	}
}
