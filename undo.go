package main

import "os"

func undoPath(path string) string {
	return path + ".undo"
}

var mutatingCommands = map[string]bool{
	"add":             true,
	"done":            true,
	"edit":            true,
	"rm":              true,
	"clear-completed": true,
	"import":          true,
}

func saveUndoSnapshot(path string) error {
	todos, err := loadTodos(path)
	if err != nil {
		return err
	}
	return saveTodos(undoPath(path), todos)
}

func hasUndoSnapshot(path string) bool {
	_, err := os.Stat(undoPath(path))
	return err == nil
}
