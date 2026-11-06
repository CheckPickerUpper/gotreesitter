package grammars

import (
	"strings"
	"testing"

	ts "github.com/odvcencio/gotreesitter"
)

// TestOdinProcedureResultsParseWithoutErrors checks that every form of an
// Odin procedure result parses without an error node and with the expected
// result shape. Unnamed result tuples such as `(int, bool)` used to end in an
// ERROR at the second entry.
func TestOdinProcedureResultsParseWithoutErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		source string
		result string
	}{
		{
			name:   "single type",
			source: "f :: proc() -> int { return 1 }",
			result: "(type (identifier))",
		},
		{
			name:   "unnamed types",
			source: "f :: proc() -> (int, bool) { return 1, true }",
			result: "(type (tuple_type (type (identifier)) (type (identifier))))",
		},
		{
			name:   "named types",
			source: "f :: proc() -> (a: int, b: bool) { return 1, true }",
			result: "(type (tuple_type (named_type (identifier) (type (identifier))) (named_type (identifier) (type (identifier)))))",
		},
		{
			name:   "names sharing a type",
			source: "f :: proc() -> (a, b: int) { return 1, 2 }",
			result: "(type (tuple_type (named_type (identifier) (identifier) (type (identifier)))))",
		},
		{
			name:   "three names sharing a type",
			source: "f :: proc() -> (a, b, c: int) { return 1, 2, 3 }",
			result: "(type (tuple_type (named_type (identifier) (identifier) (identifier) (type (identifier)))))",
		},
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
			if !strings.Contains(treeText, "(procedure (parameters) "+test.result+" (block") {
				t.Fatalf("result shape %s missing from %s", test.result, treeText)
			}
		})
	}
}
