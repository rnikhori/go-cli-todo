package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestExitCodeFor(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"nil error", nil, exitOK},
		{"plain error", errors.New("boom"), exitError},
		{"usage error", usageErrorf("bad flag"), exitUsage},
		{"wrapped usage error", wrapForTest(usageErrorf("bad flag")), exitUsage},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := exitCodeFor(c.err); got != c.want {
				t.Errorf("exitCodeFor(%v) = %d, want %d", c.err, got, c.want)
			}
		})
	}
}

func wrapForTest(err error) error {
	return fmt.Errorf("context: %w", err)
}

func TestDispatchErrorsClassifyAsUsageOrRuntime(t *testing.T) {
	path := tempStore(t)

	if err := dispatch("bogus", nil, path); exitCodeFor(err) != exitUsage {
		t.Errorf("unknown command: exitCodeFor = %d, want %d", exitCodeFor(err), exitUsage)
	}
	if err := dispatch("add", nil, path); exitCodeFor(err) != exitUsage {
		t.Errorf("add with no text: exitCodeFor = %d, want %d", exitCodeFor(err), exitUsage)
	}
	if err := dispatch("done", []string{"99"}, path); exitCodeFor(err) != exitError {
		t.Errorf("done unknown id: exitCodeFor = %d, want %d", exitCodeFor(err), exitError)
	}
}

func TestCliErrorUnwrap(t *testing.T) {
	base := errors.New("underlying")
	wrapped := &cliError{err: base, code: exitUsage}
	if !errors.Is(wrapped, base) {
		t.Error("expected cliError to unwrap to the underlying error")
	}
	if wrapped.Error() != base.Error() {
		t.Errorf("Error() = %q, want %q", wrapped.Error(), base.Error())
	}
}
