package main

import "time"

type Stats struct {
	Total      int
	Done       int
	Pending    int
	Overdue    int
	ByPriority map[Priority]int
}

func computeStats(todos []Todo, now time.Time) Stats {
	s := Stats{ByPriority: map[Priority]int{
		PriorityLow:    0,
		PriorityMedium: 0,
		PriorityHigh:   0,
	}}
	for _, t := range todos {
		s.Total++
		if t.Done {
			s.Done++
		} else {
			s.Pending++
		}
		s.ByPriority[t.EffectivePriority()]++
		if t.IsOverdue(now) {
			s.Overdue++
		}
	}
	return s
}
