package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: todo [--store path] <command> [arguments]

global options:
  --store path   path to the todo storage file (overrides TODO_STORE env var)

commands:
  add <text>   add a new todo
  list         list all todos
  done <id>    mark a todo as done
  rm <id>      remove a todo`)
}

func run() error {
	fs := flag.NewFlagSet("todo", flag.ContinueOnError)
	fs.Usage = usage
	storeFlag := fs.String("store", "", "path to the todo storage file")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) < 1 {
		usage()
		return fmt.Errorf("no command given")
	}

	path, err := resolveStorePath(*storeFlag)
	if err != nil {
		return err
	}

	cmd := rest[0]
	args := rest[1:]

	switch cmd {
	case "add":
		if len(args) == 0 {
			return fmt.Errorf("add requires text")
		}
		text := strings.Join(args, " ")
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		todo := Todo{ID: nextID(todos), Text: text}
		todos = append(todos, todo)
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("added todo %d: %s\n", todo.ID, todo.Text)

	case "list":
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		if len(todos) == 0 {
			fmt.Println("no todos")
			return nil
		}
		for _, t := range todos {
			mark := " "
			if t.Done {
				mark = "x"
			}
			fmt.Printf("[%s] %d: %s\n", mark, t.ID, t.Text)
		}

	case "done":
		if len(args) != 1 {
			return fmt.Errorf("done requires exactly one id")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		found := false
		for i := range todos {
			if todos[i].ID == id {
				todos[i].Done = true
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no todo with id %d", id)
		}
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("marked todo %d done\n", id)

	case "rm":
		if len(args) != 1 {
			return fmt.Errorf("rm requires exactly one id")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		idx := -1
		for i := range todos {
			if todos[i].ID == id {
				idx = i
				break
			}
		}
		if idx == -1 {
			return fmt.Errorf("no todo with id %d", id)
		}
		todos = append(todos[:idx], todos[idx+1:]...)
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("removed todo %d\n", id)

	default:
		usage()
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
