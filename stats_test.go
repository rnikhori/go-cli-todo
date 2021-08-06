package main

import (
	"testing"
	"time"
)

func TestComputeStats(t *testing.T) {
	now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	past := now.AddDate(0, 0, -1)

	todos := []Todo{
		{ID: 1, Done: true, Priority: PriorityHigh},
		{ID: 2, Done: false, Priority: PriorityLow},
		{ID: 3, Done: false, Priority: PriorityMedium, Due: &past},
		{ID: 4, Done: false},
	}

	s := computeStats(todos, now)

	if s.Total != 4 {
		t.Errorf("Total = %d, want 4", s.Total)
	}
	if s.Done != 1 {
		t.Errorf("Done = %d, want 1", s.Done)
	}
	if s.Pending != 3 {
		t.Errorf("Pending = %d, want 3", s.Pending)
	}
	if s.Overdue != 1 {
		t.Errorf("Overdue = %d, want 1", s.Overdue)
	}
	if s.ByPriority[PriorityHigh] != 1 {
		t.Errorf("ByPriority[high] = %d, want 1", s.ByPriority[PriorityHigh])
	}
	if s.ByPriority[PriorityLow] != 1 {
		t.Errorf("ByPriority[low] = %d, want 1", s.ByPriority[PriorityLow])
	}
	if s.ByPriority[PriorityMedium] != 2 {
		t.Errorf("ByPriority[medium] = %d, want 2 (includes default)", s.ByPriority[PriorityMedium])
	}
}

func TestComputeStatsEmpty(t *testing.T) {
	s := computeStats(nil, time.Now())
	if s.Total != 0 || s.Done != 0 || s.Pending != 0 || s.Overdue != 0 {
		t.Errorf("expected all-zero stats for empty input, got %+v", s)
	}
}
