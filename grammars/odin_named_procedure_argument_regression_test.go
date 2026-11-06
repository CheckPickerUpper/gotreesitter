package grammars

import (
	"strings"
	"testing"

	ts "github.com/odvcencio/gotreesitter"
)

// TestOdinNamedProcedureArgumentParsesWithoutErrors checks that a procedure
// literal passed after `name =` in a call parses like a positional one. It
// used to end in an ERROR that could cover the whole file.
func TestOdinNamedProcedureArgumentParsesWithoutErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		source string
	}{
		{name: "positional", source: "main :: proc() { f(proc(x: int) {}) }"},
		{name: "named", source: "main :: proc() { f(cb = proc(x: int) {}) }"},
		{name: "named after positional", source: "main :: proc() { f(1, cb = proc(x: int) { y := x }) }"},
	} {
		t.Run(test.name, func(t *testing.T) {
			language := OdinLanguage()
			source := []byte("package p\n\n" + test.source + "\n")
			tree, parseError := ts.NewParser(language).Parse(source)
			if parseError != nil {
				t.Fatalf("parse failed: %v", parseError)
			}
			defer tree.Release()
			root := tree.RootNode()
			treeText := root.SExpr(language)
			if root.HasError() || tree.ParseStopReason() != ts.ParseStopAccepted {
				t.Fatalf("stop=%s has_error=%t tree=%s", tree.ParseStopReason(), root.HasError(), treeText)
			}
			if !strings.Contains(treeText, "(call_expression") || !strings.Contains(treeText, "(procedure (parameters (parameter") {
				t.Fatalf("call with a procedure argument missing from %s", treeText)
			}
		})
	}
}
