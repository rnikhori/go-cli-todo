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
  add [--priority low|medium|high] <text>   add a new todo
  list                                      list all todos
  edit <id> <text>                          change a todo's text
  done <id>                                 mark a todo as done
  rm <id>                                   remove a todo`)
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
		addFs := flag.NewFlagSet("add", flag.ContinueOnError)
		priorityStr := addFs.String("priority", "medium", "priority: low, medium, or high")
		if err := addFs.Parse(args); err != nil {
			return err
		}
		textArgs := addFs.Args()
		if len(textArgs) == 0 {
			return fmt.Errorf("add requires text")
		}
		priority, err := ParsePriority(*priorityStr)
		if err != nil {
			return err
		}
		text := strings.Join(textArgs, " ")
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		todo := Todo{ID: nextID(todos), Text: text, Priority: priority}
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
			fmt.Printf("[%s] %d (%s): %s\n", mark, t.ID, t.EffectivePriority(), t.Text)
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
		idx := findByID(todos, id)
		if idx == -1 {
			return fmt.Errorf("no todo with id %d", id)
		}
		todos[idx].Done = true
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("marked todo %d done\n", id)

	case "edit":
		if len(args) < 2 {
			return fmt.Errorf("edit requires an id and new text")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		text := strings.Join(args[1:], " ")
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		idx := findByID(todos, id)
		if idx == -1 {
			return fmt.Errorf("no todo with id %d", id)
		}
		todos[idx].Text = text
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("updated todo %d: %s\n", id, text)

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
		idx := findByID(todos, id)
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
