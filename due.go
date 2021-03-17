package main

import (
	"fmt"
	"time"
)

const dueDateLayout = "2006-01-02"

func parseDue(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(dueDateLayout, s)
	if err != nil {
		return nil, fmt.Errorf("invalid due date %q (want YYYY-MM-DD): %w", s, err)
	}
	return &t, nil
}

func (t Todo) IsOverdue(now time.Time) bool {
	if t.Due == nil || t.Done {
		return false
	}
	return t.Due.Before(now)
}
