package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
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

func importCSV(path string) ([]Todo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	if !reflect.DeepEqual(records[0], csvHeader) {
		return nil, fmt.Errorf("unexpected CSV header: %v", records[0])
	}

	var todos []Todo
	for _, rec := range records[1:] {
		if len(rec) != len(csvHeader) {
			return nil, fmt.Errorf("malformed CSV row: %v", rec)
		}
		done, err := strconv.ParseBool(rec[2])
		if err != nil {
			return nil, fmt.Errorf("invalid done value in CSV: %s", rec[2])
		}
		priority, err := ParsePriority(rec[3])
		if err != nil {
			return nil, err
		}
		var due *time.Time
		if rec[4] != "" {
			due, err = parseDue(rec[4])
			if err != nil {
				return nil, err
			}
		}
		var tags []string
		if rec[5] != "" {
			tags = normalizeTags(strings.Split(rec[5], ";"))
		}
		var created time.Time
		if rec[6] != "" {
			created, err = time.Parse(time.RFC3339, rec[6])
			if err != nil {
				return nil, fmt.Errorf("invalid created_at in CSV: %s", rec[6])
			}
		}
		todos = append(todos, Todo{
			Text:      rec[1],
			Done:      done,
			Priority:  priority,
			Due:       due,
			Tags:      tags,
			CreatedAt: created,
		})
	}
	return todos, nil
}
