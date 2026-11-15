package grammarruntime

import (
	"slices"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// builtinLanguageRuntimeProfile contains narrow runtime decisions certified
// against the checked-in grammar/scanner pair. Keep these policies outside the
// parser core: loading a built-in grammar by name attaches its profile, while
// caller-constructed and adapted languages retain conservative zero defaults.
type builtinLanguageRuntimeProfile struct {
	errorModeAllExternalSymbols         bool
	recoveryMissingVersionTurns         bool
	externalScannerCheckpointReuse      bool
	externalScannerFullParseRetry       gotreesitter.ExternalScannerFullParseRetryPolicy
	fullParseAcceptedErrorRetryProfile  gotreesitter.FullParseAcceptedErrorRetryProfile
	automaticForestMemoryAllowance      int64
	automaticForestEnabled              bool
	fullParseArenaDensityCap            bool
	fullParseGSSConvergence             bool
	recoveryStackVersionOrder           bool
	nativeResultCompatibility           gotreesitter.ResultCompatibilityCapability
	nativeUnaryWrapperFlattening        []nativeUnaryWrapperFlatteningProfile
	visibleTerminalAliases              []visibleTerminalAliasProfile
	compactConvergedSplitDrops          bool
	compactPrimaryAcceptDerivation      bool
	compactAcceptanceStructuralElection bool
	compactMixedGSSMerge                bool
	legacyMergeAdmission                gotreesitter.LegacyMergeAdmissionPolicy
	compactLexerSkippedPrefixTiling     bool
	exactStackNodeEquivalence           bool
	compactPackedGSSVersionOrder        bool
	zeroWidthChildPaddingSpan           bool
	conflictActionVersionOrder          bool
	compactStrategy2ErrorRegion         bool
	compactMissingTokenInsertion        bool
	compactS5EOFMissingInsertion        bool
	compactFaithfulS5Recovery           bool
	compactOwnedEOFRecovery             bool
	compactRecoveryTrailingRetirement   bool
	compactRecoveryErrorModeKeyword     bool
	compactRecoveryPlainFirst           bool
	lineContinuationEscapeByte          byte
	conflictPolicies                    []gotreesitter.ConflictPolicy
	conflictPolicyExclusions            []conflictPolicyExclusionProfile
}

type conflictPolicyExclusionProfile struct {
	state         gotreesitter.StateID
	lookahead     gotreesitter.Symbol
	kind          gotreesitter.ConflictPolicyKind
	reduceSymbols []gotreesitter.Symbol
}

type nativeUnaryWrapperFlatteningProfile struct {
	publicParent        string
	wrapper             string
	leaf                string
	wrapperPreGotoState gotreesitter.StateID
}

type visibleTerminalAliasProfile struct {
	internalName string
	publicName   string
}

const (
	csharpAcceptedErrorRetryMaxEntryScratchPeak = 690_365
	csharpFreshErrorNoStacksRetryMaxStacks      = 16
	csharpGSSConvergenceErrorMergePerKey        = 12
	mesonAcceptedErrorRetryMinSourceBytes       = 2 * 1024
	vAcceptedErrorRetryMinSourceBytes           = 128 * 1024
	javascriptAutomaticForestMemoryAllowance    = 128 * 1024 * 1024
)

var builtinLanguageRuntimeProfiles = map[string]builtinLanguageRuntimeProfile{
	// Preserve C's physical version order and source-version conflict priority.
	// Exact link comparisons retain async generic-arrow histories until selection.
	"typescript": {
		exactStackNodeEquivalence:    true,
		compactPackedGSSVersionOrder: true,
		conflictActionVersionOrder:   true,
	},
	// TSX generic calls must keep the last conflict action's source-version
	// priority when a later declaration reconverges their packed histories.
	// Issue #1443 certifies the conflict rule without packed reduction ordering.
	"tsx": {
		exactStackNodeEquivalence:  true,
		conflictActionVersionOrder: true,
	},
	// The Agda table's repetition shift at state 4039 on `id` must decline so
	// the parser can reduce `_atoms` and finish the function head. Tree-sitter's
	// C runtime stops on repetition shifts in conflicts, so this exact-blob
	// exclusion restores that behavior without changing other Agda rows.
	"agda": {
		conflictPolicyExclusions: []conflictPolicyExclusionProfile{{
			state:         4039,
			lookahead:     1,
			kind:          gotreesitter.ConflictPolicyRepetitionShift,
			reduceSymbols: []gotreesitter.Symbol{173},
		}},
	},
	// The canonical compact corpus certifies Go's converged-path split drops
	// against the production parser and the tree-sitter C oracle.
	// The owned EOF bundle requires its executed recovery route before publication.
	"go": {
		// Re-certified 2026-10-24 against grammars/grammar_blobs/go.bin
		// SHA-256 4456f9b11706face73bd511d6e80bbd2f0d468055065296db6f5334fa101d379,
		// regenerated with cmd/grammargen (no -lr-split; -lr-split
		// interacts badly with the Go external ASI scanner, see
		// grammargen/README.md "Go's blob is generated without -lr-split").
		// The compact/converged-split-drops and owned-EOF-recovery
		// certifications were re-validated against this blob via the
		// cgo_harness Go compact/recovery receipts (TestGoCompactIncrementalExecution*,
		// TestGoCompactRecoveryVersionTurnsLockedC, TestStage6GoCompactCertification).
		compactConvergedSplitDrops: true,
		compactOwnedEOFRecovery:    true,
		// Let the absorbing version compete before the eagerly shifted
		// missing-target version, matching C's assignment-error election.
		recoveryMissingVersionTurns: true,
		// One version per merge key in C's physical version order, and C's
		// merge across dynamic precedence: when two readings reach one state
		// with equal precedence, the earlier version's link wins. Receipts:
		// 1570 of 1596 real Go files now match locked C exactly (1301 before,
		// with one file parsed as an error that C reads clean); the remaining
		// 26 differ only in the pre-existing blank-identifier import leaf. The
		// four benchmark fixtures match locked C, and the query_compile fixture
		// keeps 3 live stacks instead of 12.
		fullParseGSSConvergence:      true,
		compactPackedGSSVersionOrder: true,
		legacyMergeAdmission:         gotreesitter.LegacyMergeAdmitDynamicPrecedence,
		// Locked C exposes these aliased DFA leaves even when error recovery
		// pops them before the enclosing string production can reduce.
		visibleTerminalAliases: []visibleTerminalAliasProfile{
			{internalName: "_raw_string_literal_token1", publicName: "raw_string_literal_content"},
			{internalName: "_interpreted_string_literal_token1", publicName: "interpreted_string_literal_content"},
		},
	},
	// Re-certified 2026-09-20 against grammars/grammar_blobs/erlang.bin
	// SHA-256 aa63243bd946d324bfb335d84c162cbaf8cc9922feef71de17296a18186ecedf
	// (WhatsApp/tree-sitter-erlang bump to 6ba4c762eb30). Both grants were
	// re-validated against this blob: compactConvergedSplitDrops via
	// TestAdmissionCandidateErlangConvergedSplitAdversarialProbe (host,
	// admission_switch_erlang_converged_split_probe_test.go), and
	// compactPackedGSSVersionOrder via
	// TestCompactPackedGSSVersionOrderIssue984MatchesC (cgo_harness,
	// re-pinned deep digests) plus the re-derived
	// TestCTopologyReceiptErlangIssue984 anchors.
	"erlang": {
		compactConvergedSplitDrops:   true,
		compactPackedGSSVersionOrder: true,
	},
	// F# has one declaration-name state where C omits a same-span unary
	// long_identifier below long_identifier_or_op. Expression identifiers and
	// dotted long identifiers retain the wrapper.
	"fsharp": {
		nativeUnaryWrapperFlattening: []nativeUnaryWrapperFlatteningProfile{
			{
				publicParent:        "long_identifier_or_op",
				wrapper:             "long_identifier",
				leaf:                "identifier",
				wrapperPreGotoState: 4258,
			},
		},
	},
	// The internal DFA skips each JSDoc line decoration before it produces the
	// next tag token. The compact materializer admits the resulting interior
	// gap only when the accepted terminal carries that exact prefix evidence.
	"jsdoc": {
		compactLexerSkippedPrefixTiling: true,
	},
	// The compact Markdown inline route reduces three proved emphasis rows.
	// It starts an HTML tag through the sole shift in the fourth exact row.
	// All other repetition rows remain outside compact admission.
	"markdown_inline": {
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{State: 18, Lookahead: 56, Kind: gotreesitter.ConflictPolicyRepetitionReduce, CompactOnly: true, CompactMinFrontierHeaders: 2, ReduceSymbols: []gotreesitter.Symbol{107}},
			{State: 18, Lookahead: 50, Kind: gotreesitter.ConflictPolicyRepetitionReduce, CompactOnly: true, CompactMinFrontierHeaders: 2, ReduceSymbols: []gotreesitter.Symbol{107}},
			{State: 18, Lookahead: 36, Kind: gotreesitter.ConflictPolicyRepetitionReduce, CompactOnly: true, CompactMinFrontierHeaders: 2, ReduceSymbols: []gotreesitter.Symbol{107}},
			{State: 209, Lookahead: 50, Kind: gotreesitter.ConflictPolicyShift, CompactOnly: true, ReduceSymbols: []gotreesitter.Symbol{104}},
		},
	},
	// C-oracle parity certifies the complete compact recovery competition for
	// this exact artifact against the ten pinned HTML witnesses.
	"html": {
		compactStrategy2ErrorRegion:  true,
		compactMissingTokenInsertion: true,
	},
	// Objective-C keeps parity-relevant alternatives below the bounded stack
	// comparison frontier. Exact comparison preserves them until generic result
	// selection chooses the C-equivalent alias-target shape.
	"objc": {
		exactStackNodeEquivalence: true,
	},
	// JavaScript's speculative automatic forest phase is bounded separately
	// from the production fallback. The exact locked corpus retained identical
	// trees and outcomes across all files while cutting aggregate allocations;
	// explicit forest parsing and non-certified languages keep the full budget.
	//
	// The same exact artifact certifies the pinned compact recovery frontier.
	// Mixed shift/reduce states decline instead of losing reduction versions.
	"javascript": {
		automaticForestMemoryAllowance:    javascriptAutomaticForestMemoryAllowance,
		compactConvergedSplitDrops:        true,
		compactStrategy2ErrorRegion:       true,
		compactMissingTokenInsertion:      true,
		compactRecoveryTrailingRetirement: true,
		compactRecoveryErrorModeKeyword:   true,
		compactRecoveryPlainFirst:         true,
	},
	// The locked Scala artifact certifies S5 missing insertion through a
	// physical graph-head merge and a per-version lexer split. It also certifies
	// S5 reductions before one missing insertion at EOF. The profile excludes
	// stack-summary recovery and the distinct recover_eof ERROR-root route.
	"scala": {
		externalScannerCheckpointReuse:      true,
		compactStrategy2ErrorRegion:         true,
		compactMissingTokenInsertion:        true,
		compactS5EOFMissingInsertion:        true,
		compactFaithfulS5Recovery:           true,
		compactPrimaryAcceptDerivation:      true,
		compactAcceptanceStructuralElection: true,
		compactRecoveryTrailingRetirement:   true,
		compactRecoveryErrorModeKeyword:     true,
	},
	// These scanner-backed grammars have certified the first retry ladder's
	// selected accepted-error tree as authoritative. Repeating the whole ladder
	// does not improve the selected tree and imposes a full additional parse.
	//
	// Dart also certifies its three profiled repeat-boundary folds (enum
	// bodies, extension bodies, and top-level declaration lists) against this
	// exact blob. The state and reduce-symbol checks preserve the retired
	// helper's scope without a linear scan over every lookahead at each state.
	//
	// Derive each state again after a grammar bump. For one repeat symbol name,
	// take the state that holds the most shift-over-reduce cells whose reduce
	// actions all carry that symbol. Grammar be07cf7118d3 gives state 601 for
	// enum_body_repeat2, state 616 for extension_body_repeat1, and state 574
	// for program_repeat4. Each maximum is unique. The blob also carries its
	// own table-derived precedence rows, and this profile only appends to them.
	"dart": {
		// C's version order keeps nullable function types grouped together.
		compactPackedGSSVersionOrder:  true,
		conflictActionVersionOrder:    true,
		externalScannerFullParseRetry: gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
		nativeResultCompatibility:     gotreesitter.ResultCompatibilityNativeCollapsedChildren,
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{State: 601, Lookahead: gotreesitter.ConflictPolicyAnyLookahead, Kind: gotreesitter.ConflictPolicyRepetitionShift, ReduceSymbols: []gotreesitter.Symbol{511}},
			{State: 616, Lookahead: gotreesitter.ConflictPolicyAnyLookahead, Kind: gotreesitter.ConflictPolicyRepetitionShift, ReduceSymbols: []gotreesitter.Symbol{516}},
			{State: 574, Lookahead: gotreesitter.ConflictPolicyAnyLookahead, Kind: gotreesitter.ConflictPolicyRepetitionShift, ReduceSymbols: []gotreesitter.Symbol{473}},
		},
	},
	// C#'s certified low-pressure accepted-error trees are authoritative. A
	// fresh no-stacks parse instead benefits from a bounded cap-16 retry; the
	// generic cap-48 ladder exceeds the large-file memory and time budgets.
	"c_sharp": {
		// Preserve collection expressions over the element-binding alias.
		compactPackedGSSVersionOrder:  true,
		externalScannerFullParseRetry: gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
		fullParseGSSConvergence:       true,
		nativeResultCompatibility: gotreesitter.ResultCompatibilityCSharpNativeNotNull |
			gotreesitter.ResultCompatibilityCSharpNativeUnicodeIdentifiers |
			gotreesitter.ResultCompatibilityCSharpNativeScopedLambdaStatements |
			gotreesitter.ResultCompatibilityCSharpNativeScopedLambdaBlocks |
			gotreesitter.ResultCompatibilityCSharpNativeQueryExpressions |
			gotreesitter.ResultCompatibilityNativeRecoveredStructure,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry:             true,
			SkipCompleteMaxEntryScratchPeak:            csharpAcceptedErrorRetryMaxEntryScratchPeak,
			FreshErrorNoStacksRetryMaxStacks:           csharpFreshErrorNoStacksRetryMaxStacks,
			SkipInitialCompleteAcceptedErrorMergeRetry: true,
			GSSConvergenceAcceptedErrorMergePerKey:     csharpGSSConvergenceErrorMergePerKey,
		},
	},
	// Crystal's external-scanner repeat selects the same tree after the complete
	// accepted-error retry ladder. Keep the full ladder, but do not run it twice
	// for this exact built-in grammar/scanner artifact.
	"crystal": {
		externalScannerFullParseRetry: gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
	},
	"cue": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	"gitcommit": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	// Haxe's accepted-error retry ladder selects the same tree on every pass.
	// Keep the first accepted result instead of running either retry ladder.
	"haxe": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	// Large V files can reuse the clean wide pass at the recovery-wide slot.
	// Keep the narrow merge pass because it changes some final V trees.
	"v": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			ReuseCleanWideForWideRetry:   true,
			ReuseCleanWideMinSourceBytes: vAcceptedErrorRetryMinSourceBytes,
		},
	},
	// Kotlin certifies CompactPrimaryAcceptanceDerivationCertified only.
	// selectCompactAcceptanceDerivation's materiality gate
	// (parsercore_phase0_driver.go, compactAcceptanceElectionIsVacuous) is
	// what makes that grant safe: the object_declaration tied election
	// (object Singleton { fun work() = Unit }, issue #93) declines instead
	// of publishing an unproven positional primary. Full-corpus field-aware
	// C-oracle verification finds zero unadjudicated divergence under that
	// gate for this mechanism alone (A3 certification workstream,
	// spec.campaign.v7; cgo_harness/kotlin_a3_certification_sweep_test.go).
	//
	// CompactConvergedReductionSplitDropsCertified stays withheld. Review
	// found a compact-only divergence class the certification sweep's
	// 3-file real corpus did not surface: an annotated extension property
	// with a getter, followed by a trailing comment (line or block), for
	// example
	//
	//	@Deprecated("old")
	//	val Int.double: Int
	//	    get() = this * 2
	//	// trailing
	//
	// Production is C-exact on this shape. With split-drops alone forced
	// on, the compact route accepts a C-divergent tree: the annotation is
	// torn into a bare prefix_expression and the getter becomes an
	// assignment (first divergence at /source_file, child count 4 vs the
	// C oracle's 3). Primary-acceptance-derivation alone is clean on the
	// same witness: it declines the converged-path split, falling back to
	// production's correct tree. See
	// kotlinA3AdversarialSources's annotated_extension_property_getter_*
	// entries for the pinned regression sweep coverage.
	//
	// A net fidelity ledger over 1,061 truncations of the Kotlin real
	// corpus found split-drops regresses more than it improves (2 improved,
	// 5 regressed) against a clean primary-accept-only baseline (0
	// improved, 0 regressed). The two witnesses split-drops alone would
	// otherwise fix -- the platform-modifier-recovery witness (internal
	// actual fun f(): String = "x") and small__TimeZoneNative.kt -- both
	// decline under primary-accept-only and fall back to production's own
	// (C-divergent, issue #93-adjacent) tree; that reversion to the
	// pre-certification status quo is the accepted, measured cost of
	// withholding this grant. The derivation-selection defect in the
	// converged-path split mechanism itself needs its own repair lane
	// before this grant is reconsidered.
	//
	// Recertified on blob 8c618126 (tree-sitter-kotlin 1852ea17, Kotlin 2.1
	// multi-dollar interpolation, 2026-09-19). The full-corpus sweep on the
	// shipped primary-accept-only profile saw 16 files: 11 accepted, 5
	// declined at the converged-path split, 0 divergences. Upstream #280
	// ("Prefer class and object declarations over infix expressions")
	// removed the infix derivation behind issue #93: the object_declaration
	// witness has no tied election left, and the platform-modifier witness
	// (internal actual fun f(): String = "x") is now C-exact in production.
	// Split-drops stays withheld; this recertification did not re-run the
	// split-drops ledger on the new blob.
	//
	// Moved to blob 6dd6ce27 (tree-sitter-kotlin d6af96d2). The full-corpus
	// sweep has not been re-run on this blob.
	"kotlin": {
		externalScannerFullParseRetry:  gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
		nativeResultCompatibility:      gotreesitter.ResultCompatibilityNativeCollapsedChildren,
		compactPrimaryAcceptDerivation: true,
		// C keeps every reading of an ambiguous Kotlin statement as a stack
		// version, merges a reduction into any version at the same state and
		// position, and elects among the merged children with
		// ts_subtree_compare. Without C's version order the parser keeps one
		// reading early and leaves an ERROR node in valid Kotlin: a generic
		// call statement (listOf<P>()), two index assignments in a row, an
		// annotation then a modifier then a typed parameter, a member index
		// inside a lambda (f { it.v[1] }), and the suspend receiver type in
		// testdata/real/kotlin. Locked-C receipts:
		// TestKotlinCompleteStatementsMatchLockedC, and a 994-file census of
		// JetBrains/Exposed plus a Kotlin compiler-plugin reader, where Go
		// error trees that C reads cleanly fall from 292 to 0 and exact
		// clean trees rise from 622 to 975.
		compactPackedGSSVersionOrder: true,
		// ts_stack_can_merge has no dynamic-precedence term: C merges a
		// reduction into an earlier version at the same state even when
		// the two readings' dynamic precedence differs, and the merge
		// decides which reading the parse continues with. Without it an
		// annotation on its own line before a generic call with a trailing
		// lambda (@A then a<E> {}) reads as a prefix_expression where C
		// reads a statement annotation.
		legacyMergeAdmission: gotreesitter.LegacyMergeAdmitDynamicPrecedence,
		// The automatic-semicolon scanner skips a block comment before an
		// empty semicolon, and C keeps that comment in the token's padding,
		// inside the parent the semicolon ends. Without it a block comment
		// after the package header (package a, /* x */, enum class S) is
		// left between two parents, recovery wraps it in an ERROR node and
		// the parse stops at the iteration limit.
		zeroWidthChildPaddingSpan: true,
	},
	// Apex's tied class-literal election matches the C oracle once the
	// compact route selects the sole primary derivation; the compact tree is
	// strictly more faithful than production plus the compat arm here. Apex
	// has no converged-path split-drop shape, so it does not certify that
	// mechanism. Full-corpus field-aware C-oracle verification certifies this
	// exact blob (A3 certification workstream, spec.campaign.v7).
	//
	// Recertified on blob 28fc59c4 (tree-sitter-sfapex da568eee,
	// 2026-09-20). The bump adds the multi_line_string_literal rule and
	// widens line_comment; the class-literal election shape does not change.
	// The A3 full-corpus sweep on the new blob reports 0 divergences.
	"apex": {
		nativeResultCompatibility:      gotreesitter.ResultCompatibilityNativeCollapsedChildren,
		compactPrimaryAcceptDerivation: true,
	},
	"elixir": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
		// Q0 (#1311): one legacy version per merge key, and C's merge across
		// dynamic precedence. Receipts: the generated 32 KiB cliff keeps 2
		// live stacks instead of 8, and the largest pinned corpus file now
		// matches locked C. The other corpus files keep their verdicts.
		fullParseGSSConvergence: true,
		legacyMergeAdmission:    gotreesitter.LegacyMergeAdmitDynamicPrecedence,
	},
	// Q0 (#1311): one legacy version per merge key and C's
	// representation-free merge. Receipts: run-tests.php keeps 5 live stacks
	// instead of 80, and the pinned corpus, the cliff file, and 144
	// edit-session inputs keep their locked-C verdicts.
	"php": {
		fullParseGSSConvergence: true,
		legacyMergeAdmission:    gotreesitter.LegacyMergeAdmitMixedRepresentation,
	},
	// Q0 (#1311): one legacy version per merge key. Receipts: the pinned
	// sample keeps 2 live stacks instead of 4, and all 144 edit-session inputs
	// still match locked C.
	"markdown": {
		fullParseGSSConvergence: true,
	},
	"hack": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	// Q0 (#1311): one legacy version per merge key. Receipts: the pinned
	// corpus peaks at 2 live stacks instead of 12 with 83% fewer new nodes, and
	// the corpus and 144 edit-session inputs keep their locked-C verdicts.
	"kconfig": {
		fullParseGSSConvergence: true,
	},
	"ruby": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	"rust": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	// Recertified on blob b954781f (tree-sitter-r 58a22794466c, 2026-09-20).
	// The bump splits _raw_string_literal into _raw_string_open /
	// _raw_string_content / _raw_string_close and adds an identifier-
	// continuation guard to the ELSE scan. Neither change touches a
	// collapsed-child pair: collapsedChildOccurrenceRules carries no "r"
	// entry, so this grant stays the pre-bump no-op it already was.
	"r": {
		nativeResultCompatibility: gotreesitter.ResultCompatibilityNativeCollapsedChildren,
	},
	// Matlab's external-scanner repeat selects the same tree after the complete
	// accepted-error retry ladder. Keep the full ladder, but do not run it twice
	// for this exact built-in grammar/scanner artifact.
	"matlab": {
		externalScannerFullParseRetry: gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
	},
	// Odin's accepted-error retry ladder selects the same tree on every pass,
	// and its locked large test-vector witness certifies the ASCII full-arena
	// density cap. Both policies remain bound to this exact blob identity.
	// Moved to blob 85ba9481 (tree-sitter-odin 2b234040). The retry-profile
	// corpus certifier has not been re-run on this blob.
	"odin": {
		fullParseArenaDensityCap: true,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	// Python's tuple-assignment and f-string interpolation ties match the C
	// oracle when the compact route applies C's raw subtree ordering. The
	// selection is clean-only and has no recovery authority. Full-corpus,
	// field-aware C-oracle verification certifies this exact blob for the
	// structural election and the existing converged-path split-drop and
	// primary-acceptance mechanisms (A3 certification workstream,
	// spec.campaign.v7). The mixed flat/GSS receiver path is also certified
	// for this exact artifact.
	"python": {
		externalScannerFullParseRetry:       gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
		compactConvergedSplitDrops:          true,
		compactPrimaryAcceptDerivation:      true,
		compactAcceptanceStructuralElection: true,
		compactMixedGSSMerge:                true,
	},
	// Perl's tied push-list election matches the C oracle once the compact
	// route accepts after a converged-path split drop and selects the sole
	// primary derivation. Full-corpus field-aware C-oracle verification
	// certifies both mechanisms for this exact blob (A3 certification
	// workstream, spec.campaign.v7).
	//
	// Bumped to blob 86d67a08 (tree-sitter-perl 8917c6e9, 2026-09-20). That
	// bump adds _RECOVER_PAREN_CLOSE (a synthetic close-paren the C scanner
	// emits during error recovery inside an unclosed call argument list) and
	// moves perl from the action-probe external-lex-state fallback to
	// nextGLRScoredExternalToken (the new blob carries 50 external lex state
	// rows; the old one carried 0). Symbol IDs 252-289 are unchanged; the new
	// external took 290 and _ERROR moved from 290 to 291. Both grants below
	// were re-proven against this exact blob with
	// TestPerlA3CompactCertificationFullCorpusSweep (cgo_harness, Docker)
	// before this pin landed; see that test for the full-corpus receipt.
	"perl": {
		compactConvergedSplitDrops:     true,
		compactPrimaryAcceptDerivation: true,
		// C keeps the shift derivation when the declared list conflict
		// converges, and the parenthesized call remains unambiguous.
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{State: 1545, Lookahead: 16, Kind: gotreesitter.ConflictPolicyShift, ReduceSymbols: []gotreesitter.Symbol{333}},
			{State: gotreesitter.ConflictPolicyAnyState, Lookahead: 23, Kind: gotreesitter.ConflictPolicyShift, ReduceSymbols: []gotreesitter.Symbol{397}},
		},
	},
	// Ada's tied aggregate elections (positional-array and others-choice)
	// match the C oracle once the compact route accepts after a
	// converged-path split drop and selects the sole primary derivation.
	// Full-corpus field-aware C-oracle verification certifies both
	// mechanisms for this exact blob (A3 certification workstream,
	// spec.campaign.v7).
	//
	// Ada's grammar declares conflicts: [[$.component_choice_list,
	// $.discrete_choice], ...] (tree-sitter-ada grammar.js): a bare
	// "others" choice reduces the shared token as either component_choice_list
	// (162... record-aggregate reading, RM 4.3.1) or discrete_choice
	// (array-aggregate reading, RM 3.8.1), both plain REDUCE actions with no
	// dynamic precedence on either alternative. The locked C oracle keeps
	// the discrete_choice reading; this row-scoped fold (state/lookahead
	// wildcarded, matched only by the declared-conflict symbol pair,
	// component_choice_list=196, discrete_choice=256 in this blob) reproduces
	// that natively, the same mechanism ql's signatureExpr election uses.
	//
	// Recertified on blob fbb1e5cf (tree-sitter-ada dd5fa4cd, 2026-09-20).
	// The bump adds the finally_part rule and rewrites
	// handled_sequence_of_statements; it shifts every symbol above
	// finally_part by one, so the declared-conflict pair moved from 195/255
	// to 196/256. Both identifiers come from the shipped blob by name. The
	// A3 full-corpus sweep on the new blob reports files=23 accepted=20
	// declined=3 divergences=0, and the three ada C-oracle receipts pass.
	"ada": {
		compactConvergedSplitDrops:     true,
		compactPrimaryAcceptDerivation: true,
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{
				State:         gotreesitter.ConflictPolicyAnyState,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyDeclaredReduceReduceHighestSymbol,
				ReduceSymbols: []gotreesitter.Symbol{196, 256},
			},
		},
	},
	// Swift's external-scanner full-parse retry does not need to repeat: a
	// width disagreement the first pass already resolved does not change on
	// a second pass. Repeating that ladder for the external scanner does not
	// help.
	//
	// The accepted-error skip now applies only to sources of 64 KiB and
	// more. On blob 33fd9742, three small corpus files (swift-algorithms_
	// Chunked.swift 27,814 bytes, FlattenCollection.swift 9,196, and
	// Stride.swift 7,887) parse clean only through the complete retry
	// ladder, so skipping that ladder for them regresses clean to error.
	// The 104,681-byte issue-586 witness (stdlib_FloatingPointToString.swift)
	// needs the opposite: without the skip it stays unbounded under the race
	// detector (825s, over the CI shard timeout, versus 273s with the skip).
	// The size gate keeps both: small sources always take the full ladder,
	// and only large sources skip it. The B16 telemetry receipt, 25/25
	// real-corpus parity and race-timed witness were measured on blob 33fd9742
	// (tree-sitter-swift 00bbb0a2, 2026-09-20), not on the blob pinned below.
	"swift": {
		// Preserve navigation expressions with C's physical version order.
		compactPackedGSSVersionOrder: true,
		// The locked C ERROR lex row enables all 34 external symbols.
		errorModeAllExternalSymbols:   true,
		externalScannerFullParseRetry: gotreesitter.ExternalScannerFullParseRetrySkipRepeat,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry:  true,
			SkipCompleteMaxEntryScratchPeak: 3 * 64 * 1024,
			SkipCompleteMinSourceBytes:      64 * 1024,
		},
	},
	// Large ASM accepted-error parses improve during the same-stack merge
	// retry, but later widened-stack passes do not advance the selected tree.
	// Keep that first retry while avoiding the redundant wider passes.
	"asm": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			MinSourceBytes:      11 * 1024,
			InitialStackCeiling: 8,
		},
	},
	// These grammars have error-bearing real-corpus witnesses that legitimately
	// reach EOF. Re-running the accepted-error ladder does not improve their
	// selected trees, so the exact certified blobs keep the first result.
	"bash": {
		compactConvergedSplitDrops: true,
		// Retain the pipeline predecessor when its redirected tail converges.
		fullParseGSSConvergence: true,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	"caddy": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	"cpp": {
		recoveryStackVersionOrder: true,
		// Preserve expression predecessors when qualified names converge with
		// declaration readings, so nested calls remain valid call arguments.
		fullParseGSSConvergence: true,
		// Resolve sizeof's expression/type ambiguity in C's physical order.
		compactPackedGSSVersionOrder: true,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	// On the locked large accepted-error witnesses, the non-recovery and
	// recovery-enabled widened-stack passes produce the same complete tree and
	// semantic runtime state. Retain the first result instead of parsing it
	// again; exact blob identity keeps adapted and future grammars conservative.
	"enforce": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			ReuseCleanWideForWideRetry:   true,
			ReuseCleanWideMinSourceBytes: 128 * 1024,
		},
	},
	"kdl": {
		automaticForestEnabled: true,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	// These exact built-in artifacts have authenticated corpus receipts with no
	// forest/C-oracle divergence on accepted forest files, no
	// routed/production divergence on any file, and an aggregate route wall-time
	// win after conservative production fallbacks.
	"awk": {
		automaticForestEnabled: true,
	},
	"uxntal": {
		automaticForestEnabled: true,
	},
	// Meson's tied smoke election compares variableunit with var_unit. The
	// locked C runtime selects variableunit through its raw subtree comparator.
	// The retry ladder still changes selected trees on small error-bearing files.
	"meson": {
		compactPrimaryAcceptDerivation:      true,
		compactAcceptanceStructuralElection: true,
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
			SkipCompleteMinSourceBytes:     mesonAcceptedErrorRetryMinSourceBytes,
		},
	},
	"rego": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	"scss": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
	},
	// Tcl's complete accepted-error trees and first error-bearing no-stacks
	// retry are authoritative for the certified corpus. Later widened passes
	// repeat the same failure without advancing the selected tree.
	"tcl": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
			FreshErrorNoStacksMaxPasses:    1,
		},
	},
	// Haskell has three exact repeat rows where C selects the reduce arm.
	// Retaining both arms grows new graph-structured stack frontiers.
	"haskell": {
		compactConvergedSplitDrops: true,
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{
				State:         9609,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyRepetitionReduce,
				ReduceSymbols: []gotreesitter.Symbol{518},
			},
			{
				State:         10984,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyRepetitionReduce,
				ReduceSymbols: []gotreesitter.Symbol{516},
			},
			{
				State:         11192,
				Lookahead:     4,
				Kind:          gotreesitter.ConflictPolicyRepetitionReduce,
				ReduceSymbols: []gotreesitter.Symbol{500},
			},
		},
	},
	// On large accepted-error Java sources, the cap-16 same-stack merge retry
	// remains authoritative; the subsequent cap-64 clean/recovery passes do not
	// improve the selected tree. Keep this bound pinned to the exact built-in
	// blob so overrides and caller-built languages retain the generic ladder.
	"java": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			MinSourceBytes:      64 * 1024,
			InitialStackCeiling: 14,
		},
	},
	// Large Groovy and D accepted-error parses stay inside their certified
	// initial stack ceilings. Widening does not improve the selected tree and
	// reintroduces their real-corpus time/RSS cliffs. The generic profile gate
	// keeps incremental fallbacks and explicit diagnostic overrides
	// conservative.
	"groovy": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			MinSourceBytes:      64 * 1024,
			InitialStackCeiling: 2,
		},
	},
	"d": {
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			MinSourceBytes:      64 * 1024,
			InitialStackCeiling: 3,
		},
	},
	// NOTE on dot: no certified row here (unlike gomod/dart/c below) — dot's
	// retired dotRepetitionShiftConflictChoice is NOT migrated. dot isn't in
	// cRepetitionSkipOptOut, so the engine-wide C repetition-skip fold
	// already folded its stmt_list_repeat1 boundary (state 4) with a flat
	// parse stack; the retired helper's only call site was forest-only and
	// dot isn't in builtinForestDefaults, so it never actually ran for this
	// blob. A/B validation found identical accepted trees, but
	// reviving the helper's shift preference as a policy grows the LR stack
	// O(n) with statement count (406 deep on a 400-statement file vs 8 for
	// the fold already in place) for no fork-count benefit (both already
	// MaxStacksSeen=1). Retired outright instead of migrated.
	//
	// go.mod's require-list continuation forks at two states: the top-level
	// source_file list (state 3) and a parenthesized require block (state
	// 37). Replaces gomodRepetitionShiftConflictChoice.
	// Unlike dot, this is not a revival: the retired helper's dispatch call
	// site (conflict_policy.go's gomod special case, still present) ran
	// unconditionally ahead of the reuse gate and the engine-wide fold, so
	// this exact preference was already the live, shipped behavior —
	// confirmed byte-identical MaxStacksSeen/PeakStackDepth before and after
	// this migration on a synthetic 400-require go.mod.
	"gomod": {
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{State: 3, Lookahead: gotreesitter.ConflictPolicyAnyLookahead, Kind: gotreesitter.ConflictPolicyRepetitionShift, ReduceSymbols: []gotreesitter.Symbol{50}},
			{State: 37, Lookahead: gotreesitter.ConflictPolicyAnyLookahead, Kind: gotreesitter.ConflictPolicyRepetitionShift, ReduceSymbols: []gotreesitter.Symbol{52}},
		},
	},
	// C's top-level declaration list (translation_unit_repeat1) and
	// preprocessor-conditional body (preproc_if_repeat1) close only on
	// terminators with no continuation shift (EOF; #endif/#elif/#else), so
	// any lookahead offering a continuation shift makes the competing reduce
	// a zero-progress dead end. That C-faithful rule is scoped by reduce
	// symbol identity alone, not by table position: it recurs at hundreds of
	// (state, lookahead) rows across the grammar (a scan of this exact blob
	// found 3242, across 446 states), which is why this is the one certified
	// profile using the ConflictPolicyAnyState/AnyLookahead wildcards instead
	// of an enumerated per-row table. case_statement_repeat1 is deliberately
	// excluded: its list boundary is load-bearing, not a dead-end. Replaces
	// cRepetitionShiftConflictChoice; byte-for-byte C
	// parity held by TestParityCTopLevelDeclAmbiguity /
	// TestParityCPreprocConditional.
	//
	// Re-pinned 2026-09-20 to grammars/grammar_blobs/c.bin SHA-256
	// db0123b46b06dbfdb139e834dff30ceca3b8df123aac39865813a596eb819ef9
	// after the tree-sitter-c bump to b780e47fc780. The reduce symbols keep
	// their identity in the new blob: 324 translation_unit_repeat1, 326
	// preproc_if_repeat1, 343 enumerator_list_repeat1. The cgo receipts
	// TestParityCTopLevelDeclAmbiguity, TestParityCPreprocConditional,
	// TestIssue667CEnumListsMatchCReference, TestCExternCWrapperParity, and
	// TestStage6CCompactCertification pass against this blob.
	"c": {
		recoveryStackVersionOrder: true,
		// Measurements on large C witnesses show identical trees before and after
		// the retry ladder. Keep the initial tree and avoid repeated full parses
		// for this exact grammar blob.
		fullParseAcceptedErrorRetryProfile: gotreesitter.FullParseAcceptedErrorRetryProfile{
			SkipCompleteAcceptedErrorRetry: true,
		},
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{
				State:         gotreesitter.ConflictPolicyAnyState,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyRepetitionShift,
				ReduceSymbols: []gotreesitter.Symbol{324, 326},
			},
			// Fold aux_sym_enumerator_list_repeat1 at its conflict cell instead
			// of forking. Upstream never executes SHIFT_REPEAT at a conflict cell
			// (lib/src/parser.c, `if (action.shift.repetition) break;`): it
			// reduces, renumbers, and re-dispatches the same lookahead. C is
			// opted out of the engine-wide equivalent (cRepetitionSkipOptOut in
			// conflict_policy.go), so an enumerator list forks and then depends
			// on cap-one convergence. When one branch is in demoted-linear form
			// the GSS merge refuses, the depth tiebreak keeps the unfolded
			// branch, and a three-enumerator list dead-ends on `}` and invents a
			// comma. Folding this one symbol restores the upstream order.
			{
				State:         gotreesitter.ConflictPolicyAnyState,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyRepetitionReduce,
				ReduceSymbols: []gotreesitter.Symbol{343},
			},
		},
	},
	// Gleam's declared integer self conflict has plain reduce rows 732/735:
	// integer(1 child), then integer(2 children). C runs the sign-plus-digits
	// reduction immediately, preserving a negative integer instead of wrapping
	// a positive integer in integer_negation.
	"gleam": {
		conflictPolicies: []gotreesitter.ConflictPolicy{{
			State:         gotreesitter.ConflictPolicyAnyState,
			Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
			Kind:          gotreesitter.ConflictPolicyDeclaredSelfReduceReduceLongest,
			ReduceSymbols: []gotreesitter.Symbol{231},
		}},
	},
	// QL declares conflicts: [[$.simpleId, $.className], ...]. An upper-case
	// identifier in signature position (`implements Foo`, `implements
	// M::Foo`) reduces the shared _upper_id token as either simpleId (162,
	// building moduleExpr) or className (163, building typeExpr); both are
	// live GLR alternatives of signatureExpr's choice(typeExpr, moduleExpr,
	// predicateExpr) with no dynamic precedence on either production. The
	// locked C oracle (tree-sitter-ql 1fd627a4e8bff8c24c11987474bd33112bead857)
	// always keeps the className/typeExpr reading; this row-scoped fold
	// (state/lookahead wildcarded, matched only by the declared-conflict
	// symbol pair) reproduces that natively instead of relying on a
	// post-parse tree rewrite.
	"ql": {
		conflictPolicies: []gotreesitter.ConflictPolicy{
			{
				State:         gotreesitter.ConflictPolicyAnyState,
				Lookahead:     gotreesitter.ConflictPolicyAnyLookahead,
				Kind:          gotreesitter.ConflictPolicyDeclaredReduceReduceHighestSymbol,
				ReduceSymbols: []gotreesitter.Symbol{162, 163},
			},
		},
	},
	// The missing-token probe and absorbing version retain C physical turns.
	"sql": {
		recoveryMissingVersionTurns: true,
	},
	// PowerShell's backtick immediately followed by a newline is the
	// language's line-continuation escape: the C reference scanner consumes
	// it as ordinary skipped trivia (zero ERROR nodes across the sequence,
	// two children on the enclosing command_argument_sep — the same shape as
	// any other run of whitespace), never as a token of its own. Declaring it
	// here lets bytesAreParserPadding give it the same unconditional
	// treatment already applied to backslash+newline instead of falling
	// through to skipped-gap ERROR materialization. C-oracle verified on the
	// enclosing command_argument_sep shape (TestPowerShellBacktickContinuationIsParserPadding).
	"powershell": {
		lineContinuationEscapeByte: '`',
	},
}

type exactRuntimeProfileExternalScanner interface {
	externalScannerForExactRuntimeProfile() gotreesitter.ExternalScanner
}

func attachBuiltinLanguageRuntimeProfile(name string, lang *gotreesitter.Language) bool {
	if lang == nil {
		return false
	}
	profile, ok := builtinLanguageRuntimeProfiles[canonicalLanguageName(name)]
	if !ok {
		return false
	}
	changed := false
	for _, alias := range profile.visibleTerminalAliases {
		for symbol := 1; symbol < len(lang.SymbolNames) && uint32(symbol) < lang.TokenCount && symbol < len(lang.SymbolMetadata); symbol++ {
			if lang.SymbolNames[symbol] != alias.internalName {
				continue
			}
			lang.SymbolNames[symbol] = alias.publicName
			lang.SymbolMetadata[symbol].Name = alias.publicName
			lang.SymbolMetadata[symbol].Visible = true
			lang.SymbolMetadata[symbol].Named = true
			changed = true
		}
	}
	if profile.recoveryMissingVersionTurns && !lang.RecoveryMissingVersionTurnsCertified {
		lang.RecoveryMissingVersionTurnsCertified = true
		changed = true
	}
	if profile.errorModeAllExternalSymbols && len(lang.LexModes) > 0 && len(lang.ExternalSymbols) > 0 && len(lang.ExternalLexStates) < 1<<16 {
		rowIndex := int(lang.LexModes[0].ExternalLexState)
		var current []bool
		if rowIndex < len(lang.ExternalLexStates) {
			current = lang.ExternalLexStates[rowIndex]
		}
		if len(current) != len(lang.ExternalSymbols) || slices.Contains(current, false) {
			row := make([]bool, len(lang.ExternalSymbols))
			for i := range row {
				row[i] = true
			}
			// Ordinary states can share the old row. Give ERROR its own row
			// instead of enabling recovery-only symbols in those states.
			lang.LexModes = slices.Clone(lang.LexModes)
			lang.LexModes[0].ExternalLexState = uint16(len(lang.ExternalLexStates))
			lang.ExternalLexStates = append(slices.Clone(lang.ExternalLexStates), row)
			changed = true
		}
	}
	if profile.externalScannerCheckpointReuse {
		if certifier, ok := lang.ExternalScanner.(exactRuntimeProfileExternalScanner); ok {
			lang.ExternalScanner = certifier.externalScannerForExactRuntimeProfile()
			changed = true
		}
	}
	if profile.externalScannerFullParseRetry != gotreesitter.ExternalScannerFullParseRetryDefault &&
		lang.ExternalScanner != nil &&
		lang.ExternalScannerFullParseRetryPolicy != profile.externalScannerFullParseRetry {
		lang.ExternalScannerFullParseRetryPolicy = profile.externalScannerFullParseRetry
		changed = true
	}
	if profile.fullParseAcceptedErrorRetryProfile != (gotreesitter.FullParseAcceptedErrorRetryProfile{}) &&
		lang.FullParseAcceptedErrorRetryProfile != profile.fullParseAcceptedErrorRetryProfile {
		lang.FullParseAcceptedErrorRetryProfile = profile.fullParseAcceptedErrorRetryProfile
		changed = true
	}
	if profile.automaticForestMemoryAllowance > 0 &&
		lang.AutomaticForestMemoryAllowanceBytes != profile.automaticForestMemoryAllowance {
		lang.AutomaticForestMemoryAllowanceBytes = profile.automaticForestMemoryAllowance
		changed = true
	}
	if profile.automaticForestEnabled && !lang.AutomaticForestEnabledByDefault {
		lang.AutomaticForestEnabledByDefault = true
		changed = true
	}
	if profile.fullParseArenaDensityCap && !lang.FullParseArenaDensityCapEnabled {
		lang.FullParseArenaDensityCapEnabled = true
		changed = true
	}
	if profile.recoveryStackVersionOrder && !lang.RecoveryStackVersionOrderEnabled {
		lang.RecoveryStackVersionOrderEnabled = true
		changed = true
	}
	if profile.fullParseGSSConvergence && !lang.FullParseGSSConvergenceEnabled {
		lang.FullParseGSSConvergenceEnabled = true
		changed = true
	}
	if missing := profile.nativeResultCompatibility &^ lang.NativeResultCompatibility; missing != 0 {
		lang.NativeResultCompatibility |= missing
		changed = true
	}
	if rules := resolveNativeUnaryWrapperFlatteningProfile(lang, profile.nativeUnaryWrapperFlattening); len(rules) > 0 &&
		!slices.Equal(lang.NativeUnaryWrapperFlattening, rules) {
		lang.NativeUnaryWrapperFlattening = rules
		changed = true
	}
	if profile.compactConvergedSplitDrops && !lang.CompactConvergedReductionSplitDropsCertified {
		lang.CompactConvergedReductionSplitDropsCertified = true
		changed = true
	}
	if profile.compactPrimaryAcceptDerivation && !lang.CompactPrimaryAcceptanceDerivationCertified {
		lang.CompactPrimaryAcceptanceDerivationCertified = true
		changed = true
	}
	if profile.compactAcceptanceStructuralElection && !lang.CompactAcceptanceStructuralElectionCertified {
		lang.CompactAcceptanceStructuralElectionCertified = true
		changed = true
	}
	if profile.compactMixedGSSMerge && !lang.CompactMixedGSSMergeCertified {
		lang.CompactMixedGSSMergeCertified = true
		changed = true
	}
	if missing := profile.legacyMergeAdmission &^ lang.LegacyMergeAdmission; missing != 0 {
		lang.LegacyMergeAdmission |= missing
		changed = true
	}
	if profile.compactLexerSkippedPrefixTiling && !lang.CompactLexerSkippedPrefixTilingCertified {
		lang.CompactLexerSkippedPrefixTilingCertified = true
		changed = true
	}
	if profile.exactStackNodeEquivalence && !lang.ExactStackNodeEquivalenceCertified {
		lang.ExactStackNodeEquivalenceCertified = true
		changed = true
	}
	if profile.conflictActionVersionOrder && !lang.ConflictActionVersionOrderCertified {
		lang.ConflictActionVersionOrderCertified = true
		changed = true
	}
	if profile.compactPackedGSSVersionOrder && !lang.CompactPackedGSSVersionOrderCertified {
		lang.CompactPackedGSSVersionOrderCertified = true
		changed = true
	}
	if profile.zeroWidthChildPaddingSpan && !lang.ZeroWidthChildPaddingSpanCertified {
		lang.ZeroWidthChildPaddingSpanCertified = true
		changed = true
	}
	if profile.compactStrategy2ErrorRegion && !lang.CompactStrategy2ErrorRegionCertified {
		lang.CompactStrategy2ErrorRegionCertified = true
		changed = true
	}
	if profile.compactMissingTokenInsertion && !lang.CompactMissingTokenInsertionCertified {
		lang.CompactMissingTokenInsertionCertified = true
		changed = true
	}
	if profile.compactS5EOFMissingInsertion && !lang.CompactS5EOFMissingInsertionCertified {
		lang.CompactS5EOFMissingInsertionCertified = true
		changed = true
	}
	if profile.compactFaithfulS5Recovery && !lang.CompactFaithfulS5RecoveryCertified {
		lang.CompactFaithfulS5RecoveryCertified = true
		changed = true
	}
	if profile.compactOwnedEOFRecovery && !lang.CompactOwnedEOFRecoveryCertified {
		lang.CompactOwnedEOFRecoveryCertified = true
		changed = true
	}
	if profile.compactRecoveryTrailingRetirement && !lang.CompactRecoveryTrailingLineageRetirementCertified {
		lang.CompactRecoveryTrailingLineageRetirementCertified = true
		changed = true
	}
	if profile.compactRecoveryErrorModeKeyword && !lang.CompactRecoveryErrorModeKeywordCaptureCertified {
		lang.CompactRecoveryErrorModeKeywordCaptureCertified = true
		changed = true
	}
	if profile.compactRecoveryPlainFirst && !lang.CompactRecoveryPlainFirstCertified {
		lang.CompactRecoveryPlainFirstCertified = true
		changed = true
	}
	if profile.lineContinuationEscapeByte != 0 && lang.LineContinuationEscapeByte != profile.lineContinuationEscapeByte {
		lang.LineContinuationEscapeByte = profile.lineContinuationEscapeByte
		changed = true
	}
	for _, policy := range profile.conflictPolicies {
		if languageHasConflictPolicy(lang, policy) {
			continue
		}
		policy.ReduceSymbols = append([]gotreesitter.Symbol(nil), policy.ReduceSymbols...)
		lang.ConflictPolicies = append(lang.ConflictPolicies, policy)
		changed = true
	}
	if excludeBuiltinConflictPolicies(lang, profile.conflictPolicyExclusions) {
		changed = true
	}
	return changed
}

func excludeBuiltinConflictPolicies(lang *gotreesitter.Language, exclusions []conflictPolicyExclusionProfile) bool {
	if lang == nil || len(exclusions) == 0 || len(lang.ConflictPolicies) == 0 {
		return false
	}
	kept := lang.ConflictPolicies[:0]
	changed := false
	for _, policy := range lang.ConflictPolicies {
		excluded := false
		for _, exclusion := range exclusions {
			if policy.State == exclusion.state && policy.Lookahead == exclusion.lookahead &&
				policy.Kind == exclusion.kind && slices.Equal(policy.ReduceSymbols, exclusion.reduceSymbols) {
				excluded = true
				changed = true
				break
			}
		}
		if !excluded {
			kept = append(kept, policy)
		}
	}
	if changed {
		lang.ConflictPolicies = kept
	}
	return changed
}

func resolveNativeUnaryWrapperFlatteningProfile(
	lang *gotreesitter.Language,
	profile []nativeUnaryWrapperFlatteningProfile,
) []gotreesitter.UnaryWrapperFlatteningRule {
	if lang == nil || len(profile) == 0 {
		return nil
	}
	rules := make([]gotreesitter.UnaryWrapperFlatteningRule, 0, len(profile))
	for _, candidate := range profile {
		parent, parentOK := runtimeProfileNamedSymbol(lang, candidate.publicParent)
		wrapper, wrapperOK := runtimeProfileNamedSymbol(lang, candidate.wrapper)
		leaf, leafOK := runtimeProfileNamedSymbol(lang, candidate.leaf)
		if !parentOK || !wrapperOK || !leafOK {
			continue
		}
		rules = append(rules, gotreesitter.UnaryWrapperFlatteningRule{
			PublicParent:        parent,
			Wrapper:             wrapper,
			Leaf:                leaf,
			WrapperPreGotoState: candidate.wrapperPreGotoState,
		})
	}
	return rules
}

func runtimeProfileNamedSymbol(lang *gotreesitter.Language, name string) (gotreesitter.Symbol, bool) {
	if lang == nil {
		return 0, false
	}
	for index, symbolName := range lang.SymbolNames {
		if symbolName == name &&
			index < len(lang.SymbolMetadata) &&
			lang.SymbolMetadata[index].Visible &&
			lang.SymbolMetadata[index].Named {
			return gotreesitter.Symbol(index), true
		}
	}
	return 0, false
}

func runtimeProfileSymbol(lang *gotreesitter.Language, name string) (gotreesitter.Symbol, bool) {
	if lang == nil {
		return 0, false
	}
	var match gotreesitter.Symbol
	found := false
	for index, symbolName := range lang.SymbolNames {
		if symbolName != name {
			continue
		}
		if found {
			return 0, false
		}
		match = gotreesitter.Symbol(index)
		found = true
	}
	return match, found
}

func languageHasConflictPolicy(lang *gotreesitter.Language, want gotreesitter.ConflictPolicy) bool {
	if lang == nil {
		return false
	}
	for _, policy := range lang.ConflictPolicies {
		if policy.State == want.State &&
			policy.Lookahead == want.Lookahead &&
			policy.Kind == want.Kind &&
			slices.Equal(policy.ReduceSymbols, want.ReduceSymbols) {
			return true
		}
	}
	return false
}
