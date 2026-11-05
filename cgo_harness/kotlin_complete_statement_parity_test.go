//go:build cgo && treesitter_c_parity

package cgoharness

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/odvcencio/gotreesitter/grammars"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TestKotlinCompleteStatementsMatchLockedC compares valid Kotlin whose
// statements have more than one complete reading, and their neighbouring
// forms, with the grammar and runtime pinned for C.
func TestKotlinCompleteStatementsMatchLockedC(t *testing.T) {
	entry := grammars.DetectLanguageByName("kotlin")
	language := entry.Language()
	cLanguage, languageError := COracleLanguage("kotlin")
	if languageError != nil {
		t.Fatal(languageError)
	}
	cParser := sitter.NewParser()
	t.Cleanup(cParser.Close)
	if languageError := cParser.SetLanguage(cLanguage); languageError != nil {
		t.Fatal(languageError)
	}
	folder := "../grammars/testdata/kotlin_complete_statements"
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
			cTree := cParser.Parse(source, nil)
			if cTree == nil {
				t.Fatal("locked C parse returned no tree")
			}
			t.Cleanup(cTree.Close)
			assertLockedCTreeExact(t, filepath.Base(path), tree, language, cTree)
		})
	}
}
