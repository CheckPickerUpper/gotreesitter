package grammarruntime

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

func TestBuiltinExternalScannerRetryProfilesAttach(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	tests := []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "c_sharp", load: CSharpLanguage},
		{name: "crystal", load: CrystalLanguage},
		{name: "dart", load: DartLanguage},
		{name: "kotlin", load: KotlinLanguage},
		{name: "matlab", load: MatlabLanguage},
		{name: "python", load: PythonLanguage},
		{name: "swift", load: SwiftLanguage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := tt.load()
			if lang.ExternalScanner == nil {
				t.Fatal("ExternalScanner = nil, want attached scanner")
			}
			if got := lang.ExternalScannerFullParseRetryPolicy; got != gotreesitter.ExternalScannerFullParseRetrySkipRepeat {
				t.Fatalf("ExternalScannerFullParseRetryPolicy = %d, want certified skip-repeat policy", got)
			}
		})
	}
}

func TestCrystalAndMatlabKeepAcceptedErrorRetryLadder(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	for _, tt := range []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "crystal", load: CrystalLanguage},
		{name: "matlab", load: MatlabLanguage},
	} {
		t.Run(tt.name, func(t *testing.T) {
			lang := tt.load()
			if got := lang.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want complete ladder", got)
			}
		})
	}
}

func TestHTMLProfileCertifiesCompleteCompactRecovery(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := HtmlLanguage()
	if !lang.CompactStrategy2ErrorRegionCertified || !lang.CompactMissingTokenInsertionCertified {
		t.Fatal("the HTML profile did not attach both compact recovery capabilities")
	}
	if lang.CompactS5EOFMissingInsertionCertified {
		t.Fatal("the HTML profile unexpectedly enabled EOF S5 missing insertion")
	}
	if lang.CompactRecoveryPlainFirstCertified {
		t.Fatal("the HTML profile unexpectedly enabled plain-first recovery")
	}
	if lang.CompactFaithfulS5RecoveryCertified {
		t.Fatal("the HTML profile unexpectedly enabled the Scala S5 route")
	}
}

func TestJavaScriptProfileCertifiesCompactRecoveryFrontier(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := JavascriptLanguage()
	if !lang.CompactStrategy2ErrorRegionCertified || !lang.CompactMissingTokenInsertionCertified ||
		!lang.CompactRecoveryPlainFirstCertified || !lang.CompactRecoveryTrailingLineageRetirementCertified ||
		!lang.CompactRecoveryErrorModeKeywordCaptureCertified {
		t.Fatal("the JavaScript profile did not attach its compact recovery capabilities")
	}
	if lang.CompactFaithfulS5RecoveryCertified {
		t.Fatal("the JavaScript profile unexpectedly enabled the Scala S5 route")
	}
	if lang.CompactS5EOFMissingInsertionCertified {
		t.Fatal("the JavaScript profile unexpectedly enabled EOF S5 missing insertion")
	}
}

func TestScalaProfileCertifiesPackageTwoRecovery(t *testing.T) {
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := ScalaLanguage()
	if !lang.CompactStrategy2ErrorRegionCertified || !lang.CompactMissingTokenInsertionCertified ||
		!lang.CompactS5EOFMissingInsertionCertified ||
		!lang.CompactFaithfulS5RecoveryCertified ||
		!lang.CompactPrimaryAcceptanceDerivationCertified ||
		!lang.CompactAcceptanceStructuralElectionCertified ||
		!lang.CompactRecoveryTrailingLineageRetirementCertified ||
		!lang.CompactRecoveryErrorModeKeywordCaptureCertified {
		t.Fatal("the Scala profile did not attach its package-two recovery capabilities")
	}
}

func TestBuiltinGoMissingVersionTurnsAttachByName(t *testing.T) {
	lang := GoLanguage()
	if !lang.RecoveryMissingVersionTurnsCertified {
		t.Fatal("built-in Go artifact lacks its missing-version scheduling grant")
	}
	other := &gotreesitter.Language{Name: "custom"}
	if attachBuiltinLanguageRuntimeProfile("custom", other) || other.RecoveryMissingVersionTurnsCertified {
		t.Fatal("the Go profile granted missing-version scheduling to a custom language")
	}
}

func TestBuiltinRuntimeProfilesStayNarrow(t *testing.T) {
	// 42 = the prior 41 plus the F# unary-wrapper materialization profile.
	// Bash, Haskell, and JavaScript reuse their existing exact profiles.
	// Exact blob identity bounds each profile capability.
	// D and Groovy's retry ceilings
	// moved out of parser-core name switches and onto exact-blob profiles. The
	// prior gomod and C additions moved hardcoded compat-tier behavior to profiles.
	// Meson and Enforce add two exact-blob retry policies. JavaScript adds one
	// exact-blob automatic-forest memory allowance. The
	// repetition-conflict helpers were retired in favor of certified
	// ConflictPolicies rows here. AWK and Uxntal add exact-blob automatic forest
	// routing; KDL reuses its existing retry profile. dot's helper was
	// retired outright, not migrated (see the "NOTE on dot" comment above
	// the gomod entry), so it does not add a map entry. Crystal and Matlab add
	// exact-blob external-scanner repeat suppression while retaining the full
	// 44 = the prior 42 plus two new A3 certification-workstream entries:
	// Perl and Ada (spec.campaign.v7, finding
	// tied-election-family-compact-retirement). Apex and Python add compact
	// primary-acceptance-derivation certification (Python also adds
	// converged-split-drop certification); Perl and Ada are new entries
	// carrying both certifications. Kotlin's pre-existing entry stays
	// unchanged in count (it does not add a map entry) but later adds both
	// grants in place once selectCompactAcceptanceDerivation's materiality
	// gate resolves the interaction that withheld them; see the "kotlin" map
	// entry comment.
	// 45 = the prior 44 plus the B3 stage S3 html entry, certifying native
	// strategy-2 error-region recovery for the html_erroneous_end_tag class
	// (spec.compact-recovery-ownership.v1).
	// 46 = the prior 45 plus the powershell entry, declaring backtick as the
	// language's line-continuation escape byte (Language.LineContinuationEscapeByte)
	// so bytesAreParserPadding classifies backtick+newline as padding, matching
	// the C oracle (spore.2026-08-02.birch-g.powershell-bisect).
	// 47 = the prior 46 plus the V entry. It reuses a certified clean wide
	// result during the accepted-error retry ladder.
	// 48 = the prior 47 plus the ql entry. It carries a
	// ConflictPolicyDeclaredReduceReduceHighestSymbol row that resolves the
	// grammar-declared simpleId/className reduce-reduce conflict to the
	// C-native reading during dispatch.
	// 49 = the prior 48 plus the JSDoc entry. It certifies exact internal-DFA
	// skipped-prefix evidence for compact reduce tiling.
	// 50 = the prior 49 plus the Markdown inline entry. It certifies four exact
	// compact conflict rows while production parsing ignores them.
	// 51 = the prior 50 plus the irreducible YAML recover_eof EOF-root entry.
	// 52 = the prior 51 plus the certified Scala package-two recovery entry.
	// 54 = the prior 52 plus the Q0 PHP and Markdown entries (#1311). Both
	// certify one legacy version per merge key; PHP also admits
	// mixed-representation merges.
	// 55 = the prior 54 plus the Q0 Kconfig entry (#1311), which certifies one
	// legacy version per merge key.
	// 56 = the prior 55 plus the Agda entry. It excludes one exact-blob
	// ConflictPolicyRepetitionShift row (state 4039, lookahead id) so the parser
	// reduces there as the C runtime does, which never takes repetition shifts.
	// 57 = the prior 56 plus TypeScript's exact-blob C version-order grant.
	// TypeScript's locked four-file fresh set improves from 3/4 to 4/4.
	// 58 adds SQL's exact-blob missing-version-turn certification (2/4 to 4/4).
	// 59 retains Gleam's exact-blob negative-integer grant from main.
	// The locked fresh set improves from 3/4 to 4/4 without counter changes.
	// 60 adds TSX's exact-blob conflict-action priority grant (#1443),
	// without enabling its broader packed-reduction transaction ordering.
	// 57 = the prior 60 less YAML's recover_eof entry and the empty HTTP and
	// Robot entries.
	if got, want := len(builtinLanguageRuntimeProfiles), 57; got != want {
		t.Fatalf("builtinLanguageRuntimeProfiles has %d entries, want %d", got, want)
	}
	lang := &gotreesitter.Language{ExternalScanner: KotlinExternalScanner{}}
	if attachBuiltinLanguageRuntimeProfile("custom", lang) {
		t.Fatal("unknown runtime profile reported an attachment")
	}
	if got := lang.ExternalScannerFullParseRetryPolicy; got != gotreesitter.ExternalScannerFullParseRetryDefault {
		t.Fatalf("unknown runtime profile changed policy to %d", got)
	}
	if got := lang.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
		t.Fatalf("unknown runtime profile changed accepted-error retry profile to %+v", got)
	}
	if got := lang.AutomaticForestMemoryAllowanceBytes; got != 0 {
		t.Fatalf("unknown runtime profile changed automatic forest allowance to %d", got)
	}
	if lang.FullParseArenaDensityCapEnabled {
		t.Fatal("unknown runtime profile enabled the full-parse arena density cap")
	}
	if lang.CompactConvergedReductionSplitDropsCertified {
		t.Fatal("unknown runtime profile enabled compact converged-split drops")
	}
	if lang.CompactPrimaryAcceptanceDerivationCertified {
		t.Fatal("unknown runtime profile enabled compact primary derivation selection")
	}
	if lang.CompactLexerSkippedPrefixTilingCertified {
		t.Fatal("unknown runtime profile enabled compact lexer skipped-prefix tiling")
	}
	if lang.ExactStackNodeEquivalenceCertified {
		t.Fatal("unknown runtime profile enabled exact stack-node equivalence")
	}
	if lang.CompactPackedGSSVersionOrderCertified {
		t.Fatal("unknown runtime profile enabled packed GSS version ordering")
	}
	if lang.CompactRecoveryPlainFirstCertified {
		t.Fatal("unknown runtime profile enabled compact plain-first recovery")
	}
	if lang.CompactRecoveryTrailingLineageRetirementCertified {
		t.Fatal("unknown runtime profile enabled compact trailing-lineage retirement")
	}
	if lang.CompactRecoveryErrorModeKeywordCaptureCertified {
		t.Fatal("unknown runtime profile enabled compact error-mode keyword capture")
	}
	if lang.LineContinuationEscapeByte != 0 {
		t.Fatalf("unknown runtime profile set a line-continuation escape byte: %q", lang.LineContinuationEscapeByte)
	}
}

func TestRuntimeProfileSymbolRejectsDuplicateNames(t *testing.T) {
	lang := &gotreesitter.Language{SymbolNames: []string{"first", "duplicate", "duplicate"}}
	if symbol, ok := runtimeProfileSymbol(lang, "first"); !ok || symbol != 0 {
		t.Fatalf("unique symbol=%d ok=%t, want 0/true", symbol, ok)
	}
	if symbol, ok := runtimeProfileSymbol(lang, "duplicate"); ok {
		t.Fatalf("duplicate symbol=%d ok=true, want fail-closed", symbol)
	}
}

func TestBuiltinObjcExactStackEquivalenceProfileAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := ObjcLanguage()
	if !builtin.ExactStackNodeEquivalenceCertified {
		t.Fatal("exact built-in Objective-C artifact did not receive exact stack-node equivalence")
	}

	custom := &gotreesitter.Language{Name: "objc"}
	AttachLanguageSupport("objc", custom)
	if custom.ExactStackNodeEquivalenceCertified {
		t.Fatal("same-name custom Objective-C grammar enabled exact stack-node equivalence")
	}
}

func TestBuiltinCompactAcceptanceProfilesAttach(t *testing.T) {
	tests := []struct {
		name string
		load func() *gotreesitter.Language
		want func(*gotreesitter.Language) bool
	}{
		{
			name: "meson", load: MesonLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified &&
					lang.CompactAcceptanceStructuralElectionCertified
			},
		},
		// The tied-election family (A3 certification workstream,
		// spec.campaign.v7, finding
		// tied-election-family-compact-retirement): full-corpus field-aware
		// C-oracle verification certifies primary-acceptance-derivation
		// selection for all five languages. Python also certifies C's raw
		// subtree ordering for clean ties. Kotlin's grant lands under
		// selectCompactAcceptanceDerivation's materiality gate
		// (parsercore_phase0_driver.go, compactAcceptanceElectionIsVacuous);
		// see the runtime_profiles.go "kotlin" entry comment.
		{
			name: "python", load: PythonLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified &&
					lang.CompactAcceptanceStructuralElectionCertified &&
					lang.CompactMixedGSSMergeCertified
			},
		},
		{
			name: "apex", load: ApexLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified
			},
		},
		{
			name: "perl", load: PerlLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified
			},
		},
		{
			name: "ada", load: AdaLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified
			},
		},
		{
			name: "kotlin", load: KotlinLanguage,
			want: func(lang *gotreesitter.Language) bool {
				return lang.CompactPrimaryAcceptanceDerivationCertified
			},
		},
	}

	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if lang := tt.load(); !tt.want(lang) {
				t.Fatal("exact built-in artifact did not receive compact acceptance certification")
			}

			custom := &gotreesitter.Language{Name: tt.name}
			AttachLanguageSupport(tt.name, custom)
			if tt.want(custom) {
				t.Fatal("same-name custom grammar enabled compact acceptance certification")
			}
		})
	}
}

func TestBuiltinJsdocLexerSkippedPrefixProfileAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := JsdocLanguage()
	if !builtin.CompactLexerSkippedPrefixTilingCertified {
		t.Fatal("exact built-in JSDoc artifact did not receive skipped-prefix certification")
	}

	custom := &gotreesitter.Language{Name: "jsdoc"}
	AttachLanguageSupport("jsdoc", custom)
	if custom.CompactLexerSkippedPrefixTilingCertified {
		t.Fatal("same-name custom JSDoc grammar enabled skipped-prefix certification")
	}

	attached := &gotreesitter.Language{Name: "jsdoc"}
	if !attachBuiltinLanguageRuntimeProfile("jsdoc", attached) {
		t.Fatal("the JSDoc profile did not attach")
	}
	if !attached.CompactLexerSkippedPrefixTilingCertified {
		t.Fatal("the JSDoc profile did not enable skipped-prefix certification")
	}
}

func TestBuiltinCompactConvergedSplitProfilesAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	tests := []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "bash", load: BashLanguage},
		{name: "erlang", load: ErlangLanguage},
		{name: "go", load: GoLanguage},
		{name: "haskell", load: HaskellLanguage},
		{name: "javascript", load: JavascriptLanguage},
		{name: "python", load: PythonLanguage},
		// A3 certification workstream (spec.campaign.v7, finding
		// tied-election-family-compact-retirement): Perl and Ada certify
		// converged-path split-drop acceptance after full-corpus field-aware
		// C-oracle verification. Apex does not need this certification and
		// stays out of this table. Kotlin's grant stays withheld: review
		// found a compact-only divergence class (an annotated extension
		// property with a getter and a trailing comment); see the
		// runtime_profiles.go "kotlin" entry comment.
		{name: "perl", load: PerlLanguage},
		{name: "ada", load: AdaLanguage},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builtin := test.load()
			if !builtin.CompactConvergedReductionSplitDropsCertified {
				t.Fatal("exact built-in artifact did not receive compact converged-split certification")
			}

			custom := &gotreesitter.Language{Name: test.name}
			AttachLanguageSupport(test.name, custom)
			if custom.CompactConvergedReductionSplitDropsCertified {
				t.Fatal("same-name custom grammar enabled compact converged-split drops")
			}

			attached := &gotreesitter.Language{Name: test.name}
			if !attachBuiltinLanguageRuntimeProfile(test.name, attached) {
				t.Fatal("the profile did not attach")
			}
			if !attached.CompactConvergedReductionSplitDropsCertified {
				t.Fatal("the profile did not enable compact converged-split drops")
			}
		})
	}
}

func TestBuiltinCompactPackedGSSVersionOrderProfileAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := ErlangLanguage()
	if !builtin.CompactPackedGSSVersionOrderCertified {
		t.Fatal("exact Erlang artifact did not receive packed GSS version-order certification")
	}

	custom := &gotreesitter.Language{Name: "erlang"}
	AttachLanguageSupport("erlang", custom)
	if custom.CompactPackedGSSVersionOrderCertified {
		t.Fatal("same-name custom grammar enabled packed GSS version ordering")
	}
}

func TestBuiltinErlangPackedGSSVersionOrderIssue984Tree(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := ErlangLanguage()
	fresh, err := LoadLanguage("erlang", BlobByName("erlang"))
	if err != nil {
		t.Fatalf("load fresh Erlang language: %v", err)
	}
	for _, test := range []struct {
		name     string
		language *gotreesitter.Language
	}{
		{name: "fresh_decode", language: fresh},
		{name: "cached_builtin", language: builtin},
	} {
		t.Run(test.name, func(t *testing.T) {
			parser := gotreesitter.NewParser(test.language)
			parser.SetAdmissionCandidateRoute(false)
			tree, err := parser.Parse([]byte("000\"0A!A \"A\"=0:A0!)A\"0%0000"))
			if err != nil {
				t.Fatal(err)
			}
			defer tree.Release()
			if got, want := tree.RootNode().SExpr(test.language), `(source_file (integer) (string) (var) (string) (integer) (comment))`; got != want {
				t.Fatalf("Erlang issue #984 tree = %s, want %s", got, want)
			}
		})
	}
}

func BenchmarkBuiltinErlangCompactIssue984(b *testing.B) {
	language := ErlangLanguage()
	parser := gotreesitter.NewParser(language)
	parser.SetAdmissionCandidateRoute(true)
	source := []byte("000\"0A!A \"A\"=0:A0!)A\"0%0000")
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		tree, err := parser.Parse(source)
		if err != nil {
			b.Fatal(err)
		}
		tree.Release()
	}
}

func TestBuiltinErlangPackedGSSVersionOrderPreservesLinearRecovery(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	const source = "f(X) ->\n" +
		"  try 1 of\n" +
		"    X -> ok;\n" +
		"    Y when Z -> ok\n" +
		"  catch\n" +
		"    Pattern -> ok;\n" +
		"    error:Pattern -> ok;\n" +
		"    throw:Pattern:Stack -> ok;\n" +
		"    exit:Complex + Pattern:Stack -> ok\n" +
		"  after\n" +
		"    ok\n" +
		"  end,\n" +
		"  try _\n" +
		"  catch\n" +
		"  end.\n"
	builtin := ErlangLanguage()
	linear, err := LoadLanguage("erlang", BlobByName("erlang"))
	if err != nil {
		t.Fatalf("load fresh Erlang language: %v", err)
	}
	linear.CompactPackedGSSVersionOrderCertified = false
	var sexprs [2]string
	for i, language := range []*gotreesitter.Language{builtin, linear} {
		parser := gotreesitter.NewParser(language)
		tree, err := parser.Parse([]byte(source))
		if err != nil {
			t.Fatalf("route %d parse: %v", i, err)
		}
		if tree.RootNode().HasError() {
			got := tree.RootNode().SExpr(language)
			tree.Release()
			t.Fatalf("route %d produced an error tree: %s", i, got)
		}
		sexprs[i] = tree.RootNode().SExpr(language)
		tree.Release()
	}
	if sexprs[0] != sexprs[1] {
		t.Fatalf("packed and linear Erlang trees differ:\npacked: %s\nlinear: %s", sexprs[0], sexprs[1])
	}
}

func TestBuiltinOdinArenaDensityProfileAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := OdinLanguage()
	if !builtin.FullParseArenaDensityCapEnabled {
		t.Fatal("exact built-in Odin artifact did not receive the arena density certification")
	}

	custom := &gotreesitter.Language{Name: "odin"}
	AttachLanguageSupport("odin", custom)
	if custom.FullParseArenaDensityCapEnabled {
		t.Fatal("same-name custom Odin grammar enabled the arena density cap")
	}

	attached := &gotreesitter.Language{Name: "odin"}
	if !attachBuiltinLanguageRuntimeProfile("odin", attached) {
		t.Fatal("the Odin profile did not attach")
	}
	if !attached.FullParseArenaDensityCapEnabled {
		t.Fatal("the Odin profile did not enable the arena density cap")
	}
}

func TestBuiltinAutomaticForestProfilesAttach(t *testing.T) {
	tests := []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "awk", load: AwkLanguage},
		{name: "kdl", load: KdlLanguage},
		{name: "uxntal", load: UxntalLanguage},
	}

	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builtin := tt.load()
			if !builtin.AutomaticForestEnabledByDefault {
				t.Fatal("exact built-in artifact did not receive automatic forest certification")
			}

			adapted := &gotreesitter.Language{Name: tt.name}
			AttachLanguageSupport(tt.name, adapted)
			if adapted.AutomaticForestEnabledByDefault {
				t.Fatal("same-name adapted language enabled automatic forest routing")
			}
		})
	}
}

func TestEnforceDuplicateWideProfileAttaches(t *testing.T) {
	want := gotreesitter.FullParseAcceptedErrorRetryProfile{
		ReuseCleanWideForWideRetry:   true,
		ReuseCleanWideMinSourceBytes: 128 * 1024,
	}

	adapted := &gotreesitter.Language{Name: "enforce"}
	AttachLanguageSupport("enforce", adapted)
	if got := adapted.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
		t.Fatalf("adapted Enforce language attached retry profile %+v", got)
	}

	attached := &gotreesitter.Language{Name: "enforce"}
	if !attachBuiltinLanguageRuntimeProfile("enforce", attached) {
		t.Fatal("the Enforce profile did not attach")
	}
	if got := attached.FullParseAcceptedErrorRetryProfile; got != want {
		t.Fatalf("Enforce retry profile = %+v, want %+v", got, want)
	}
}

func TestBuiltinHaskellConflictPolicyAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	lang := HaskellLanguage()
	want := map[gotreesitter.StateID]struct {
		lookahead gotreesitter.Symbol
		reduce    gotreesitter.Symbol
	}{
		9609:  {lookahead: gotreesitter.ConflictPolicyAnyLookahead, reduce: 518},
		10984: {lookahead: gotreesitter.ConflictPolicyAnyLookahead, reduce: 516},
		11192: {lookahead: 4, reduce: 500},
	}
	seen := make(map[gotreesitter.StateID]bool, len(want))
	for _, policy := range lang.ConflictPolicies {
		expected, ok := want[policy.State]
		if !ok {
			continue
		}
		seen[policy.State] = true
		if policy.Lookahead != expected.lookahead {
			t.Fatalf("Haskell conflict policy at state %d has lookahead %d, want %d", policy.State, policy.Lookahead, expected.lookahead)
		}
		if policy.Kind != gotreesitter.ConflictPolicyRepetitionReduce ||
			len(policy.ReduceSymbols) != 1 || policy.ReduceSymbols[0] != expected.reduce {
			t.Fatalf("Haskell conflict policy at state %d = %+v, want repetition reduce over symbol %d", policy.State, policy, expected.reduce)
		}
	}
	for state := range want {
		if !seen[state] {
			t.Fatalf("Haskell conflict policy for state %d was not attached", state)
		}
	}
}

func TestBuiltinHaskellConflictPolicyAttachesOnce(t *testing.T) {
	lang := &gotreesitter.Language{Name: "haskell"}
	if !attachBuiltinLanguageRuntimeProfile("haskell", lang) {
		t.Fatal("the Haskell profile did not attach")
	}
	if got := len(lang.ConflictPolicies); got != 3 {
		t.Fatalf("Haskell conflict policies = %d, want 3", got)
	}
	if attachBuiltinLanguageRuntimeProfile("haskell", lang) {
		t.Fatal("reattaching the same Haskell profile reported a change")
	}
	if got := len(lang.ConflictPolicies); got != 3 {
		t.Fatalf("reattached Haskell conflict policies = %d, want 3", got)
	}
}

// TestBuiltinDartConflictPoliciesAttach covers the three certified dart
// repeat-boundary rows (enum bodies, extension bodies, top-level declaration
// lists) that replaced the hardcoded dartRepetitionShiftConflictChoice helper
// through the generic policy path.
func TestBuiltinDartConflictPoliciesAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	lang := DartLanguage()
	// Pin the reduce symbol by name as well as by identifier. A grammar bump
	// renumbers symbols, so the identifier alone cannot show that a row still
	// covers the same repeat boundary.
	wantStates := map[gotreesitter.StateID]gotreesitter.Symbol{601: 511, 616: 516, 574: 473}
	wantNames := map[gotreesitter.StateID]string{
		601: "enum_body_repeat2",
		616: "extension_body_repeat1",
		574: "program_repeat4",
	}
	gotStates := map[gotreesitter.StateID]bool{}
	for _, policy := range lang.ConflictPolicies {
		wantReduce, known := wantStates[policy.State]
		if !known {
			continue
		}
		gotStates[policy.State] = true
		if policy.Lookahead != gotreesitter.ConflictPolicyAnyLookahead ||
			policy.Kind != gotreesitter.ConflictPolicyRepetitionShift ||
			len(policy.ReduceSymbols) != 1 || policy.ReduceSymbols[0] != wantReduce {
			t.Fatalf("dart conflict policy at state %d = %+v, want repetition-shift over reduce symbol %d", policy.State, policy, wantReduce)
		}
		gotName := ""
		if idx := int(policy.ReduceSymbols[0]); idx >= 0 && idx < len(lang.SymbolNames) {
			gotName = lang.SymbolNames[idx]
		}
		if gotName != wantNames[policy.State] {
			t.Fatalf("dart conflict policy at state %d reduces symbol %d (%q), want %q",
				policy.State, policy.ReduceSymbols[0], gotName, wantNames[policy.State])
		}
	}
	for state := range wantStates {
		if !gotStates[state] {
			t.Fatalf("dart conflict policy for state %d was not attached", state)
		}
	}
	// The blob carries its own table-derived precedence rows. The runtime
	// profile must add exactly the three certified repeat-boundary rows on top
	// of them, and must not drop or duplicate a derived row.
	blob := BlobByName("dart")
	if len(blob) == 0 {
		t.Fatal("BlobByName(dart) returned no data")
	}
	raw, err := gotreesitter.LoadLanguage(blob)
	if err != nil {
		t.Fatalf("LoadLanguage(dart blob): %v", err)
	}
	if got, want := len(lang.ConflictPolicies), len(raw.ConflictPolicies)+len(wantStates); got != want {
		t.Fatalf("dart ConflictPolicies = %d rows, want %d (%d table-derived + %d certified)",
			got, want, len(raw.ConflictPolicies), len(wantStates))
	}
}

func TestBuiltinAgdaConflictPolicyExclusion(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	lang, err := LoadLanguage("agda", BlobByName("agda"))
	if err != nil {
		t.Fatalf("load Agda language: %v", err)
	}
	neighboringRowRetained := false
	for _, policy := range lang.ConflictPolicies {
		if policy.State != 4039 || policy.Kind != gotreesitter.ConflictPolicyRepetitionShift {
			continue
		}
		if policy.Lookahead == 1 && len(policy.ReduceSymbols) == 1 && policy.ReduceSymbols[0] == 173 {
			t.Fatalf("unsafe Agda repetition-shift policy remains attached: %+v", policy)
		}
		if policy.Lookahead == 8 {
			neighboringRowRetained = true
		}
	}
	if !neighboringRowRetained {
		t.Fatal("Agda exclusion removed repetition-shift policies beyond the id lookahead row")
	}
}

func TestBuiltinMarkdownInlineConflictPolicyAttaches(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	builtin := MarkdownInlineLanguage()
	if got := len(builtin.ConflictPolicies); got != 4 {
		t.Fatalf("Markdown inline conflict policies = %d, want 4", got)
	}
	fold := builtin.ConflictPolicies[0]
	if fold.State != 18 || fold.Lookahead != 56 || fold.Kind != gotreesitter.ConflictPolicyRepetitionReduce ||
		!fold.CompactOnly || fold.CompactMinFrontierHeaders != 2 || len(fold.ReduceSymbols) != 1 || fold.ReduceSymbols[0] != 107 {
		t.Fatalf("Markdown inline fold policy = %+v, want state 18 reduce symbol 107 for lookahead 56", fold)
	}
	wordFold := builtin.ConflictPolicies[1]
	if wordFold.State != 18 || wordFold.Lookahead != 50 || wordFold.Kind != gotreesitter.ConflictPolicyRepetitionReduce ||
		!wordFold.CompactOnly || wordFold.CompactMinFrontierHeaders != 2 || len(wordFold.ReduceSymbols) != 1 || wordFold.ReduceSymbols[0] != 107 {
		t.Fatalf("Markdown inline word fold policy = %+v, want state 18 reduce symbol 107 for lookahead 50", wordFold)
	}
	lineEndFold := builtin.ConflictPolicies[2]
	if lineEndFold.State != 18 || lineEndFold.Lookahead != 36 || lineEndFold.Kind != gotreesitter.ConflictPolicyRepetitionReduce ||
		!lineEndFold.CompactOnly || lineEndFold.CompactMinFrontierHeaders != 2 || len(lineEndFold.ReduceSymbols) != 1 || lineEndFold.ReduceSymbols[0] != 107 {
		t.Fatalf("Markdown inline line-end fold policy = %+v, want state 18 reduce symbol 107 for lookahead 36", lineEndFold)
	}
	policy := builtin.ConflictPolicies[3]
	if policy.State != 209 || policy.Lookahead != 50 || policy.Kind != gotreesitter.ConflictPolicyShift ||
		!policy.CompactOnly || len(policy.ReduceSymbols) != 1 || policy.ReduceSymbols[0] != 104 {
		t.Fatalf("Markdown inline conflict policy = %+v, want state 209 shift over symbol 104 for lookahead 50", policy)
	}

	custom := &gotreesitter.Language{Name: "markdown_inline"}
	AttachLanguageSupport("markdown_inline", custom)
	if len(custom.ConflictPolicies) != 0 {
		t.Fatal("same-name custom grammar received the Markdown inline conflict policy")
	}
}

// TestBuiltinGomodConflictPoliciesAttach covers the certified go.mod
// require-list rows that replaced gomodRepetitionShiftConflictChoice
// through the generic policy path.
func TestBuiltinGomodConflictPoliciesAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	lang := GomodLanguage()
	var state3, state37 int
	for _, policy := range lang.ConflictPolicies {
		if policy.Kind != gotreesitter.ConflictPolicyRepetitionShift {
			t.Fatalf("gomod conflict policy = %+v, want repetition-shift kind", policy)
		}
		if policy.Lookahead != gotreesitter.ConflictPolicyAnyLookahead {
			t.Fatalf("gomod conflict policy = %+v, want wildcard lookahead", policy)
		}
		switch policy.State {
		case 3:
			if len(policy.ReduceSymbols) != 1 || policy.ReduceSymbols[0] != 50 {
				t.Fatalf("gomod state-3 conflict policy = %+v, want source_file_repeat1 (symbol 50)", policy)
			}
			state3++
		case 37:
			if len(policy.ReduceSymbols) != 1 || policy.ReduceSymbols[0] != 52 {
				t.Fatalf("gomod state-37 conflict policy = %+v, want require_directive_repeat1 (symbol 52)", policy)
			}
			state37++
		default:
			t.Fatalf("unexpected gomod conflict policy state %d", policy.State)
		}
	}
	if state3 != 1 {
		t.Fatalf("gomod state-3 conflict policies = %d, want 1", state3)
	}
	if state37 != 1 {
		t.Fatalf("gomod state-37 conflict policies = %d, want 1", state37)
	}
}

// TestBuiltinDotHasNoConflictPolicies documents that dot's retired
// compat-tier helper (dotRepetitionShiftConflictChoice) was NOT migrated to a
// certified ConflictPolicy: dot already relied on the engine-wide C
// repetition-skip fold (it is not in cRepetitionSkipOptOut), and A/B testing
// found the helper's repetition-shift preference grows the LR parse-stack
// depth O(n) with statement count for no behavioral benefit over that
// already-flat fold (see the "NOTE on dot" comment in runtime_profiles.go).
func TestBuiltinDotHasNoConflictPolicies(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	if got := len(DotLanguage().ConflictPolicies); got != 0 {
		t.Fatalf("dot ConflictPolicies = %d rows, want 0 (helper retired, not migrated)", got)
	}
}

// TestBuiltinCConflictPolicyAttachesWildcard covers the certified C
// translation_unit_repeat1/preproc_if_repeat1 fold, the one profile using the
// ConflictPolicyAnyState/AnyLookahead wildcards because the C-faithful rule
// is scoped by reduce-symbol identity alone (thousands of real table rows).
// Replaces cRepetitionShiftConflictChoice.
func TestBuiltinCConflictPolicyAttachesWildcard(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	lang := CLanguage()
	if got, want := len(lang.ConflictPolicies), 2; got != want {
		t.Fatalf("c ConflictPolicies = %d rows, want %d", got, want)
	}
	policy := lang.ConflictPolicies[0]
	if policy.State != gotreesitter.ConflictPolicyAnyState || policy.Lookahead != gotreesitter.ConflictPolicyAnyLookahead ||
		policy.Kind != gotreesitter.ConflictPolicyRepetitionShift || len(policy.ReduceSymbols) != 2 ||
		policy.ReduceSymbols[0] != 324 || policy.ReduceSymbols[1] != 326 {
		t.Fatalf("c conflict policy = %+v, want wildcard repetition-shift over translation_unit_repeat1(324)/preproc_if_repeat1(326)", policy)
	}
	// C is held out of the engine-wide repetition fold
	// (cRepetitionSkipConflictChoice), so it carves out the symbols that are
	// proven safe one row at a time, the way haskell does. This row folds
	// aux_sym_enumerator_list_repeat1 and stops a three-enumerator list from
	// inventing a comma; see grammars/c_enum_missing_comma_pin_test.go.
	fold := lang.ConflictPolicies[1]
	if fold.State != gotreesitter.ConflictPolicyAnyState || fold.Lookahead != gotreesitter.ConflictPolicyAnyLookahead ||
		fold.Kind != gotreesitter.ConflictPolicyRepetitionReduce || len(fold.ReduceSymbols) != 1 ||
		fold.ReduceSymbols[0] != 343 {
		t.Fatalf("c conflict policy[1] = %+v, want wildcard repetition-reduce over enumerator_list_repeat1(343)", fold)
	}
}

func TestBuiltinCompleteAcceptedErrorRetryProfilesAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	tests := []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "bash", load: BashLanguage},
		{name: "c", load: CLanguage},
		{name: "caddy", load: CaddyLanguage},
		{name: "c_sharp", load: CSharpLanguage},
		{name: "cpp", load: CppLanguage},
		{name: "haxe", load: HaxeLanguage},
		{name: "kdl", load: KdlLanguage},
		{name: "meson", load: MesonLanguage},
		{name: "odin", load: OdinLanguage},
		{name: "rego", load: RegoLanguage},
		{name: "scss", load: ScssLanguage},
		{name: "swift", load: SwiftLanguage},
		{name: "tcl", load: TclLanguage},
		{name: "v", load: VLanguage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := tt.load().FullParseAcceptedErrorRetryProfile
			if tt.name == "v" {
				if profile.SkipInitialCompleteAcceptedErrorMergeRetry || !profile.ReuseCleanWideForWideRetry ||
					profile.SkipFreshCompleteAcceptedErrorRetry || profile.SkipCompleteAcceptedErrorRetry {
					t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want clean-wide reuse certification", profile)
				}
			} else if !profile.SkipCompleteAcceptedErrorRetry {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want skip-complete certification", profile)
			}
			if tt.name == "swift" && profile.SkipCompleteMaxEntryScratchPeak != 3*64*1024 {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want first-growth entry-scratch ceiling", profile)
			}
			if tt.name == "meson" && profile.SkipCompleteMinSourceBytes != 2*1024 {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want %d-byte skip minimum", profile, 2*1024)
			}
			if tt.name == "v" && profile.ReuseCleanWideMinSourceBytes != vAcceptedErrorRetryMinSourceBytes {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want %d-byte reuse minimum", profile, vAcceptedErrorRetryMinSourceBytes)
			}
			if tt.name == "c_sharp" && (profile.SkipCompleteMaxEntryScratchPeak != csharpAcceptedErrorRetryMaxEntryScratchPeak ||
				profile.FreshErrorNoStacksRetryMaxStacks != csharpFreshErrorNoStacksRetryMaxStacks ||
				profile.GSSConvergenceAcceptedErrorMergePerKey != csharpGSSConvergenceErrorMergePerKey ||
				!profile.SkipInitialCompleteAcceptedErrorMergeRetry) {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want bounded C# retry profile", profile)
			}
			if tt.name == "tcl" && profile.FreshErrorNoStacksMaxPasses != 1 {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want one fresh no-stacks retry", profile)
			}
		})
	}
}

func TestBuiltinBoundedAcceptedErrorRetryProfilesAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })

	tests := []struct {
		name          string
		load          func() *gotreesitter.Language
		want          gotreesitter.FullParseAcceptedErrorRetryProfile
		wantNoScanner bool
	}{
		{
			name: "asm",
			load: AsmLanguage,
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				MinSourceBytes:      11 * 1024,
				InitialStackCeiling: 8,
			},
		},
		{
			name: "java",
			load: JavaLanguage,
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				MinSourceBytes:      64 * 1024,
				InitialStackCeiling: 14,
			},
			wantNoScanner: true,
		},
		{
			name: "groovy",
			load: GroovyLanguage,
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				MinSourceBytes:      64 * 1024,
				InitialStackCeiling: 2,
			},
			wantNoScanner: true,
		},
		{
			name: "d",
			load: DLanguage,
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				MinSourceBytes:      64 * 1024,
				InitialStackCeiling: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := tt.load()
			if tt.wantNoScanner && lang.ExternalScanner != nil {
				t.Fatalf("%s ExternalScanner != nil; profile must not depend on scanner attachment", tt.name)
			}
			if got := lang.FullParseAcceptedErrorRetryProfile; got != tt.want {
				t.Fatalf("FullParseAcceptedErrorRetryProfile = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestBuiltinExternalScannerRetryProfilesAttachByName(t *testing.T) {
	tests := []struct {
		name    string
		scanner gotreesitter.ExternalScanner
	}{
		{name: "crystal", scanner: CrystalExternalScanner{}},
		{name: "kotlin", scanner: KotlinExternalScanner{}},
		{name: "matlab", scanner: MatlabExternalScanner{}},
		{name: "swift", scanner: SwiftExternalScanner{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := &gotreesitter.Language{Name: tt.name, ExternalScanner: tt.scanner}
			if !attachBuiltinLanguageRuntimeProfile(tt.name, lang) {
				t.Fatal("the profile did not attach")
			}
			if got := lang.ExternalScannerFullParseRetryPolicy; got != gotreesitter.ExternalScannerFullParseRetrySkipRepeat {
				t.Fatalf("profile policy = %d, want skip-repeat", got)
			}
		})
	}
}

func TestBuiltinCSharpRetryProfileAttaches(t *testing.T) {
	lang := &gotreesitter.Language{ExternalScanner: CSharpExternalScanner{}}
	if !attachBuiltinLanguageRuntimeProfile("c_sharp", lang) {
		t.Fatal("the C# profile did not attach")
	}
	if got := lang.ExternalScannerFullParseRetryPolicy; got != gotreesitter.ExternalScannerFullParseRetrySkipRepeat {
		t.Fatalf("certified C# scanner retry policy = %d, want skip-repeat", got)
	}
	want := gotreesitter.FullParseAcceptedErrorRetryProfile{
		SkipCompleteAcceptedErrorRetry:             true,
		SkipCompleteMaxEntryScratchPeak:            csharpAcceptedErrorRetryMaxEntryScratchPeak,
		FreshErrorNoStacksRetryMaxStacks:           csharpFreshErrorNoStacksRetryMaxStacks,
		SkipInitialCompleteAcceptedErrorMergeRetry: true,
		GSSConvergenceAcceptedErrorMergePerKey:     csharpGSSConvergenceErrorMergePerKey,
	}
	if got := lang.FullParseAcceptedErrorRetryProfile; got != want {
		t.Fatalf("certified C# retry profile = %+v, want %+v", got, want)
	}
	if !lang.FullParseGSSConvergenceEnabled {
		t.Fatal("certified C# blob did not enable full-parse GSS convergence")
	}
}

// TestBuiltinLegacyMergeAdmissionProfilesAttach pins the Q0 per-language
// merge admission entries (#1311).
func TestBuiltinLegacyMergeAdmissionProfilesAttach(t *testing.T) {
	tests := []struct {
		name        string
		policy      gotreesitter.LegacyMergeAdmissionPolicy
		convergence bool
	}{
		{name: "elixir", policy: gotreesitter.LegacyMergeAdmitDynamicPrecedence, convergence: true},
		{name: "php", policy: gotreesitter.LegacyMergeAdmitMixedRepresentation, convergence: true},
		{name: "markdown", convergence: true},
		{name: "kconfig", convergence: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exact := &gotreesitter.Language{}
			if !attachBuiltinLanguageRuntimeProfile(test.name, exact) {
				t.Fatalf("the %s profile did not attach", test.name)
			}
			if exact.LegacyMergeAdmission != test.policy {
				t.Fatalf("%s merge admission = %d, want %d", test.name, exact.LegacyMergeAdmission, test.policy)
			}
			if exact.FullParseGSSConvergenceEnabled != test.convergence {
				t.Fatalf("%s convergence = %t, want %t", test.name, exact.FullParseGSSConvergenceEnabled, test.convergence)
			}
		})
	}
}

func TestBuiltinRecoveryVersionOrderAttaches(t *testing.T) {
	for _, name := range []string{"c", "cpp"} {
		t.Run(name, func(t *testing.T) {
			attached := &gotreesitter.Language{}
			attachBuiltinLanguageRuntimeProfile(name, attached)
			if !attached.RecoveryStackVersionOrderEnabled {
				t.Fatal("the profile did not enable recovery version ordering")
			}
		})
	}
}

func TestResidualRetryProfilesAttach(t *testing.T) {
	tests := []struct {
		name string
		want gotreesitter.FullParseAcceptedErrorRetryProfile
	}{
		{
			name: "groovy",
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{MinSourceBytes: 64 * 1024, InitialStackCeiling: 2},
		},
		{
			name: "d",
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{MinSourceBytes: 64 * 1024, InitialStackCeiling: 3},
		},
		{
			name: "c_sharp",
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				SkipCompleteAcceptedErrorRetry:             true,
				SkipCompleteMaxEntryScratchPeak:            csharpAcceptedErrorRetryMaxEntryScratchPeak,
				FreshErrorNoStacksRetryMaxStacks:           csharpFreshErrorNoStacksRetryMaxStacks,
				SkipInitialCompleteAcceptedErrorMergeRetry: true,
				GSSConvergenceAcceptedErrorMergePerKey:     csharpGSSConvergenceErrorMergePerKey,
			},
		},
		{
			name: "meson",
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				SkipCompleteAcceptedErrorRetry: true,
				SkipCompleteMinSourceBytes:     2 * 1024,
			},
		},
		{
			name: "v",
			want: gotreesitter.FullParseAcceptedErrorRetryProfile{
				ReuseCleanWideForWideRetry:   true,
				ReuseCleanWideMinSourceBytes: vAcceptedErrorRetryMinSourceBytes,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapted := &gotreesitter.Language{Name: tt.name}
			AttachLanguageSupport(tt.name, adapted)
			if got := adapted.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
				t.Fatalf("adapted same-name language attached retry profile %+v", got)
			}

			attached := &gotreesitter.Language{Name: tt.name}
			if !attachBuiltinLanguageRuntimeProfile(tt.name, attached) {
				t.Fatal("the profile did not attach")
			}
			if got := attached.FullParseAcceptedErrorRetryProfile; got != tt.want {
				t.Fatalf("profile retry profile = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCSharpNativeResultCompatibilityAttaches(t *testing.T) {
	want := gotreesitter.ResultCompatibilityCSharpNativeNotNull |
		gotreesitter.ResultCompatibilityCSharpNativeUnicodeIdentifiers |
		gotreesitter.ResultCompatibilityCSharpNativeScopedLambdaStatements |
		gotreesitter.ResultCompatibilityCSharpNativeScopedLambdaBlocks |
		gotreesitter.ResultCompatibilityCSharpNativeQueryExpressions |
		gotreesitter.ResultCompatibilityNativeRecoveredStructure

	adapted := &gotreesitter.Language{Name: "c_sharp"}
	AttachLanguageSupport("c_sharp", adapted)
	if got := adapted.NativeResultCompatibility; got != 0 {
		t.Fatalf("AttachLanguageSupport certified native result compatibility: %#x", got)
	}

	attached := &gotreesitter.Language{Name: "c_sharp"}
	if !attachBuiltinLanguageRuntimeProfile("c_sharp", attached) {
		t.Fatal("the C# profile did not attach")
	}
	if got := attached.NativeResultCompatibility; got != want {
		t.Fatalf("profile NativeResultCompatibility = %#x, want %#x", got, want)
	}

	loaded, err := LoadLanguage("c_sharp", BlobByName("c_sharp"))
	if err != nil {
		t.Fatalf("LoadLanguage(c_sharp): %v", err)
	}
	if got := loaded.NativeResultCompatibility; got != want {
		t.Fatalf("LoadLanguage NativeResultCompatibility = %#x, want %#x", got, want)
	}
	if got := CSharpLanguage().NativeResultCompatibility; got != want {
		t.Fatalf("embedded CSharpLanguage NativeResultCompatibility = %#x, want %#x", got, want)
	}
}

func TestCollapsedChildNativeCapabilityAttaches(t *testing.T) {
	want := gotreesitter.ResultCompatibilityNativeCollapsedChildren
	tests := []struct {
		name string
		load func() *gotreesitter.Language
	}{
		{name: "apex", load: ApexLanguage},
		{name: "dart", load: DartLanguage},
		{name: "elixir", load: ElixirLanguage},
		{name: "hack", load: HackLanguage},
		{name: "kotlin", load: KotlinLanguage},
		{name: "ruby", load: RubyLanguage},
	}
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapted := &gotreesitter.Language{Name: tt.name}
			AttachLanguageSupport(tt.name, adapted)
			if adapted.NativeResultCompatibility&want != 0 {
				t.Fatal("same-name adapted language attached collapsed-child capability")
			}

			exact := tt.load()
			if exact.NativeResultCompatibility&want == 0 {
				t.Fatal("exact built-in artifact omitted collapsed-child capability")
			}
		})
	}
}

func TestAttachLanguageSupportDoesNotAttachRuntimeProfiles(t *testing.T) {
	base := KotlinLanguage()
	lang := &gotreesitter.Language{
		Name:            base.Name,
		ExternalSymbols: append([]gotreesitter.Symbol(nil), base.ExternalSymbols...),
		SymbolNames:     append([]string(nil), base.SymbolNames...),
	}
	if !AttachLanguageSupport("kotlin", lang) {
		t.Fatal("AttachLanguageSupport(kotlin) did not attach scanner support")
	}
	if got := lang.ExternalScannerFullParseRetryPolicy; got != gotreesitter.ExternalScannerFullParseRetryDefault {
		t.Fatalf("AttachLanguageSupport certified policy without blob identity: %d", got)
	}

	java := &gotreesitter.Language{Name: "java"}
	AttachLanguageSupport("java", java)
	if got := java.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
		t.Fatalf("AttachLanguageSupport certified Java retry profile without blob identity: %+v", got)
	}

	rego := &gotreesitter.Language{Name: "rego"}
	AttachLanguageSupport("rego", rego)
	if got := rego.FullParseAcceptedErrorRetryProfile; got != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) {
		t.Fatalf("AttachLanguageSupport certified Rego retry profile without blob identity: %+v", got)
	}
}

func TestNativeUnaryWrapperFlatteningProfileCensus(t *testing.T) {
	configuredLanguages := make(map[string]int)
	for name, profile := range builtinLanguageRuntimeProfiles {
		if len(profile.nativeUnaryWrapperFlattening) == 0 {
			continue
		}
		configuredLanguages[name] = len(profile.nativeUnaryWrapperFlattening)
	}
	if len(configuredLanguages) != 1 || configuredLanguages["fsharp"] != 1 {
		t.Fatalf("native unary-wrapper profile census = %v, want fsharp:1", configuredLanguages)
	}

	language := FsharpLanguage()
	if len(language.NativeUnaryWrapperFlattening) != 1 {
		t.Fatalf("F# native unary-wrapper rules = %d, want 1", len(language.NativeUnaryWrapperFlattening))
	}
	rule := language.NativeUnaryWrapperFlattening[0]
	for symbol, want := range map[gotreesitter.Symbol]string{
		rule.PublicParent: "long_identifier_or_op",
		rule.Wrapper:      "long_identifier",
		rule.Leaf:         "identifier",
	} {
		if int(symbol) >= len(language.SymbolNames) || language.SymbolNames[symbol] != want {
			t.Fatalf("F# unary-wrapper symbol %d = %q, want %q", symbol, language.SymbolNames[symbol], want)
		}
	}
	if got, want := rule.WrapperPreGotoState, gotreesitter.StateID(4258); got != want {
		t.Fatalf("F# wrapper pre-goto state = %d, want %d", got, want)
	}
}

func TestSwiftErrorModeExternalSymbolsAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	lang := SwiftLanguage()
	row := lang.ExternalLexStates[lang.LexModes[0].ExternalLexState]
	if len(row) != 34 {
		t.Fatalf("Swift ERROR row has %d symbols, want the locked C's 34", len(row))
	}
	for i, valid := range row {
		if !valid {
			t.Fatalf("Swift ERROR row disables external symbol %d", i)
		}
	}
	rows := len(lang.ExternalLexStates)
	attachBuiltinLanguageRuntimeProfile("swift", lang)
	if len(lang.ExternalLexStates) != rows {
		t.Fatal("reattaching the profile grew the scanner table")
	}

	modes := []gotreesitter.LexMode{{ExternalLexState: 1}, {ExternalLexState: 1}}
	states := [][]bool{{false, false}, {true, false}}
	custom := &gotreesitter.Language{Name: "swift", LexModes: modes, ExternalLexStates: states, ExternalSymbols: []gotreesitter.Symbol{1, 2}}
	attachBuiltinLanguageRuntimeProfile("swift", custom)
	if modes[0].ExternalLexState != 1 || states[1][1] || custom.LexModes[1].ExternalLexState != 1 {
		t.Fatal("ERROR row changed an ordinary scanner row or shared table")
	}
	if got := len(custom.ExternalLexStates); got != 3 {
		t.Fatalf("scanner rows=%d, want 3", got)
	}
	adapted := &gotreesitter.Language{Name: "swift", LexModes: []gotreesitter.LexMode{{}}, ExternalSymbols: []gotreesitter.Symbol{1}}
	AttachLanguageSupport("swift", adapted)
	if len(adapted.ExternalLexStates) != 0 {
		t.Fatal("same-name adapted grammar acquired the certified ERROR row")
	}
}

func TestRecoveryMissingVersionTurnsAttach(t *testing.T) {
	PurgeEmbeddedLanguageCache()
	t.Cleanup(func() { PurgeEmbeddedLanguageCache() })
	if !Language("sql").RecoveryMissingVersionTurnsCertified {
		t.Fatal("embedded SQL grammar did not acquire recovery turns")
	}
	adapted := &gotreesitter.Language{Name: "sql"}
	AttachLanguageSupport("sql", adapted)
	if adapted.RecoveryMissingVersionTurnsCertified {
		t.Fatal("same-name adapted grammar acquired recovery-turn certification")
	}
}

func TestBuiltinConflictActionVersionOrderAttaches(t *testing.T) {
	for _, tc := range []struct {
		name             string
		load             func() *gotreesitter.Language
		packedOrder      bool
		exactEquivalence bool
	}{{"typescript", TypescriptLanguage, true, true}, {"dart", DartLanguage, true, false}, {"tsx", TsxLanguage, false, true}} {
		name := tc.name
		t.Run(name, func(t *testing.T) {
			builtin := tc.load()
			if builtin == nil || !builtin.ConflictActionVersionOrderCertified || builtin.CompactPackedGSSVersionOrderCertified != tc.packedOrder || builtin.ExactStackNodeEquivalenceCertified != tc.exactEquivalence {
				t.Fatal("exact artifact must certify only its proven conflict and reduction ordering")
			}
			custom := &gotreesitter.Language{Name: name}
			AttachLanguageSupport(name, custom)
			if custom.ConflictActionVersionOrderCertified || custom.CompactPackedGSSVersionOrderCertified || custom.ExactStackNodeEquivalenceCertified {
				t.Fatal("same-name custom artifact received certification")
			}
		})
	}
}
