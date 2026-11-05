package gotreesitter

// packedVersionAdvance ports the stop rule of ts_parser_parse's version loop
// for the packed stack-version order. C advances each version until its
// position passes last_position, the furthest position an earlier advance
// reached; a version other than version 0 also stops on reaching it. C merges
// versions in ts_parser__condense_stack only after every version has stopped.
// When version 0 shifts a zero-width token at last_position, it has not
// stopped: C advances it through the next token, which raises last_position,
// so every later version advances through that token as well before the
// versions meet. The Go loop dispatches every version against one lookahead
// per round, so after such a round its versions must not merge until a round
// passes last_position.
type packedVersionAdvance struct {
	// lastPosition is C's last_position: the furthest byte a round reached.
	lastPosition uint32
	// awaiting is true while version 0 has shifted a zero-width token at
	// lastPosition and no later round has passed it, so neither the
	// end-of-round condense nor the next round's merge may run.
	awaiting bool
}

// finishRound records the position a round's lookahead reached and decides
// whether the versions still wait to pass lastPosition. Rounds at end of
// input, rounds without a lookahead, and rounds with a version in recovery
// keep C's ordinary condense.
func (advance *packedVersionAdvance) finishRound(packedVersionOrder bool, lookahead Token, stacks []glrStack) {
	if lookahead.EndByte > advance.lastPosition {
		advance.lastPosition = lookahead.EndByte
		advance.awaiting = false
		return
	}
	if !packedVersionOrder || len(stacks) < 2 || lookahead.Symbol == 0 || lookahead.NoLookahead ||
		cRecoveryRelevantStack(stacks) {
		advance.awaiting = false
		return
	}
	if advance.awaiting {
		return
	}
	versionZero := &stacks[0]
	if lookahead.StartByte != lookahead.EndByte || versionZero.dead || versionZero.cPaused {
		return
	}
	top := stackEntryNode(versionZero.top())
	advance.awaiting = top != nil && top.symbol == lookahead.Symbol &&
		top.startByte == lookahead.StartByte && top.endByte == lookahead.EndByte
}
