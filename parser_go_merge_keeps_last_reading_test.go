package gotreesitter_test

import "testing"

// TestParseGoKeepsTheLastReadingThatCanAccept covers valid Go whose readings
// fill a merge key's survivor budget before the one reading that survives to
// the end arrives. Each source is minimized from a real file go/parser and the
// C runtime both accept; each must parse with no ERROR node at the default
// budget.
func TestParseGoKeepsTheLastReadingThatCanAccept(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{
			name: "a for over a call after two calls, whose body passes a generic map literal",
			source: `package c
func e(){g.e(t)
g.e(t)
for s.d(s.s(n)){e(map[g]o.n[g]{})}}
`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree, lang := parseGo(t, tc.source)
			defer tree.Release()
			root := tree.RootNode()
			if root.HasError() {
				t.Fatalf("root has error:\n%s", root.SExpr(lang))
			}
			if got, want := root.EndByte(), uint32(len(tc.source)); got != want {
				t.Fatalf("root end = %d, want %d", got, want)
			}
			if findNamedChild(lang, root, "for_statement") == nil {
				t.Fatalf("missing for_statement:\n%s", root.SExpr(lang))
			}
		})
	}
}
