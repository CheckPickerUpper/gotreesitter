package gotreesitter_test

import (
	"testing"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

// TestParseGoFinishedTreeEndsRecoveryOfACostlierPausedVersion covers a parse
// whose finished tree is cheaper than the paused version the condense step
// would resume. C stops there, so recovery from the paused version never
// accepts a later tree that would win the equal-cost tie. For each source the
// C runtime returns the ERROR node as an extra child of source_file; both
// admission routes must return the same.
func TestParseGoFinishedTreeEndsRecoveryOfACostlierPausedVersion(t *testing.T) {
	entry := grammars.DetectLanguageByName("go")
	if entry == nil {
		t.Fatal("Go is not registered")
	}
	lang := entry.Language()
	for _, source := range []string{",", ",\n"} {
		for _, compact := range []bool{false, true} {
			parser := gts.NewParser(lang)
			parser.SetAdmissionCandidateRoute(compact)
			tree, err := parser.Parse([]byte(source))
			if err != nil {
				t.Fatalf("source=%q compact=%t: %v", source, compact, err)
			}
			root := tree.RootNode()
			if root.Type(lang) != "source_file" || root.ChildCount() != 1 {
				tree.Release()
				t.Fatalf("source=%q compact=%t: root = %s, want one child under source_file", source, compact, root.SExpr(lang))
			}
			child := root.Child(0)
			if !child.IsError() || !child.IsExtra() || child.StartByte() != 0 || child.EndByte() != 1 {
				tree.Release()
				t.Fatalf("source=%q compact=%t: child %s[%d:%d] error=%t extra=%t, want an extra ERROR[0:1]",
					source, compact, child.Type(lang), child.StartByte(), child.EndByte(), child.IsError(), child.IsExtra())
			}
			tree.Release()
		}
	}
}
