package main

import (
	"testing"
	"time"
)

func mustDue(t *testing.T, s string) *time.Time {
	t.Helper()
	d, err := parseDue(s)
	if err != nil {
		t.Fatalf("parseDue(%q) failed: %v", s, err)
	}
	return d
}

func TestSortTodosByPriority(t *testing.T) {
	todos := []Todo{
		{ID: 1, Priority: PriorityLow},
		{ID: 2, Priority: PriorityHigh},
		{ID: 3, Priority: PriorityMedium},
	}
	if err := sortTodos(todos, sortByPriority); err != nil {
		t.Fatalf("sortTodos returned error: %v", err)
	}
	wantOrder := []int{2, 3, 1}
	for i, id := range wantOrder {
		if todos[i].ID != id {
			t.Errorf("position %d = id %d, want id %d", i, todos[i].ID, id)
		}
	}
}

func TestSortTodosByCreated(t *testing.T) {
	now := time.Now()
	todos := []Todo{
		{ID: 1, CreatedAt: now.Add(2 * time.Hour)},
		{ID: 2, CreatedAt: now},
		{ID: 3, CreatedAt: now.Add(time.Hour)},
	}
	if err := sortTodos(todos, sortByCreated); err != nil {
		t.Fatalf("sortTodos returned error: %v", err)
	}
	wantOrder := []int{2, 3, 1}
	for i, id := range wantOrder {
		if todos[i].ID != id {
			t.Errorf("position %d = id %d, want id %d", i, todos[i].ID, id)
		}
	}
}

func TestSortTodosByDue(t *testing.T) {
	todos := []Todo{
		{ID: 1, Due: mustDue(t, "2026-03-01")},
		{ID: 2, Due: nil},
		{ID: 3, Due: mustDue(t, "2026-01-01")},
	}
	if err := sortTodos(todos, sortByDue); err != nil {
		t.Fatalf("sortTodos returned error: %v", err)
	}
	wantOrder := []int{3, 1, 2}
	for i, id := range wantOrder {
		if todos[i].ID != id {
			t.Errorf("position %d = id %d, want id %d", i, todos[i].ID, id)
		}
	}
}

func TestSortTodosEmptyKeyIsNoop(t *testing.T) {
	todos := []Todo{{ID: 2}, {ID: 1}}
	if err := sortTodos(todos, ""); err != nil {
		t.Fatalf("sortTodos returned error: %v", err)
	}
	if todos[0].ID != 2 || todos[1].ID != 1 {
		t.Errorf("expected order unchanged, got %+v", todos)
	}
}

func TestSortTodosInvalidKey(t *testing.T) {
	todos := []Todo{{ID: 1}}
	if err := sortTodos(todos, "bogus"); err == nil {
		t.Fatal("expected error for invalid sort key")
	}
}
