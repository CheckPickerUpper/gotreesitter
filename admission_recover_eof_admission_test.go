//go:build !gts_no_parsercorephase0

package gotreesitter_test

import (
	"testing"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

func TestAdmissionRecoverEOFPayloadCountFallsBackForUnfinishedYAML(t *testing.T) {
	lang := grammars.YamlLanguage()
	for _, source := range []string{"a: [1\n", "\"a\": [\"x\"\n"} {
		t.Run(source, func(t *testing.T) {
			production := gts.NewParser(lang)
			production.SetAdmissionCandidateRoute(false)
			want, err := production.Parse([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			defer want.Release()

			candidate := gts.NewParser(lang)
			candidate.SetAdmissionCandidateRoute(true)
			got, err := candidate.Parse([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			defer got.Release()
			if got.RootNode().SExpr(lang) != want.RootNode().SExpr(lang) {
				t.Fatalf("unfinished YAML fallback tree=%s want=%s", got.RootNode().SExpr(lang), want.RootNode().SExpr(lang))
			}
		})
	}
}

func TestAdmissionRecoverEOFIncrementalMarkerFallsBackForUncertifiedScanner(t *testing.T) {
	lang := grammars.YamlLanguage()
	parser := gts.NewParser(lang)
	parser.SetAdmissionCandidateRoute(true)
	oldSource := []byte("[\n")
	oldTree, err := parser.Parse(oldSource)
	if err != nil {
		t.Fatal(err)
	}
	if oldTree == nil || oldTree.RootNode() == nil || !oldTree.RootNode().IsError() {
		if oldTree != nil {
			oldTree.Release()
		}
		t.Fatal("recover_eof source did not produce its ERROR root")
	}
	defer oldTree.Release()
	oldTree.Edit(gts.InputEdit{
		StartByte:   0,
		OldEndByte:  1,
		NewEndByte:  1,
		StartPoint:  gts.Point{},
		OldEndPoint: gts.Point{Column: 1},
		NewEndPoint: gts.Point{Column: 1},
	})
	newSource := []byte("]\n")
	got, profile, err := parser.ParseIncrementalProfiled(newSource, oldTree)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.RootNode() == nil {
		if got != nil {
			got.Release()
		}
		t.Fatal("incremental marker fallback returned no tree")
	}
	defer got.Release()
	freshParser := gts.NewParser(lang)
	freshParser.SetAdmissionCandidateRoute(false)
	want, err := freshParser.Parse(newSource)
	if err != nil {
		t.Fatal(err)
	}
	defer want.Release()
	gotRoot, wantRoot := got.RootNode(), want.RootNode()
	if gotRoot.SExpr(lang) != wantRoot.SExpr(lang) || gotRoot.Range() != wantRoot.Range() {
		t.Fatalf("incremental marker fallback tree=%s range=%+v want fresh=%s range=%+v",
			gotRoot.SExpr(lang), gotRoot.Range(), wantRoot.SExpr(lang), wantRoot.Range())
	}
	requireIncrementalDeepTreeMatchesFresh(t, got, want, lang)
	if profile.ReuseUnsupportedReason == "old tree carried a compact recover_eof EOF runtime" {
		t.Fatalf("recover_eof marker reported the obsolete permanent bar: %+v", profile)
	}
	if !profile.ReuseUnsupported || profile.ReuseUnsupportedReason != "external_scanner_unsupported" ||
		profile.OldTreeReuseRoute || profile.ReusedSubtrees != 0 || profile.ReusedBytes != 0 {
		t.Fatalf("recover_eof marker did not fail closed for the uncertified scanner: %+v", profile)
	}
	runtime := got.ParseRuntime()
	if runtime.StopReason != gts.ParseStopAccepted || !runtime.LastTokenWasEOF ||
		runtime.SourceLen != uint32(len(newSource)) || runtime.LastTokenEndByte != uint32(len(newSource)) ||
		runtime.ExpectedEOFByte != uint32(len(newSource)) || runtime.RootEndByte != uint32(len(newSource)) ||
		runtime.RootEndByte != got.RootNode().EndByte() ||
		runtime.IncrementalOldTreeReuseRoute {
		t.Fatalf("recover_eof fallback runtime=%+v, want fresh EOF at %d", runtime, len(newSource))
	}
}

func TestAdmissionRecoverEOFUnfinishedYAMLIncrementalFallsBackCorrectly(t *testing.T) {
	lang := grammars.YamlLanguage()
	parser := gts.NewParser(lang)
	parser.SetAdmissionCandidateRoute(true)
	oldSource := []byte("\"a\": [\"x\"\n")
	oldTree, err := parser.Parse(oldSource)
	if err != nil {
		t.Fatal(err)
	}
	defer oldTree.Release()
	oldTree.Edit(gts.InputEdit{
		StartByte:   7,
		OldEndByte:  8,
		NewEndByte:  8,
		StartPoint:  gts.Point{Column: 7},
		OldEndPoint: gts.Point{Column: 8},
		NewEndPoint: gts.Point{Column: 8},
	})
	newSource := []byte("\"a\": [\"y\"\n")
	got, _, err := parser.ParseIncrementalProfiled(newSource, oldTree)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Release()
	freshParser := gts.NewParser(lang)
	freshParser.SetAdmissionCandidateRoute(false)
	want, err := freshParser.Parse(newSource)
	if err != nil {
		t.Fatal(err)
	}
	defer want.Release()
	if got.RootNode().SExpr(lang) != want.RootNode().SExpr(lang) {
		t.Fatalf("unfinished YAML incremental tree=%s want fresh=%s", got.RootNode().SExpr(lang), want.RootNode().SExpr(lang))
	}
}

// TestAdmissionRecoverEOFTieKeepsFirstFinishedTreeAsIncumbent pins the
// nightly-fuzz witness Go "," (FuzzAdmissionRouteEquality seed
// 0d6dbe8c80002ec6). Compact recovery finishes two trees at the same error
// cost: the recover_eof ERROR root first, then source_file over an extra
// ERROR after the summary recovery returns to the start state. C's
// ts_parser__select_tree keeps the earlier finished tree as the incumbent and
// takes the later tree on an equal positive cost, so the locked C oracle
// publishes (source_file (ERROR)). The compact route must fold its finished
// trees in that acceptance order, not in header order.
func TestAdmissionRecoverEOFTieKeepsFirstFinishedTreeAsIncumbent(t *testing.T) {
	lang := grammars.GoLanguage()
	tree, telemetry, ok, reason := gts.TryCompactRecoverEOFAcceptTelemetryForTest(
		gts.NewParser(lang), []byte(","),
	)
	if tree != nil {
		defer tree.Release()
	}
	if !ok || reason != "" || tree == nil || tree.RootNode() == nil {
		t.Fatalf("compact route ok=%t reason=%q tree=%t, want the compact route to serve the input", ok, reason, tree != nil)
	}
	if telemetry.RecoverEOFAccepts != 1 {
		t.Fatalf("RecoverEOFAccepts=%d, want exactly one", telemetry.RecoverEOFAccepts)
	}
	root := tree.RootNode()
	if root.Type(lang) != "source_file" || root.IsError() || !root.HasError() ||
		root.StartByte() != 0 || root.EndByte() != 1 || root.ChildCount() != 1 {
		t.Fatalf("root=%s error=%t has_error=%t span=%d..%d children=%d, want source_file[0:1] with one child",
			root.Type(lang), root.IsError(), root.HasError(), root.StartByte(), root.EndByte(), root.ChildCount())
	}
	child := root.Child(0)
	if !child.IsError() || !child.IsExtra() || child.StartByte() != 0 || child.EndByte() != 1 {
		t.Fatalf("child=%s error=%t extra=%t span=%d..%d, want an extra ERROR over 0..1",
			child.Type(lang), child.IsError(), child.IsExtra(), child.StartByte(), child.EndByte())
	}
}
