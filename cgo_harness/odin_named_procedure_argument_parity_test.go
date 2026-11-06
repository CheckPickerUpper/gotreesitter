//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TestOdinNamedProcedureArgumentLockedCParity checks that a procedure literal
// passed positionally or after `name =` parses without errors in the locked C
// parser and that the built-in Odin blob produces the same tree.
func TestOdinNamedProcedureArgumentLockedCParity(t *testing.T) {
	cLanguage, err := COracleLanguage("odin")
	if err != nil {
		t.Fatal(err)
	}
	cParser := sitter.NewParser()
	defer cParser.Close()
	if err := cParser.SetLanguage(cLanguage); err != nil {
		t.Fatal(err)
	}
	language := grammars.OdinLanguage()

	for _, tc := range []struct {
		name   string
		source string
	}{
		{"positional", "main :: proc() { f(proc(x: int) {}) }"},
		{"named", "main :: proc() { f(cb = proc(x: int) {}) }"},
		{"named after positional", "main :: proc() { f(1, cb = proc(x: int) { y := x }) }"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte("package p\n\n" + tc.source + "\n")
			cTree := cParser.Parse(source, nil)
			if cTree == nil || cTree.RootNode() == nil {
				t.Fatal("locked C parser returned no tree")
			}
			defer cTree.Close()
			cRoot := cTree.RootNode()
			if cRoot.HasError() || cRoot.EndByte() != uint(len(source)) {
				t.Fatalf("locked C tree contains errors or omits source bytes: %s", cRoot.ToSexp())
			}
			tree, err := gotreesitter.NewParser(language).Parse(source)
			if err != nil {
				t.Fatal(err)
			}
			if tree == nil || tree.RootNode() == nil {
				t.Fatal("Go parser returned no tree")
			}
			defer tree.Release()
			if difference := FirstDivergenceDumpV1(tree.RootNode(), language, cRoot); difference != nil {
				t.Fatalf("Go Odin tree differs from locked C: %+v", difference)
			}
		})
	}
}
