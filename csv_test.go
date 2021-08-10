package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExportCSV(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "out.csv")

	due := mustDue(t, "2026-01-15")
	todos := []Todo{
		{ID: 1, Text: "buy milk", Done: false, Priority: PriorityHigh, Due: due, Tags: []string{"home", "errand"}, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Text: "write report", Done: true},
	}

	if err := exportCSV(csvPath, todos); err != nil {
		t.Fatalf("exportCSV returned error: %v", err)
	}

	data, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("failed to read exported CSV: %v", err)
	}
	content := string(data)

	if !strings.HasPrefix(content, "id,text,done,priority,due,tags,created_at") {
		t.Errorf("expected header row, got %q", content)
	}
	if !strings.Contains(content, "buy milk") {
		t.Errorf("expected exported text, got %q", content)
	}
	if !strings.Contains(content, "home;errand") {
		t.Errorf("expected semicolon-joined tags, got %q", content)
	}
	if !strings.Contains(content, "2026-01-15") {
		t.Errorf("expected due date, got %q", content)
	}
	if !strings.Contains(content, "write report") || !strings.Contains(content, "true") {
		t.Errorf("expected second row with done=true, got %q", content)
	}
}

func TestExportCSVEmpty(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "empty.csv")

	if err := exportCSV(csvPath, nil); err != nil {
		t.Fatalf("exportCSV returned error: %v", err)
	}

	data, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("failed to read exported CSV: %v", err)
	}
	if strings.TrimSpace(string(data)) != strings.Join(csvHeader, ",") {
		t.Errorf("expected only header row for empty input, got %q", data)
	}
}
