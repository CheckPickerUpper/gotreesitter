//go:build cgo && treesitter_c_parity

package cgoharness

// C-oracle coverage for the included-ranges route. No test covered
// SetIncludedRanges for any language before this one, which is why a proposed
// deletion of the Go arm's normalizeGoSourceFileRoot member passed a full
// green suite while it turned this route's root from `source_file` into
// `ERROR`.
//
// Production reaches the route through injection: injection.go calls
// childParser.SetIncludedRanges for every injected child, and
// grammars/markdown_injection_register.go maps the `go` and `golang` fence
// languages to the Go grammar.
//
// The parser's padding scan clips to the included ranges, so the route does
// not force recovery between two included ranges. On all four geometries
// below gotreesitter's root has C's symbol, start, child count and error
// flag, and the Go arm's root member does not fire.
//
// READ THIS BEFORE CITING THIS TEST AS EVIDENCE. It compares the root only,
// not the tree under it, and on two geometries gotreesitter's root ends past
// C's (goRootEndsPastC).
//
// Run: GTS_PARITY_ALLOW_HOST=1 go test ./cgo_harness -tags treesitter_c_parity \
//        -run TestIncludedRangesGo -v

import (
	"os"
	"path/filepath"
	"testing"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func includedRangesPointAt(src []byte, off int) (uint, uint) {
	var row, col uint
	for i := 0; i < off && i < len(src); i++ {
		if src[i] == '\n' {
			row++
			col = 0
			continue
		}
		col++
	}
	return row, col
}

// includedRangesRootObservation is one measured root, pinned exactly.
type includedRangesRootObservation struct {
	Kind       string
	StartByte  uint32
	EndByte    uint32
	ChildCount int
	HasError   bool
}

type includedRangesGeometry struct {
	name  string
	spans [2][2]int
	// armRewrites is the node count the Go arm's root member rewrote. A
	// nonzero value proves the member repaired an ERROR root on this
	// geometry, because normalizeGoSourceFileRoot returns early for any other
	// root symbol.
	armRewrites uint64
	c           includedRangesRootObservation
	// goRootEndsPastC is how many bytes after C's root end gotreesitter's
	// root ends. It is nonzero where the last included range ends inside a
	// token: Go has external tokens, so its token source lexes the whole
	// source and filters tokens by range, and the token that starts inside
	// the last range is kept whole. C's lexer stops at the range end.
	goRootEndsPastC uint32
	// note records what this geometry proves and what it does not.
	note string
}

// includedRangesGoGeometries pins C's root for every measured geometry.
// gotreesitter's root is compared with the root C returns in the same run,
// apart from goRootEndsPastC.
//
// The padding scan clips to the included ranges: it does not scan the
// excluded bytes between two ranges as if they had to be whitespace, so a
// non-whitespace gap between ranges does not kill every GLR stack and force
// recovery. Every tested geometry reaches a `source_file` root, and the Go
// arm's root member does not fire because none supplies an ERROR root.
var includedRangesGoGeometries = []includedRangesGeometry{
	{
		name:        "interior_anchors_arm_live",
		spans:       [2][2]int{{26, 150}, {203, 276}},
		armRewrites: 0,
		c:           includedRangesRootObservation{"source_file", 26, 276, 7, true},
		note: "The realistic injection shape: both ranges start at a positive " +
			"offset.",
	},
	{
		name:        "anchored_at_zero_and_eof",
		spans:       [2][2]int{{0, 150}, {203, 276}},
		armRewrites: 0,
		c:           includedRangesRootObservation{"source_file", 0, 276, 7, true},
		note: "A shape an injection child never receives: a Markdown fence " +
			"cannot start at byte 0, because the opening fence line always " +
			"precedes the content.",
	},
	{
		name:            "trimmed_tail_child_count_diverges",
		spans:           [2][2]int{{0, 150}, {203, 250}},
		armRewrites:     0,
		c:               includedRangesRootObservation{"source_file", 0, 250, 6, true},
		goRootEndsPastC: 1,
		note: "The last range ends at byte 250, inside `true` (247 to 251). " +
			"gotreesitter keeps the whole token.",
	},
	{
		name:            "padding_kills_stacks",
		spans:           [2][2]int{{40, 150}, {203, 260}},
		armRewrites:     0,
		c:               includedRangesRootObservation{"source_file", 40, 260, 5, true},
		goRootEndsPastC: 4,
		note: "The last range ends at byte 260, inside `syscall` (257 to 264). " +
			"gotreesitter keeps the whole token.",
	},
}

func measureIncludedRangesRoots(
	t *testing.T,
	src []byte,
	cLang *sitter.Language,
	goLang *gts.Language,
	spans [2][2]int,
) (includedRangesRootObservation, includedRangesRootObservation, uint64) {
	t.Helper()
	cRanges := make([]sitter.Range, 0, len(spans))
	goRanges := make([]gts.Range, 0, len(spans))
	for _, span := range spans {
		start, end := span[0], span[1]
		if end > len(src) {
			end = len(src)
		}
		sr, sc := includedRangesPointAt(src, start)
		er, ec := includedRangesPointAt(src, end)
		cRanges = append(cRanges, sitter.Range{
			StartByte:  uint(start),
			EndByte:    uint(end),
			StartPoint: sitter.Point{Row: sr, Column: sc},
			EndPoint:   sitter.Point{Row: er, Column: ec},
		})
		goRanges = append(goRanges, gts.Range{
			StartByte:  uint32(start),
			EndByte:    uint32(end),
			StartPoint: gts.Point{Row: uint32(sr), Column: uint32(sc)},
			EndPoint:   gts.Point{Row: uint32(er), Column: uint32(ec)},
		})
	}

	cParser := sitter.NewParser()
	defer cParser.Close()
	if err := cParser.SetLanguage(cLang); err != nil {
		t.Fatalf("set C language: %v", err)
	}
	if err := cParser.SetIncludedRanges(cRanges); err != nil {
		t.Fatalf("C SetIncludedRanges: %v", err)
	}
	cTree := cParser.Parse(src, nil)
	if cTree == nil {
		t.Fatal("C parse returned nil")
	}
	defer cTree.Close()
	cRoot := cTree.RootNode()

	goParser := gts.NewParser(goLang)
	goParser.SetIncludedRanges(goRanges)
	goTree, err := goParser.Parse(src)
	if err != nil {
		t.Fatalf("go parse: %v", err)
	}
	if goTree == nil || goTree.RootNode() == nil {
		t.Fatal("go parse returned no tree")
	}
	defer goTree.Release()
	goRoot := goTree.RootNode()

	var armRewrites uint64
	if passes := goTree.ParseRuntime().NormalizationPasses; passes != nil {
		for _, pass := range *passes {
			if pass.Name == "dispatch.go.source-file-root" {
				armRewrites = pass.NodesRewritten
			}
		}
	}

	return includedRangesRootObservation{
			Kind:       cRoot.Kind(),
			StartByte:  uint32(cRoot.StartByte()),
			EndByte:    uint32(cRoot.EndByte()),
			ChildCount: int(cRoot.ChildCount()),
			HasError:   cRoot.HasError(),
		}, includedRangesRootObservation{
			Kind:       goRoot.Type(goLang),
			StartByte:  goRoot.StartByte(),
			EndByte:    goRoot.EndByte(),
			ChildCount: goRoot.ChildCount(),
			HasError:   goRoot.HasError(),
		}, armRewrites
}

func loadIncludedRangesGoFixture(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "testdata", "included_ranges", "go_two_fences.go")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read included-ranges fixture: %v", err)
	}
	return src
}

// TestIncludedRangesGoRootParity compares gotreesitter's root with the root C
// returns in the same run on the included-ranges route, per geometry, and
// pins C's root so a change in the oracle fails here too. The one allowed
// difference is goRootEndsPastC.
func TestIncludedRangesGoRootParity(t *testing.T) {
	cLang := loadCanonicalGoCLanguage(t)
	goLang := grammars.GoLanguage()
	if goLang == nil {
		t.Skip("go grammar unavailable")
	}
	src := loadIncludedRangesGoFixture(t)

	for _, geometry := range includedRangesGoGeometries {
		geometry := geometry
		t.Run(geometry.name, func(t *testing.T) {
			t.Setenv("GTS_DISPATCHER_CENSUS", "1")
			gotC, gotGo, armRewrites := measureIncludedRangesRoots(t, src, cLang, goLang, geometry.spans)

			t.Logf("note: %s", geometry.note)
			t.Logf("C   %+v", gotC)
			t.Logf("GTS %+v", gotGo)
			t.Logf("dispatch.go.source-file-root rewrote %d nodes", armRewrites)

			if gotC != geometry.c {
				t.Errorf("C oracle root moved: got %+v, pinned %+v", gotC, geometry.c)
			}
			wantGo := gotC
			wantGo.EndByte += geometry.goRootEndsPastC
			if gotGo != wantGo {
				t.Errorf("gotreesitter root %+v, want C's root from this run with the end %d bytes later: %+v", gotGo, geometry.goRootEndsPastC, wantGo)
			}
			if armRewrites != geometry.armRewrites {
				t.Errorf("dispatch.go.source-file-root rewrites = %d, pinned %d", armRewrites, geometry.armRewrites)
			}
		})
	}
}

// TestIncludedRangesGoArmGuardsRootSymbol is the root-symbol regression
// guard for this route. It is separate from the pin table above so its
// intent cannot be misread: on every pinned geometry, the root symbol must
// keep matching C, and the Go arm's root member must stay inert, because the
// padding-scan fix keeps this route out of recovery and the member no longer
// has an ERROR root to repair. If a future change reopens the route to
// recovery on any pinned geometry, the member fires again, this guard's
// armRewrites check fails, and — if the member's repair does not fully
// recover the root symbol — the Kind check fails too. Nothing here claims
// the route is at span or child-count parity.
func TestIncludedRangesGoArmGuardsRootSymbol(t *testing.T) {
	cLang := loadCanonicalGoCLanguage(t)
	goLang := grammars.GoLanguage()
	if goLang == nil {
		t.Skip("go grammar unavailable")
	}
	src := loadIncludedRangesGoFixture(t)

	for _, geometry := range includedRangesGoGeometries {
		geometry := geometry
		t.Run(geometry.name, func(t *testing.T) {
			t.Setenv("GTS_DISPATCHER_CENSUS", "1")
			gotC, gotGo, armRewrites := measureIncludedRangesRoots(t, src, cLang, goLang, geometry.spans)
			if armRewrites != 0 {
				t.Fatalf("dispatch.go.source-file-root rewrote %d nodes; want 0 now that the padding scan clips to included ranges", armRewrites)
			}
			if gotGo.Kind != gotC.Kind {
				t.Fatalf("root symbol diverged from C: gotreesitter=%q C=%q", gotGo.Kind, gotC.Kind)
			}
			if gotGo.Kind == "ERROR" {
				t.Fatalf("root is an ERROR node while C publishes %q", gotC.Kind)
			}
		})
	}
}
