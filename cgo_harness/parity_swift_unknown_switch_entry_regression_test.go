//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// TestParitySwiftUnknownAttributeStartsSwitchEntry checks that an `@unknown`
// switch entry after another entry's statements ends that entry and starts a
// new switch_entry with modifiers, as locked C tree-sitter-swift reads it.
func TestParitySwiftUnknownAttributeStartsSwitchEntry(t *testing.T) {
	tests := []struct {
		name   string
		source string
		// statementEmptyTuple is set when the source holds `()` statements.
		// docs/grammar-ownership.md owns that divergence: locked C reads
		// each as a tuple_expression holding `value: (bang (MISSING "!"))`,
		// and this grammar reads it with no children. Go must read the
		// source clean and match C's tree with only those children removed.
		// Once tree-sitter-swift's own grammar reads `()` clean, C matches
		// Go exactly and this field goes.
		statementEmptyTuple bool
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
			statementEmptyTuple: true,
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
					if test.statementEmptyTuple {
						assertLockedCTreeExceptEmptyTupleBang(t, "Swift @unknown switch entry", goTree, goLang, cTree)
						return
					}
					assertLockedCTreeExact(t, "Swift @unknown switch entry", goTree, goLang, cTree)
				})
			}
		})
	}
}

// assertLockedCTreeExceptEmptyTupleBang checks that Go reads the source clean
// and that Go's tree equals locked C's tree after one transform: each
// tuple_expression's `value: (bang (MISSING "!"))` child is removed. Every
// other difference fails, and so does a C tree with no such child, because C
// then reads `()` clean and the case belongs back on exact parity.
func assertLockedCTreeExceptEmptyTupleBang(t *testing.T, label string, goTree *gotreesitter.Tree, goLang *gotreesitter.Language, cTree *sitter.Tree) {
	t.Helper()
	goRoot := goTree.RootNode()
	if goRoot.HasError() {
		t.Fatalf("%s: Go reads an error: %s", label, goRoot.SExpr(goLang))
	}
	goNodes := DumpV1FromGo(goRoot, goLang).Nodes
	var cNodes []DumpV1Node
	removed := 0
	cRoot := cTree.RootNode()
	walkDumpV1CWithoutEmptyTupleBang(cRoot, rootPath(cRoot.Kind()), "", &cNodes, &removed)
	if removed == 0 {
		t.Fatalf("%s: locked C reads every () clean; return this case to exact parity", label)
	}
	for i := 0; i < len(goNodes) || i < len(cNodes); i++ {
		switch {
		case i >= len(goNodes):
			t.Fatalf("%s: C has extra node %+v", label, cNodes[i])
		case i >= len(cNodes):
			t.Fatalf("%s: Go has extra node %+v", label, goNodes[i])
		case goNodes[i] != cNodes[i]:
			t.Fatalf("%s: node %d differs:\nGo: %+v\nC:  %+v", label, i, goNodes[i], cNodes[i])
		}
	}
	t.Logf("%s: Go matches locked C at %d nodes, with %d empty-tuple bang children removed from C", label, len(goNodes), removed)
}

// walkDumpV1CWithoutEmptyTupleBang flattens a C tree as walkDumpV1C does,
// leaving out each tuple_expression child that is `value: (bang (MISSING
// "!"))`. Child counts, child indices and HasError are computed on the tree
// without them. It returns whether the node's kept subtree holds an error.
func walkDumpV1CWithoutEmptyTupleBang(node *sitter.Node, path, field string, out *[]DumpV1Node, removed *int) bool {
	start := node.StartPosition()
	end := node.EndPosition()
	index := len(*out)
	*out = append(*out, DumpV1Node{
		Path:      path,
		Type:      node.Kind(),
		Field:     field,
		IsNamed:   node.IsNamed(),
		IsExtra:   node.IsExtra(),
		StartByte: uint32(node.StartByte()),
		EndByte:   uint32(node.EndByte()),
		StartPoint: DumpV1Point{
			Row:    uint32(start.Row),
			Column: uint32(start.Column),
		},
		EndPoint: DumpV1Point{
			Row:    uint32(end.Row),
			Column: uint32(end.Column),
		},
	})
	hasError := node.IsError() || node.IsMissing()
	kept := 0
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(uint(i))
		childField := node.FieldNameForChild(uint32(i))
		if node.Kind() == "tuple_expression" && isEmptyTupleMissingBang(child, childField) {
			*removed++
			continue
		}
		if walkDumpV1CWithoutEmptyTupleBang(child, childPath(path, child.Kind(), kept), childField, out, removed) {
			hasError = true
		}
		kept++
	}
	(*out)[index].ChildCount = kept
	(*out)[index].HasError = hasError
	return hasError
}

// isEmptyTupleMissingBang reports whether a tuple_expression child is the
// `value: (bang (MISSING "!"))` locked C writes for `()`.
func isEmptyTupleMissingBang(child *sitter.Node, field string) bool {
	if field != "value" || child.Kind() != "bang" || child.ChildCount() != 1 {
		return false
	}
	only := child.Child(0)
	return only.IsMissing() && only.Kind() == "!" && only.ChildCount() == 0
}
