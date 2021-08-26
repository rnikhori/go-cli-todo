package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportCSVRoundTrip(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "todos.csv")

	due := mustDue(t, "2026-02-01")
	original := []Todo{
		{ID: 1, Text: "buy milk", Done: false, Priority: PriorityHigh, Due: due, Tags: []string{"home", "errand"}},
		{ID: 2, Text: "write report", Done: true, Priority: PriorityLow},
	}

	if err := exportCSV(csvPath, original); err != nil {
		t.Fatalf("exportCSV returned error: %v", err)
	}

	imported, err := importCSV(csvPath)
	if err != nil {
		t.Fatalf("importCSV returned error: %v", err)
	}
	if len(imported) != len(original) {
		t.Fatalf("imported %d todos, want %d", len(imported), len(original))
	}
	if imported[0].Text != "buy milk" || imported[0].Priority != PriorityHigh {
		t.Errorf("imported[0] = %+v, mismatch with original", imported[0])
	}
	if imported[0].Due == nil || !imported[0].Due.Equal(*due) {
		t.Errorf("imported[0].Due = %v, want %v", imported[0].Due, due)
	}
	if len(imported[0].Tags) != 2 {
		t.Errorf("imported[0].Tags = %v, want 2 tags", imported[0].Tags)
	}
	if !imported[1].Done || imported[1].Priority != PriorityLow {
		t.Errorf("imported[1] = %+v, mismatch with original", imported[1])
	}
}

func TestImportCSVRejectsBadHeader(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "bad.csv")
	writeFile(t, csvPath, "not,the,right,header\n1,2,3,4\n")

	if _, err := importCSV(csvPath); err == nil {
		t.Fatal("expected error for CSV with unexpected header")
	}
}

func TestImportCSVRejectsMalformedRow(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "malformed.csv")
	writeFile(t, csvPath, "id,text,done,priority,due,tags,created_at\n1,bad row,notabool,high,,,\n")

	if _, err := importCSV(csvPath); err == nil {
		t.Fatal("expected error for row with invalid done value")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
}
