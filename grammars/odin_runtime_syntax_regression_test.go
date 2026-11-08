package grammars

import (
	"testing"

	ts "github.com/odvcencio/gotreesitter"
)

// TestOdinRuntimeSyntaxParsesWithoutErrors checks constructs the Odin core
// library uses, each of which used to end in an ERROR; some, such as a
// procedure type's result in a struct field, swallowed the declarations after
// them.
func TestOdinRuntimeSyntaxParsesWithoutErrors(t *testing.T) {
	for _, test := range []struct {
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
		t.Run(test.name, func(t *testing.T) {
			language := OdinLanguage()
			source := []byte("package p\n\n" + test.source + "\n")
			tree, parseError := ts.NewParser(language).Parse(source)
			if parseError != nil {
				t.Fatalf("parse failed: %v", parseError)
			}
			defer tree.Release()
			root := tree.RootNode()
			if root.HasError() || tree.ParseStopReason() != ts.ParseStopAccepted {
				t.Fatalf("stop=%s has_error=%t tree=%s", tree.ParseStopReason(), root.HasError(), root.SExpr(language))
			}
		})
	}
}
