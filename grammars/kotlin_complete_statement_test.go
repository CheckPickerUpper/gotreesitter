//go:build !grammar_subset || grammar_subset_kotlin

package grammars

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKotlinValidSourcesParseWithoutErrors parses valid Kotlin whose
// statements have more than one complete reading. The C runtime reads each of
// these files without an error node; so must the Go runtime.
func TestKotlinValidSourcesParseWithoutErrors(t *testing.T) {
	entry := DetectLanguageByName("kotlin")
	if entry == nil {
		t.Fatal("kotlin is not registered")
	}
	folder := "testdata/kotlin_complete_statements"
	files, readError := os.ReadDir(folder)
	if readError != nil {
		t.Fatal(readError)
	}
	paths := []string{"../internal/benchfixtures/testdata/real/kotlin"}
	for _, file := range files {
		paths = append(paths, filepath.Join(folder, file.Name()))
	}
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, readError := os.ReadFile(path)
			if readError != nil {
				t.Fatal(readError)
			}
			tree, parseError := entry.Parse(source)
			if parseError != nil {
				t.Fatal(parseError)
			}
			t.Cleanup(tree.Release)
			root := tree.RootNode()
			if root.HasErrorOrMissing() {
				t.Fatalf("valid Kotlin contains an error or missing node:\n%s", root.SExpr(tree.Language()))
			}
			if root.StartByte() != 0 || root.EndByte() != uint32(len(source)) {
				t.Fatalf("root covers %d..%d, want 0..%d", root.StartByte(), root.EndByte(), len(source))
			}
		})
	}
}
