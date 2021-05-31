package main

import (
	"fmt"
	"strconv"
)

func parseIDs(args []string) ([]int, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("requires at least one id")
	}
	ids := make([]int, 0, len(args))
	for _, a := range args {
		id, err := strconv.Atoi(a)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", a)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func removeByIDs(todos []Todo, ids []int) []Todo {
	remove := make(map[int]bool, len(ids))
	for _, id := range ids {
		remove[id] = true
	}
	out := todos[:0:0]
	for _, t := range todos {
		if !remove[t.ID] {
			out = append(out, t)
		}
	}
	return out
}
