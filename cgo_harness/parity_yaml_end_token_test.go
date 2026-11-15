//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"testing"

	"github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

// TestParityYAMLEndTokenBeforeTrailingBytes checks inputs where YAML's lex
// state after an open bracket or brace has no token, so tree-sitter's lexer
// accepts the end token there, before the trailing newline or space. C ends
// the parse and its ERROR root at that token; the Go parser must do the same
// instead of reading the trailing byte as an error run.
func TestParityYAMLEndTokenBeforeTrailingBytes(t *testing.T) {
	cLanguage, err := ParityCLanguage("yaml")
	if err != nil {
		t.Fatalf("load C yaml: %v", err)
	}
	goLanguage := grammars.YamlLanguage()
	for _, source := range []string{"[\n", "[ ", "{\n"} {
		t.Run(source, func(t *testing.T) {
			cTree := compactT3ParseC(t, cLanguage, []byte(source))
			defer cTree.Close()
			parser := gotreesitter.NewParser(goLanguage)
			parser.SetAdmissionCandidateRoute(false)
			goTree, err := parser.Parse([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			defer goTree.Release()
			if divergences := compactT3StructuralDivergences(goTree.RootNode(), goLanguage, cTree.RootNode()); len(divergences) != 0 {
				t.Fatalf("Go tree %s [%d,%d) differs from C %s [%d,%d): %+v",
					goTree.RootNode().SExpr(goLanguage), goTree.RootNode().StartByte(), goTree.RootNode().EndByte(),
					cTree.RootNode().ToSexp(), cTree.RootNode().StartByte(), cTree.RootNode().EndByte(), divergences)
			}
		})
	}
}
