package grammarruntime

import "testing"

func TestGoMissingTurnRangeExclusionBindsToGoScanner(t *testing.T) {
	lang := GoLanguage()
	if !lang.RecoveryMissingVersionTurnsCertified {
		t.Fatal("Go lost its missing-version turn certification")
	}
	bound, ok := lang.ExternalScanner.(interface {
		ExcludesMissingVersionTurnsForIncludedRanges() bool
	})
	if !ok || !bound.ExcludesMissingVersionTurnsForIncludedRanges() {
		t.Fatal("the bound scanner lost its included-range mode exclusion")
	}
	if !builtinLanguageRuntimeProfiles["sql"].recoveryMissingVersionTurns {
		t.Fatal("the existing SQL profile grant changed")
	}
}
