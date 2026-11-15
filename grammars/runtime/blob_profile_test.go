package grammarruntime

import (
	"github.com/odvcencio/gotreesitter"
	"testing"
)

func TestJavaScriptForestAllowanceRequiresBuiltinProfile(t *testing.T) {
	blob := BlobByName("javascript")
	custom, err := gotreesitter.LoadLanguage(blob)
	if err != nil {
		t.Fatalf("decode caller-owned javascript blob: %v", err)
	}
	custom.WantsForest = true
	if got := custom.AutomaticForestMemoryAllowanceBytes; got != 0 {
		t.Fatalf("caller-owned WantsForest allowance = %d, want zero/full-budget behavior", got)
	}
}
