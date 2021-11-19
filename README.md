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

### list

```sh
go run . list [--tag t] [--search text] [--sort priority|created|due] [--no-color]
```

- `--tag` shows only todos carrying the given tag.
- `--search` shows only todos whose text contains the given substring
  (case-insensitive).
- `--sort` orders the output by priority (high first), creation time
  (oldest first), or due date (soonest first, with undated todos last).
- `--no-color` disables ANSI colors; colors are also disabled automatically
  when the `NO_COLOR` environment variable is set or stdout is not a
  terminal (e.g. when piping output).

```sh
go run . list --tag home
go run . list --search "report"
go run . list --sort priority
go run . list --sort due --no-color
```

### Tags

Attach tags to a todo with repeated `--tag` flags on `add`. Tags are stored
lowercased and deduplicated.

```sh
go run . add --tag work --tag urgent "ship release"
go run . list --tag work
```

### done

```sh
go run . done <id> [id...]
```

Marks one or more todos as done in a single call.

```sh
go run . done 1
go run . done 1 2 3
```

### rm

```sh
go run . rm [--force] <id> [id...]
```

Removes one or more todos. By default this prompts for confirmation and
lists every todo that would be removed; pass `--force` to skip the prompt
(useful in scripts).

```sh
go run . rm 2
go run . rm --force 2 3 4
```

### clear-completed

```sh
go run . clear-completed [--force]
```

Removes every todo currently marked done. Like `rm`, it prompts with a
summary unless `--force` is given.

```sh
go run . clear-completed
go run . clear-completed --force
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
