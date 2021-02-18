package main

import (
	"fmt"
	"strings"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func ParsePriority(s string) (Priority, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low", "l":
		return PriorityLow, nil
	case "medium", "med", "m", "":
		return PriorityMedium, nil
	case "high", "h":
		return PriorityHigh, nil
	default:
		return "", fmt.Errorf("invalid priority %q (want low, medium, or high)", s)
	}
}

func (p Priority) rank() int {
	switch p {
	case PriorityHigh:
		return 2
	case PriorityLow:
		return 0
	default:
		return 1
	}
}

func (t Todo) EffectivePriority() Priority {
	if t.Priority == "" {
		return PriorityMedium
	}
	return t.Priority
}
