package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNextID(t *testing.T) {
	cases := []struct {
		name  string
		todos []Todo
		want  int
	}{
		{"empty", nil, 1},
		{"single", []Todo{{ID: 1}}, 2},
		{"gap", []Todo{{ID: 1}, {ID: 5}}, 6},
		{"unordered", []Todo{{ID: 3}, {ID: 1}, {ID: 2}}, 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nextID(c.todos); got != c.want {
				t.Errorf("nextID(%v) = %d, want %d", c.todos, got, c.want)
			}
		})
	}
}

func TestSaveAndLoadTodos(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")

	want := []Todo{
		{ID: 1, Text: "buy milk", Done: false},
		{ID: 2, Text: "write report", Done: true},
	}

	if err := saveTodos(path, want); err != nil {
		t.Fatalf("saveTodos returned error: %v", err)
	}

	got, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("loadTodos returned %d todos, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("todo %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestLoadTodosMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.json")

	got, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error for missing file: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("loadTodos on missing file = %v, want empty slice", got)
	}
}

func TestLoadTodosEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty file: %v", err)
	}

	got, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error for empty file: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("loadTodos on empty file = %v, want empty slice", got)
	}
}
