package main

import (
	"fmt"
	"sort"
)

const (
	sortByPriority = "priority"
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
	default:
		return fmt.Errorf("invalid sort key %q (want priority)", by)
	}
}
