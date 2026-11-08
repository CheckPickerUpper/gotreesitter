//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TestParitySwiftUnknownAttributeStartsSwitchEntry checks that an `@unknown`
// switch entry after another entry's statements ends that entry and starts a
// new switch_entry with modifiers, as locked C tree-sitter-swift reads it.
func TestParitySwiftUnknownAttributeStartsSwitchEntry(t *testing.T) {
	tests := []struct {
		name   string
		source string
		// lockedCReadsMissingToken is set when locked C itself reads the
		// source with a MISSING token, so Go must match that tree exactly,
		// errors included.
		lockedCReadsMissingToken bool
	}{
		{
			name: "unknown-default-after-case",
			source: "func f(e: E) -> String {\n" +
				"  switch e {\n" +
				"  case .a: return \"a\"\n" +
				"  @unknown default: return \"x\"\n" +
				"  }\n" +
				"}\n",
		},
		{
			name: "unknown-case-after-case",
			source: "func f(e: E) -> String {\n" +
				"  switch e {\n" +
				"  case .a:\n" +
				"    return \"a\"\n" +
				"  @unknown case .b:\n" +
				"    return \"b\"\n" +
				"  }\n" +
				"}\n",
		},
		{
			name: "unknown-default-after-empty-tuple",
			source: "switch x {\n" +
				"case .a:\n" +
				"    ()\n" +
				"@unknown default:\n" +
				"    ()\n" +
				"}\n",
			lockedCReadsMissingToken: true,
		},
		{
			name: "unknown-default-first",
			source: "func f(e: E) -> String {\n" +
				"  switch e {\n" +
				"  @unknown default: return \"x\"\n" +
				"  }\n" +
				"}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := []byte(test.source)
			cLang, err := ParityCLanguage("swift")
			if err != nil {
				t.Fatalf("load locked Swift C parser: %v", err)
			}
			cParser := sitter.NewParser()
			t.Cleanup(cParser.Close)
			if err := cParser.SetLanguage(cLang); err != nil {
				t.Fatalf("set locked Swift C language: %v", err)
			}
			cTree := cParser.Parse(source, nil)
			if cTree == nil || cTree.RootNode() == nil {
				t.Fatal("locked Swift C parser returned no tree")
			}
			t.Cleanup(cTree.Close)

			for _, route := range []struct {
				name      string
				candidate bool
			}{
				{name: "production", candidate: false},
				{name: "compact", candidate: true},
			} {
				t.Run(route.name, func(t *testing.T) {
					candidate := route.candidate
					goTree, goLang, err := parseWithGo(parityCase{
						name:           "swift",
						source:         test.source,
						candidateRoute: &candidate,
					}, source, nil)
					if err != nil {
						t.Fatalf("parse Swift with Go: %v", err)
					}
					t.Cleanup(func() { releaseGoTree(goTree) })
					if goTree.RootNode().HasError() {
						t.Logf("Go tree: %s", goTree.RootNode().SExpr(goLang))
					}
					if test.lockedCReadsMissingToken {
						if !cTree.RootNode().HasError() {
							t.Fatal("locked C no longer reads a MISSING token here; drop lockedCReadsMissingToken")
						}
						assertLockedCTreeExactWithErrors(t, "Swift @unknown switch entry", goTree, goLang, cTree)
						return
					}
					assertLockedCTreeExact(t, "Swift @unknown switch entry", goTree, goLang, cTree)
				})
			}
		})
	}
}
