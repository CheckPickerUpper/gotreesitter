package grammarruntime

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

func TestGoOwnedEOFRecoveryProfileAttaches(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := GoLanguage()
	if !lang.CompactOwnedEOFRecoveryCertified || !lang.CompactConvergedReductionSplitDropsCertified {
		t.Fatal("exact Go profile lacks the owned EOF bundle or existing split drops")
	}
	if lang.CompactStrategy2ErrorRegionCertified ||
		lang.CompactMissingTokenInsertionCertified || lang.CompactS5EOFMissingInsertionCertified ||
		lang.CompactFaithfulS5RecoveryCertified || lang.CompactRecoveryPlainFirstCertified {
		t.Fatal("the Go bundle leaked independent shared-recovery capabilities")
	}
	custom := &gotreesitter.Language{Name: "go"}
	AttachLanguageSupport("go", custom)
	if custom.CompactOwnedEOFRecoveryCertified {
		t.Fatal("same-name custom Go received the owned EOF bundle")
	}
	attached := &gotreesitter.Language{Name: "go"}
	if !attachBuiltinLanguageRuntimeProfile("go", attached) || !attached.CompactOwnedEOFRecoveryCertified {
		t.Fatal("the Go profile did not attach the owned EOF bundle")
	}
	for name, profile := range builtinLanguageRuntimeProfiles {
		if name != "go" && profile.compactOwnedEOFRecovery {
			t.Fatalf("the owned EOF bundle unexpectedly certified %s", name)
		}
	}
	if !builtinLanguageRuntimeProfiles["scala"].compactFaithfulS5Recovery {
		t.Fatal("the independent Scala recovery profile changed")
	}
}
