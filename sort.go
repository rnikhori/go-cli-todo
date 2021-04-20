package main

import (
	"fmt"
	"sort"
)

const (
	sortByPriority = "priority"
	sortByCreated  = "created"
	sortByDue      = "due"
)

func sortTodos(todos []Todo, by string) error {
	switch by {
	case "":
		return nil
	case sortByPriority:
		sort.SliceStable(todos, func(i, j int) bool {
			return todos[i].EffectivePriority().rank() > todos[j].EffectivePriority().rank()
		})
		return nil
	case sortByCreated:
		sort.SliceStable(todos, func(i, j int) bool {
			return todos[i].CreatedAt.Before(todos[j].CreatedAt)
		})
		return nil
	case sortByDue:
		sort.SliceStable(todos, func(i, j int) bool {
			a, b := todos[i].Due, todos[j].Due
			if a == nil {
				return false
			}
			if b == nil {
				return true
			}
			return a.Before(*b)
		})
		return nil
	default:
		return fmt.Errorf("invalid sort key %q (want priority, created, or due)", by)
	}
}
