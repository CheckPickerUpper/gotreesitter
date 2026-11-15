package grammarruntime

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

func TestGleamSelfConflictProfileAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := Language("gleam")
	if len(lang.ConflictPolicies) != 1 {
		t.Fatalf("Gleam policies=%d, want 1", len(lang.ConflictPolicies))
	}
	p := lang.ConflictPolicies[0]
	if p.Kind != gotreesitter.ConflictPolicyDeclaredSelfReduceReduceLongest ||
		p.State != gotreesitter.ConflictPolicyAnyState || p.Lookahead != gotreesitter.ConflictPolicyAnyLookahead ||
		len(p.ReduceSymbols) != 1 || lang.SymbolNames[p.ReduceSymbols[0]] != "integer" {
		t.Fatalf("unexpected Gleam policy: %+v", p)
	}
}
