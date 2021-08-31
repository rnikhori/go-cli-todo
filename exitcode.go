package main

import (
	"errors"
	"fmt"
)

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

type cliError struct {
	err  error
	code int
}

func (e *cliError) Error() string { return e.err.Error() }
func (e *cliError) Unwrap() error { return e.err }
func (e *cliError) ExitCode() int { return e.code }

func usageErrorf(format string, args ...any) error {
	return &cliError{err: fmt.Errorf(format, args...), code: exitUsage}
}

func exitCodeFor(err error) int {
	if err == nil {
		return exitOK
	}
	var ce *cliError
	if errors.As(err, &ce) {
		return ce.code
	}
	return exitError
}
