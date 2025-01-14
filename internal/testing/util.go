package testing

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var compiled = false

// MustCompileOnce ensures geo is built from the current source
func MustCompileOnce(t *testing.T) {
	if compiled {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	mk := exec.CommandContext(ctx, "make", "clean", "all")
	mk.Dir = "../.."
	out, err := mk.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to compile geo woth error: \n%s\n\nOutput:\n%s\n", err, out)
	}
	compiled = true
}

type TestingT interface {
	Fatalf(format string, args ...any)
	Helper()
}

func NewOutputMatcher(s string) *OutputMatcher {
	return &OutputMatcher{
		Original:  s,
		remainder: s,
	}
}

type OutputMatcher struct {
	Original  string
	remainder string
}

func (s *OutputMatcher) MatchText(t TestingT, text string) {
	t.Helper()
	pos := strings.Index(s.remainder, text)
	if pos < 0 {
		t.Fatalf("could not match '%s' in remaining output:\n%s\n", text, s.remainder)
	}
	s.remainder = s.remainder[pos+len(text):]
}
