package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: todo [--store path] <command> [arguments]

global options:
  --store path   path to the todo storage file (overrides TODO_STORE env var)

commands:
  add [--priority low|medium|high] [--due YYYY-MM-DD] [--tag t]... <text>   add a new todo
  list [--tag t] [--search text] [--sort k] [--no-color]   list todos
  edit <id> <text>                          change a todo's text
  done <id> [id...]                         mark one or more todos as done
  rm [--force] <id> [id...]                 remove one or more todos (prompts unless --force)
  clear-completed [--force]                 remove all completed todos (prompts unless --force)
  undo                                      undo the last add, done, edit, rm, or clear-completed`)
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

	return dispatch(rest[0], rest[1:], path)
}

func dispatch(cmd string, args []string, path string) error {
	if mutatingCommands[cmd] {
		if err := saveUndoSnapshot(path); err != nil {
			return err
		}
	}

	switch cmd {
	case "add":
		addFs := flag.NewFlagSet("add", flag.ContinueOnError)
		priorityStr := addFs.String("priority", "medium", "priority: low, medium, or high")
		dueStr := addFs.String("due", "", "due date in YYYY-MM-DD format")
		var tags stringList
		addFs.Var(&tags, "tag", "tag to attach to the todo (repeatable)")
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
		due, err := parseDue(*dueStr)
		if err != nil {
			return err
		}
		text := strings.Join(textArgs, " ")
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		todo := Todo{
			ID:        nextID(todos),
			Text:      text,
			Priority:  priority,
			CreatedAt: time.Now(),
			Due:       due,
			Tags:      normalizeTags(tags),
		}
		todos = append(todos, todo)
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("added todo %d: %s\n", todo.ID, todo.Text)

	case "list":
		listFs := flag.NewFlagSet("list", flag.ContinueOnError)
		tagFilter := listFs.String("tag", "", "only show todos with this tag")
		searchFilter := listFs.String("search", "", "only show todos whose text contains this substring (case-insensitive)")
		sortBy := listFs.String("sort", "", "sort by: priority, created, or due")
		noColor := listFs.Bool("no-color", false, "disable colored output")
		if err := listFs.Parse(args); err != nil {
			return err
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		if *tagFilter != "" {
			filtered := todos[:0:0]
			for _, t := range todos {
				if hasTag(t.Tags, *tagFilter) {
					filtered = append(filtered, t)
				}
			}
			todos = filtered
		}
		if *searchFilter != "" {
			filtered := todos[:0:0]
			needle := strings.ToLower(*searchFilter)
			for _, t := range todos {
				if strings.Contains(strings.ToLower(t.Text), needle) {
					filtered = append(filtered, t)
				}
			}
			todos = filtered
		}
		if err := sortTodos(todos, *sortBy); err != nil {
			return err
		}
		if len(todos) == 0 {
			fmt.Println("no todos")
			return nil
		}
		now := time.Now()
		useColor := colorEnabled(*noColor)
		for _, t := range todos {
			mark := " "
			if t.Done {
				mark = colorize(useColor, ansiGreen, "x")
			}
			priorityText := colorize(useColor, priorityColor(t.EffectivePriority()), string(t.EffectivePriority()))
			line := fmt.Sprintf("[%s] %d (%s): %s", mark, t.ID, priorityText, t.Text)
			if t.Due != nil {
				line += fmt.Sprintf(" (due %s)", t.Due.Format(dueDateLayout))
			}
			if t.IsOverdue(now) {
				line += colorize(useColor, ansiBold+ansiRed, " [OVERDUE]")
			}
			if len(t.Tags) > 0 {
				line += colorize(useColor, ansiCyan, fmt.Sprintf(" #%s", strings.Join(t.Tags, " #")))
			}
			fmt.Println(line)
		}

	case "done":
		ids, err := parseIDs(args)
		if err != nil {
			return fmt.Errorf("done: %w", err)
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		for _, id := range ids {
			idx := findByID(todos, id)
			if idx == -1 {
				return fmt.Errorf("no todo with id %d", id)
			}
			todos[idx].Done = true
		}
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("marked %d todo(s) done: %v\n", len(ids), ids)

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
		rmFs := flag.NewFlagSet("rm", flag.ContinueOnError)
		force := rmFs.Bool("force", false, "skip the confirmation prompt")
		if err := rmFs.Parse(args); err != nil {
			return err
		}
		ids, err := parseIDs(rmFs.Args())
		if err != nil {
			return fmt.Errorf("rm: %w", err)
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		toRemove := make([]Todo, 0, len(ids))
		for _, id := range ids {
			idx := findByID(todos, id)
			if idx == -1 {
				return fmt.Errorf("no todo with id %d", id)
			}
			toRemove = append(toRemove, todos[idx])
		}
		if !*force {
			var summary strings.Builder
			fmt.Fprintf(&summary, "remove %d todo(s)?", len(toRemove))
			for _, t := range toRemove {
				fmt.Fprintf(&summary, "\n  %d: %q", t.ID, t.Text)
			}
			ok, err := confirm(summary.String())
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("aborted")
				return nil
			}
		}
		todos = removeByIDs(todos, ids)
		if err := saveTodos(path, todos); err != nil {
			return err
		}
		fmt.Printf("removed %d todo(s): %v\n", len(ids), ids)

	case "clear-completed":
		clearFs := flag.NewFlagSet("clear-completed", flag.ContinueOnError)
		force := clearFs.Bool("force", false, "skip the confirmation prompt")
		if err := clearFs.Parse(args); err != nil {
			return err
		}
		todos, err := loadTodos(path)
		if err != nil {
			return err
		}
		var completed []Todo
		for _, t := range todos {
			if t.Done {
				completed = append(completed, t)
			}
		}
		if len(completed) == 0 {
			fmt.Println("no completed todos to clear")
			return nil
		}
		if !*force {
			var summary strings.Builder
			fmt.Fprintf(&summary, "clear %d completed todo(s)?", len(completed))
			for _, t := range completed {
				fmt.Fprintf(&summary, "\n  %d: %q", t.ID, t.Text)
			}
			ok, err := confirm(summary.String())
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("aborted")
				return nil
			}
		}
		remaining := todos[:0:0]
		for _, t := range todos {
			if !t.Done {
				remaining = append(remaining, t)
			}
		}
		if err := saveTodos(path, remaining); err != nil {
			return err
		}
		fmt.Printf("cleared %d completed todo(s)\n", len(completed))

	case "undo":
		if !hasUndoSnapshot(path) {
			return fmt.Errorf("nothing to undo")
		}
		snapshot, err := loadTodos(undoPath(path))
		if err != nil {
			return err
		}
		if err := saveTodos(path, snapshot); err != nil {
			return err
		}
		if err := os.Remove(undoPath(path)); err != nil {
			return err
		}
		fmt.Println("undo: restored previous state")

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
