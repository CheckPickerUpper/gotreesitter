package grammarruntime

import "testing"

func TestCompactGraduationGrantDenominatorAfterEOFRetirement(t *testing.T) {
	activeGrammars := make(map[string]struct{})
	activeFlags := 0
	for name, profile := range builtinLanguageRuntimeProfiles {
		flags := [...]bool{
			profile.compactConvergedSplitDrops,
			profile.compactPrimaryAcceptDerivation,
			profile.compactStrategy2ErrorRegion,
			profile.compactMissingTokenInsertion,
			profile.compactRecoveryPlainFirst,
		}
		for _, active := range flags {
			if !active {
				continue
			}
			activeFlags++
			activeGrammars[name] = struct{}{}
		}
	}

	if activeFlags != 22 || len(activeGrammars) != 13 {
		t.Errorf(
			"compact profile denominator is %d flags across %d grammars, want 22 across 13",
			activeFlags,
			len(activeGrammars),
		)
	}
}
