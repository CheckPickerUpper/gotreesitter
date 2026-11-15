package gotreesitter

import "testing"

// TestLexerAcceptsEndAtAnyPosition checks that a state accepting the end
// token accepts it wherever the scan reaches it, as tree-sitter's
// ACCEPT_TOKEN(ts_builtin_sym_end) does, and that the outcome tells the end
// token apart from a skipped run.
func TestLexerAcceptsEndAtAnyPosition(t *testing.T) {
	states := []LexState{{AcceptEOF: true, Default: -1, EOF: -1}}
	lex := NewLexer(states, []byte("x"))
	for _, position := range []int{0, 1} {
		if token, outcome := lex.scanContiguous(0, position, 0, uint32(position)); outcome != lexScanEnd || token.Symbol != 0 || token.StartByte != uint32(position) || token.EndByte != uint32(position) {
			t.Fatalf("position %d: token=%+v outcome=%d, want the end token there", position, token, outcome)
		}
	}
	lex.setIncludedRanges([]Range{{StartByte: 0, EndByte: 1, EndPoint: Point{Column: 1}}})
	for _, position := range []int{0, 1} {
		if token, outcome := lex.scanIncluded(0, position, 0, uint32(position)); outcome != lexScanEnd || token.Symbol != 0 || token.StartByte != uint32(position) || token.EndByte != uint32(position) {
			t.Fatalf("included position %d: token=%+v outcome=%d, want the end token there", position, token, outcome)
		}
	}
	skip := []LexState{{Default: -1, EOF: -1, Transitions: []LexTransition{{Lo: ' ', Hi: ' ', NextState: 1}}}, {Skip: true, Default: -1, EOF: -1}}
	if _, outcome := NewLexer(skip, []byte(" x")).scanContiguous(0, 0, 0, 0); outcome != lexScanSkipped {
		t.Fatalf("separator outcome=%d, want a skip", outcome)
	}
}

// TestLexerNextReturnsEndBeforeRemainingInput checks that Next stops at an
// end token the lex state accepts before the input ends, instead of reading
// the remaining bytes as an error.
func TestLexerNextReturnsEndBeforeRemainingInput(t *testing.T) {
	states := []LexState{{AcceptEOF: true, Default: -1, EOF: -1}}
	lex := NewLexer(states, []byte("\n"))
	token := lex.NextWithErrorRuns(0)
	if token.Symbol != 0 || token.StartByte != 0 || token.EndByte != 0 || lex.pos != 0 {
		t.Fatalf("token=%+v at %d, want the end token at byte 0", token, lex.pos)
	}
}

func TestLexerEOFAcceptReplacesEqualWidthToken(t *testing.T) {
	for _, tc := range []struct {
		name string
		end  LexState
		want Symbol
	}{
		{"end", LexState{AcceptEOF: true, Default: -1, EOF: -1}, 0},
		{"token", LexState{AcceptToken: 2, Default: -1, EOF: -1}, 2},
		{"nonaccepting", LexState{Default: -1, EOF: -1}, 1},
		{"lower_priority", LexState{AcceptEOF: true, AcceptPriority: 1, Default: -1, EOF: -1}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			states := []LexState{
				{Default: -1, EOF: -1, Transitions: []LexTransition{{Lo: 'a', Hi: 'a', NextState: 1}}},
				{AcceptToken: 1, Default: -1, EOF: 2},
				tc.end,
			}
			for _, included := range []bool{false, true} {
				source := []byte("a")
				if included {
					source = []byte("a!")
				}
				lex := NewLexer(states, source)
				var token Token
				var outcome lexScanOutcome
				if included {
					lex.setIncludedRanges([]Range{{StartByte: 0, EndByte: 1, EndPoint: Point{Column: 1}}})
					token, outcome = lex.scanIncluded(0, 0, 0, 0)
				} else {
					token, outcome = lex.scanContiguous(0, 0, 0, 0)
				}
				if outcome == lexScanFailed || token.Symbol != tc.want || token.StartByte != 0 || token.EndByte != 1 {
					t.Fatalf("included=%v: token=%+v outcome=%d, want symbol %d", included, token, outcome, tc.want)
				}
			}
		})
	}
}

func TestKeywordEOFAcceptReplacesEqualWidthToken(t *testing.T) {
	for _, end := range []LexState{
		{AcceptEOF: true, Default: -1, EOF: -1},
		{AcceptToken: 2, Default: -1, EOF: -1},
	} {
		lang := &Language{KeywordLexStates: []LexState{
			{Default: -1, EOF: -1, Transitions: []LexTransition{{Lo: 'a', Hi: 'a', NextState: 1}}},
			{AcceptToken: 1, Default: -1, EOF: 2},
			end,
		}}
		source := &dfaTokenSource{language: lang}
		token, ok := source.lexKeywordSource([]byte("a"))
		if ok != !end.AcceptEOF || token.Symbol != end.AcceptToken {
			t.Fatalf("end=%+v: token=%+v ok=%v", end, token, ok)
		}
	}
}
