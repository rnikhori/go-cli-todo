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
	if err := dispatch("rm", []string{"--force", "1"}, path); err != nil {
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

func TestDispatchListSearchFilter(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"buy", "milk"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("add", []string{"write", "report"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := dispatch("list", []string{"--search", "MILK"}, path); err != nil {
			t.Fatalf("dispatch list returned error: %v", err)
		}
	})

	if !strings.Contains(out, "buy milk") {
		t.Errorf("output %q should contain %q", out, "buy milk")
	}
	if strings.Contains(out, "write report") {
		t.Errorf("output %q should not contain %q", out, "write report")
	}
}

func TestDispatchListSearchNoMatches(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"buy", "milk"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := dispatch("list", []string{"--search", "nothing-matches"}, path); err != nil {
			t.Fatalf("dispatch list returned error: %v", err)
		}
	})

	if !strings.Contains(out, "no todos") {
		t.Errorf("output %q should report no todos", out)
	}
}

func TestDispatchRmPromptsAndAborts(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	withStdin(t, "n\n", func() {
		captureStdout(t, func() {
			if err := dispatch("rm", []string{"1"}, path); err != nil {
				t.Fatalf("dispatch rm returned error: %v", err)
			}
		})
	})

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected todo to survive an aborted rm, got %d todos", len(todos))
	}
}

func TestDispatchRmPromptsAndConfirms(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	withStdin(t, "y\n", func() {
		captureStdout(t, func() {
			if err := dispatch("rm", []string{"1"}, path); err != nil {
				t.Fatalf("dispatch rm returned error: %v", err)
			}
		})
	})

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected todo to be removed after confirming, got %d todos", len(todos))
	}
}

func TestDispatchDoneBulkMarksAllComplete(t *testing.T) {
	path := tempStore(t)
	for _, text := range []string{"task one", "task two", "task three"} {
		if err := dispatch("add", []string{text}, path); err != nil {
			t.Fatalf("setup add failed: %v", err)
		}
	}

	if err := dispatch("done", []string{"1", "3"}, path); err != nil {
		t.Fatalf("dispatch done returned error: %v", err)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if !todos[0].Done || todos[1].Done || !todos[2].Done {
		t.Errorf("unexpected done states: %+v", todos)
	}
}

func TestDispatchDoneUnknownIDLeavesStoreUnchanged(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("done", []string{"1", "99"}, path); err == nil {
		t.Fatal("expected error for unknown id in bulk done")
	}
}

func TestDispatchRmBulkWithForceRemovesAll(t *testing.T) {
	path := tempStore(t)
	for _, text := range []string{"task one", "task two", "task three"} {
		if err := dispatch("add", []string{text}, path); err != nil {
			t.Fatalf("setup add failed: %v", err)
		}
	}

	if err := dispatch("rm", []string{"--force", "1", "3"}, path); err != nil {
		t.Fatalf("dispatch rm returned error: %v", err)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 || todos[0].ID != 2 {
		t.Errorf("expected only todo 2 to remain, got %+v", todos)
	}
}

func TestDispatchRmBulkConfirmationShowsAllIDs(t *testing.T) {
	path := tempStore(t)
	for _, text := range []string{"task one", "task two"} {
		if err := dispatch("add", []string{text}, path); err != nil {
			t.Fatalf("setup add failed: %v", err)
		}
	}

	var out string
	withStdin(t, "y\n", func() {
		out = captureStdout(t, func() {
			if err := dispatch("rm", []string{"1", "2"}, path); err != nil {
				t.Fatalf("dispatch rm returned error: %v", err)
			}
		})
	})

	if !strings.Contains(out, "task one") || !strings.Contains(out, "task two") {
		t.Errorf("expected confirmation to list both todos, got %q", out)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected both todos removed, got %+v", todos)
	}
}

func TestDispatchClearCompletedWithForce(t *testing.T) {
	path := tempStore(t)
	for _, text := range []string{"done task", "pending task"} {
		if err := dispatch("add", []string{text}, path); err != nil {
			t.Fatalf("setup add failed: %v", err)
		}
	}
	if err := dispatch("done", []string{"1"}, path); err != nil {
		t.Fatalf("setup done failed: %v", err)
	}

	if err := dispatch("clear-completed", []string{"--force"}, path); err != nil {
		t.Fatalf("dispatch clear-completed returned error: %v", err)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 || todos[0].ID != 2 {
		t.Errorf("expected only the pending todo to remain, got %+v", todos)
	}
}

func TestDispatchClearCompletedAbort(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"done task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("done", []string{"1"}, path); err != nil {
		t.Fatalf("setup done failed: %v", err)
	}

	withStdin(t, "n\n", func() {
		captureStdout(t, func() {
			if err := dispatch("clear-completed", nil, path); err != nil {
				t.Fatalf("dispatch clear-completed returned error: %v", err)
			}
		})
	})

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected todo to survive an aborted clear, got %d todos", len(todos))
	}
}

func TestDispatchClearCompletedNoneCompleted(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"pending task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	out := captureStdout(t, func() {
		if err := dispatch("clear-completed", nil, path); err != nil {
			t.Fatalf("dispatch clear-completed returned error: %v", err)
		}
	})
	if !strings.Contains(out, "no completed") {
		t.Errorf("expected 'no completed' message, got %q", out)
	}
}

func TestDispatchUndoRestoresAfterRm(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task one"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("add", []string{"task two"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("rm", []string{"--force", "1"}, path); err != nil {
		t.Fatalf("setup rm failed: %v", err)
	}

	if err := dispatch("undo", nil, path); err != nil {
		t.Fatalf("dispatch undo returned error: %v", err)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos after undo, got %d", len(todos))
	}
}

func TestDispatchUndoIsSingleStep(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("rm", []string{"--force", "1"}, path); err != nil {
		t.Fatalf("setup rm failed: %v", err)
	}

	if err := dispatch("undo", nil, path); err != nil {
		t.Fatalf("first undo returned error: %v", err)
	}
	if err := dispatch("undo", nil, path); err == nil {
		t.Fatal("expected second consecutive undo to fail with nothing to undo")
	}
}

func TestDispatchUndoWithNoHistory(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("undo", nil, path); err == nil {
		t.Fatal("expected undo with no history to return an error")
	}
}

func TestDispatchUndoRestoresAfterClearCompleted(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}
	if err := dispatch("done", []string{"1"}, path); err != nil {
		t.Fatalf("setup done failed: %v", err)
	}
	if err := dispatch("clear-completed", []string{"--force"}, path); err != nil {
		t.Fatalf("setup clear-completed failed: %v", err)
	}

	if err := dispatch("undo", nil, path); err != nil {
		t.Fatalf("dispatch undo returned error: %v", err)
	}

	todos, err := loadTodos(path)
	if err != nil {
		t.Fatalf("loadTodos returned error: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo restored after undo, got %d", len(todos))
	}
}

func TestDispatchErrorMessagesAreCommandPrefixed(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("add", []string{"task"}, path); err != nil {
		t.Fatalf("setup add failed: %v", err)
	}

	cases := []struct {
		name   string
		cmd    string
		args   []string
		prefix string
	}{
		{"done unknown id", "done", []string{"99"}, "done:"},
		{"edit unknown id", "edit", []string{"99", "text"}, "edit:"},
		{"edit invalid id", "edit", []string{"abc", "text"}, "edit:"},
		{"rm unknown id", "rm", []string{"--force", "99"}, "rm:"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := dispatch(c.cmd, c.args, path)
			if err == nil {
				t.Fatalf("expected an error for %s %v", c.cmd, c.args)
			}
			if !strings.HasPrefix(err.Error(), c.prefix) {
				t.Errorf("error %q does not start with %q", err.Error(), c.prefix)
			}
		})
	}
}

func TestDispatchVersionPrintsVersion(t *testing.T) {
	path := tempStore(t)
	out := captureStdout(t, func() {
		if err := dispatch("version", nil, path); err != nil {
			t.Fatalf("dispatch version returned error: %v", err)
		}
	})
	if strings.TrimSpace(out) != version {
		t.Errorf("version output = %q, want %q", out, version)
	}
}

func TestDispatchUnknownCommand(t *testing.T) {
	path := tempStore(t)
	if err := dispatch("bogus", nil, path); err == nil {
		t.Fatal("dispatch with unknown command should return an error")
	}
}
