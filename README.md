# go-cli-todo

A minimal command-line todo list manager written in Go, using only the standard library.

Todos are persisted as JSON in `~/.go-cli-todo.json`.

## Build

```sh
go build ./...
```

## Usage

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
