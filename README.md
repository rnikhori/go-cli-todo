# go-cli-todo

A command-line todo list manager written in Go, using only the standard library.

Todos are persisted as JSON, by default in `~/.go-cli-todo.json`.

## Build

```sh
go build ./...
```

## Storage location

By default todos are stored in `~/.go-cli-todo.json`. You can override this:

```sh
# via flag (highest precedence)
go run . --store /path/to/todos.json list

# via environment variable
export TODO_STORE=/path/to/todos.json
go run . list
```

## Commands

### add

```sh
go run . add [--priority low|medium|high] [--due YYYY-MM-DD] [--tag t]... <text>
```

Adds a new todo. Priority defaults to `medium` if not specified. `--tag` may be
repeated to attach multiple tags.

```sh
go run . add "buy milk"
go run . add --priority high --due 2026-12-01 --tag home --tag errand "renew passport"
```

### edit

```sh
go run . edit <id> <text>
```

Replaces the text of an existing todo.

```sh
go run . edit 1 "buy oat milk"
```

### Due dates

Due dates are given as `YYYY-MM-DD` on `add`. Any undone todo whose due date has
passed is marked `[OVERDUE]` in `list` output.

```sh
go run . add --due 2026-01-01 "file taxes"
go run . list
```

## Basic workflow example

```sh
go run . add "buy milk"
go run . add "write report"
go run . list
go run . done 1
go run . list
go run . rm 2
go run . list
```

Or, after building:

```sh
./go-cli-todo add "buy milk"
./go-cli-todo list
./go-cli-todo done 1
./go-cli-todo rm 2
```
