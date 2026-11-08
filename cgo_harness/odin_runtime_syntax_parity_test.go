//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TestOdinRuntimeSyntaxLockedCParity checks that constructs the Odin core
// library uses parse without errors in the locked C parser and that the
// built-in Odin blob produces the same tree.
func TestOdinRuntimeSyntaxLockedCParity(t *testing.T) {
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
		{name: "procedure field type with a result", source: "Map_Info :: struct {\n\tkey_equal: proc \"contextless\" (lhs, rhs: rawptr) -> bool,\n\tkey_hasher: proc \"contextless\" (key: rawptr, seed: uintptr) -> uintptr,\n}"},
		{name: "fixed capacity dynamic array", source: "clear_fixed :: proc(array: ^$D/[dynamic; $N]$E) {\n\tbuffer: [dynamic; 8]int\n}"},
		{name: "procedure group entry with where", source: "clear :: proc{\n\tclear_dynamic_array,\n\tclear_map where MAP_ENABLED,\n}"},
		{name: "declaration with the type left out", source: "f :: proc() {\n\tsd: = N - 1\n}"},
		{name: "struct alignment from an expression", source: "Buddy_Block :: struct #align(align_of(uint)) {\n\tsize: uint,\n}"},
		{name: "simd literal", source: "f :: proc() -> V {\n\treturn transmute(V)#simd[4]i32{1, 2, 3, 4}\n}"},
		{name: "polymorphic parameter default", source: "Chan :: struct($T: typeid, $D: Direction = Direction.Both) {\n\timpl: ^Raw_Chan,\n}"},
		{name: "implicit selector type argument", source: "as_send :: proc(c: $C/Chan($T, $D)) -> (s: Chan(T, .Send)) {\n}"},
		{name: "iteration after a statement", source: "f :: proc() {\n\tfor first := true; key, value in m {\n\t\tfirst = false\n\t}\n}"},
		{name: "tagged statement after defer", source: "f :: proc() {\n\tdefer #no_bounds_check for i in 0 ..< n {\n\t}\n}"},
		{name: "digit separator after the base", source: "O_CREAT :: 0o0_000_100\nPATH :: 0o_10000000\nMASK :: 0b_1010"},
		{name: "array type compared in where and when", source: "encode :: proc(dst: $T) where T == io.Writer || T == []byte {\n\twhen T == []byte {\n\t\tn := 1\n\t}\n}"},
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
