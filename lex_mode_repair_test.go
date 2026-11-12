package gotreesitter

import "testing"

// TestRepairNoLookaheadLexModesKeepsStatesWithRealTokenActions checks that a
// state is marked no-lookahead only when every token other than EOF, internal
// or external, is trivia. A no-lookahead state reduces by its EOF actions
// without lexing, so a state where an external token shifts (or takes any
// other real action) must keep its lex state: reducing there by EOF would
// discard the shift the external token needs.
func TestRepairNoLookaheadLexModesKeepsStatesWithRealTokenActions(t *testing.T) {
	const (
		eof             = 0
		externalToken   = 2
		internalExtra   = 3
		internalToken   = 4
		reduceEntry     = 1
		shiftEntry      = 2
		shiftExtraEntry = 3
	)
	lang := &Language{
		TokenCount:      5,
		StateCount:      4,
		ExternalSymbols: []Symbol{externalToken},
		LexModes: []LexMode{
			{},
			{LexState: 7},
			{LexState: 8},
			{LexState: 9},
		},
		ParseActions: []ParseActionEntry{
			{},
			{Actions: []ParseAction{{Type: ParseActionReduce}}},
			{Actions: []ParseAction{{Type: ParseActionShift}}},
			{Actions: []ParseAction{{Type: ParseActionShift, Extra: true}}},
		},
		ParseTable: [][]uint16{
			make([]uint16, 5),
			// EOF reduces; the external and internal extras are trivia.
			{eof: reduceEntry, externalToken: shiftExtraEntry, internalExtra: shiftExtraEntry},
			// EOF reduces; the external token shifts.
			{eof: reduceEntry, externalToken: shiftEntry, internalExtra: shiftExtraEntry},
			// EOF reduces; an internal token shifts.
			{eof: reduceEntry, internalExtra: shiftExtraEntry, internalToken: shiftEntry},
		},
	}

	RepairNoLookaheadLexModes(lang)

	for _, want := range []struct {
		state     int
		lexState  uint32
		situation string
	}{
		{state: 1, lexState: noLookaheadLexState, situation: "only trivia beside the EOF reduce"},
		{state: 2, lexState: 8, situation: "an external token shifts"},
		{state: 3, lexState: 9, situation: "an internal token shifts"},
	} {
		if got := lang.LexModes[want.state].LexStateIndex(); got != want.lexState {
			t.Errorf("state %d (%s): LexStateIndex = %d, want %d", want.state, want.situation, got, want.lexState)
		}
	}
}
