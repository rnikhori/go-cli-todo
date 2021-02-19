package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Todo struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Done     bool     `json:"done"`
	Priority Priority `json:"priority,omitempty"`
}

func defaultStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".go-cli-todo.json"), nil
}

const storeEnvVar = "TODO_STORE"

func resolveStorePath(flagVal string) (string, error) {
	if flagVal != "" {
		return flagVal, nil
	}
	if env := os.Getenv(storeEnvVar); env != "" {
		return env, nil
	}
	return defaultStorePath()
}

func loadTodos(path string) ([]Todo, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Todo{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return []Todo{}, nil
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func saveTodos(path string, todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func findByID(todos []Todo, id int) int {
	for i := range todos {
		if todos[i].ID == id {
			return i
		}
	}
	return -1
}

func nextID(todos []Todo) int {
	max := 0
	for _, t := range todos {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}
