package grammarruntime

import "testing"

func TestGoStringTerminalsRemainVisibleBeforeReduction(t *testing.T) {
	lang := GoLanguage()
	for _, name := range []string{"raw_string_literal_content", "interpreted_string_literal_content"} {
		found := false
		for symbol, symbolName := range lang.SymbolNames {
			if uint32(symbol) >= lang.TokenCount || symbolName != name {
				continue
			}
			found = true
			if !lang.SymbolMetadata[symbol].Visible || !lang.SymbolMetadata[symbol].Named {
				t.Fatalf("terminal %s is hidden before its enclosing string reduces", name)
			}
		}
		if !found {
			t.Fatalf("missing visible terminal %s", name)
		}
	}
}
