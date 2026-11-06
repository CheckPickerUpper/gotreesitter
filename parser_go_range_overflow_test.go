package gotreesitter_test

import "testing"

// TestParseGoRangeBodyKeptPastMergeOverflow covers a for-range over a bare
// identifier whose body starts with a call. At the body's "{" the parser
// forks: the identifier ends the range clause and "{" opens the body, or the
// identifier is a composite literal's type and "{" opens its literal value.
// The literal reading forks again at every "." and "(" of the call, so many
// same-state, same-offset stacks compete for the per-key survivor budget. A
// survivor budget ranked by stack depth kept only the deeper literal readings,
// which die at the loop's closing brace, and the file came back with an ERROR
// node. Each source is minimized from a real file go/parser accepts.
func TestParseGoRangeBodyKeptPastMergeOverflow(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{
			name: "reading in a func literal with a generic result after a parenthesized func literal",
			source: `package p

func f() {
	for a = range b {
		(c) = d{}
	}
	(func() {})
	g(h(), func() i[j] {
		for k = range l {
			m.n(o.p(q.r(s.t((u)))), func() {})
		}
	})
}
`,
		},
		{
			name: "two ranges, one inside a parenthesized func literal argument",
			source: `package typescript

import ()

func TestMain() {
	mo.Fold((func() Either[Failure, []string] {
		for range files {
			mo.Fold(mo.TupleToResult(os.Stat(file)), func() {})
		}
	}), func() []string {})
	for file = range filled {
		mo((file), func() {}, func() {})
	}
}
`,
		},
		{
			name: "ranges inside a func literal argument after an empty block",
			source: `package python

import ()

func TestMain() {
	{
	}
	filled.Fold(withUvStores(tool, func() mo[string] {
		for _ = range files {
			mo((WriteFile()), func() {})
		}
		for file = range files {
			mo.Fold(mo.TupleToResult(os.Stat(file)), func() {})
		}
	}), func() {})
	{
		Fold((file), func() {}, bool{})
	}
}
`,
		},
		{
			name: "ranges after assigning to a parenthesized composite literal",
			source: `package golang

import ()

func TestMain() {
	(int{}) = mo.Fold((func() Either[outcomes] {
		for file = range files {
			(file)
		}
		for _ = range files {
			mo.Fold(mo.TupleToResult(os.Stat(file)), func() {})
		}
	}), func() {})
}
`,
		},
		{
			name: "a bare range inside a func literal argument after a composite literal",
			source: `package kotlin_test

func TestKotlinSyntaxReadsEveryConcept() {
	expect("calls", []string{})
	for range entries {
		t(entry, func() { mo.Fold(mo.TupleToResult(os.ReadFile(filepath.Join(entry))), func([]byte) {}) })
	}
}
`,
		},
		{
			name: "a range after a deep call chain and a composite literal argument",
			source: `package kotlin_test

import ()

func TestASecondBuildReusesTheDaemonTheFirstStarted() {
	if (language); host {
		(t)()
		{
		}
	}
	mo((failure), func() {})
}
func registeredDaemons() {
	Fold((pattern), string{})
	for log = range logs {
		Fold((strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(filepath.Base(log))))), func() {})
	}
}
`,
		},
		{
			name: "an assigned range over a call with nested calls and a func literal",
			source: `package python_test

import ()

func totalsLookedUpAfter() {
	for contents = range rewritten {
		mo.Fold(mo.EmptyableToOption(os.MkdirAll((filepath.Join(root)))), func() {})
	}
}
`,
		},
		{
			name: "a bare range whose call passes func literals holding calls",
			source: `package odin

import ()

func checkArguments() {}
func () kind() {
	for range files {
		mo.Fold(mo.TupleToResult(os.ReadFile(filepath.Join((file)))), func() { Fold((tables), func() {}, func() {}) }, func() {})
	}
}
`,
		},
		{
			name: "an assigned range after a composite literal assignment in a range body",
			source: `package generatedcurrent

import ()

func () readsOf() {
	for listedPackage = range module {
		(Dir) = command{}
	}
	(func() {})
	Fold(WorkingTreeFiles(), func() mo[moduleFiles] {
		for listedPath = range listed {
			mo.Fold(mo.TupleToResult(os.ReadFile(filepath.Join((listedPath)))), func() {})
		}
	})
}
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
