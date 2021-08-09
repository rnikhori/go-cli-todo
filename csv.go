package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

var csvHeader = []string{"id", "text", "done", "priority", "due", "tags", "created_at"}

func exportCSV(path string, todos []Todo) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(csvHeader); err != nil {
		return err
	}
	for _, t := range todos {
		due := ""
		if t.Due != nil {
			due = t.Due.Format(dueDateLayout)
		}
		created := ""
		if !t.CreatedAt.IsZero() {
			created = t.CreatedAt.Format(time.RFC3339)
		}
		record := []string{
			strconv.Itoa(t.ID),
			t.Text,
			strconv.FormatBool(t.Done),
			string(t.EffectivePriority()),
			due,
			strings.Join(t.Tags, ";"),
			created,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
