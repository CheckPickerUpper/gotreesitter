package grammargen

// LALR(1) gives every LR(0) state one lookahead set per completed item: the
// union over every state the item's production can have started in (its
// lanes). A rule built on a generated repeat, such as a list of statements,
// is reduced in a state that also shifts the token that starts the repeat's
// next element. When that token reaches the reduce through only some lanes
// (it follows the rule in one context and not in another), resolving the
// conflict for the lanes that have it, by the repeat's associativity, also
// drops the shift for the lanes that never had the token, which canonical
// LR(1), as tree-sitter builds it, keeps. Such a cell keeps the dropped shift
// beside the resolved actions, so the parser forks there and the lane that
// cannot reduce on the token continues the repeat.

// laneDependentConflictCell is one (state, token) action cell of a merged
// LALR state whose resolution would drop a shift some lanes need.
type laneDependentConflictCell struct {
	state int
	token int
}

// findLaneDependentConflicts returns every action cell of tables where a
// completed item's lookahead comes from only some of its lanes, resolving the
// cell keeps that item's reduce and drops another action, and a dropped shift
// continues a generated repeat that the reduced rule's productions use.
func (ctx *lrContext) findLaneDependentConflicts(tables *LRTables) map[laneDependentConflictCell]bool {
	lanesByItem := make(map[[2]int][]int)
	for _, lookback := range ctx.lalrLookbacks {
		key := [2]int{int(lookback.stateIdx), int(lookback.coreIdx)}
		lanesByItem[key] = append(lanesByItem[key], int(lookback.ntIdx))
	}
	repeatsByRule := generatedRepeatsByRule(ctx.ng)
	cells := make(map[laneDependentConflictCell]bool)
	for key, lanes := range lanesByItem {
		if len(lanes) < 2 {
			continue
		}
		state, coreIdx := key[0], key[1]
		core := ctx.itemSets[state].cores[coreIdx]
		prodIdx := int(core.prodIdx)
		repeats := repeatsByRule[ctx.ng.Productions[prodIdx].LHS]
		if len(repeats) == 0 {
			continue
		}
		core.lookaheads.forEach(func(token int) {
			actions := tables.ActionTable[state][token]
			if !shiftContinuesRepeat(actions, repeats) || !ctx.someLaneLacks(lanes, token) {
				return
			}
			if resolutionDropsActionBesideReduce(token, actions, prodIdx, ctx.ng) {
				cells[laneDependentConflictCell{state: state, token: token}] = true
			}
		})
	}
	return cells
}

// generatedRepeatsByRule maps each rule to the generated repeat helpers its
// productions use.
func generatedRepeatsByRule(ng *NormalizedGrammar) map[int]map[int]bool {
	repeatsByRule := make(map[int]map[int]bool)
	for _, production := range ng.Productions {
		for _, symbol := range production.RHS {
			if symbol >= len(ng.Symbols) || !ng.Symbols[symbol].GeneratedRepeatAux {
				continue
			}
			repeats := repeatsByRule[production.LHS]
			if repeats == nil {
				repeats = make(map[int]bool)
				repeatsByRule[production.LHS] = repeats
			}
			repeats[symbol] = true
		}
	}
	return repeatsByRule
}

// shiftContinuesRepeat reports whether a shift in actions comes from an item
// of one of repeats.
func shiftContinuesRepeat(actions []lrAction, repeats map[int]bool) bool {
	for _, action := range actions {
		if action.kind != lrShift {
			continue
		}
		if repeats[action.lhsSym] {
			return true
		}
		for _, symbol := range action.lhsSyms {
			if repeats[symbol] {
				return true
			}
		}
	}
	return false
}

// someLaneLacks reports whether token is missing from the follow set of at
// least one of lanes.
func (ctx *lrContext) someLaneLacks(lanes []int, token int) bool {
	for _, lane := range lanes {
		if !ctx.lalrFollowSets[lane].contains(token) {
			return true
		}
	}
	return false
}

// resolutionDropsActionBesideReduce reports whether resolving actions keeps
// the reduce of prodIdx and drops some other action.
func resolutionDropsActionBesideReduce(token int, actions []lrAction, prodIdx int, ng *NormalizedGrammar) bool {
	if len(actions) < 2 {
		return false
	}
	resolved, err := resolveActionConflict(token, actions, ng)
	if err != nil || len(resolved) >= len(actions) {
		return false
	}
	return lrActionListHasReduce(resolved, prodIdx)
}

// resolveActionCell resolves the conflicting actions of one action cell of
// tables, keeping every shift resolution drops when the cell is lane
// dependent.
func (tables *LRTables) resolveActionCell(state, token int, actions []lrAction, ng *NormalizedGrammar) ([]lrAction, error) {
	resolved, err := resolveActionConflict(token, actions, ng)
	if err != nil || !tables.LaneDependentConflicts[laneDependentConflictCell{state: state, token: token}] {
		return resolved, err
	}
	return keepDroppedShifts(actions, resolved), nil
}

// keepDroppedShifts returns resolved followed by every shift of actions that
// resolution dropped, in their original order. A dropped reduce stays
// dropped: only the shift is what the lanes without the token need.
func keepDroppedShifts(actions, resolved []lrAction) []lrAction {
	kept := append([]lrAction(nil), resolved...)
	for _, action := range actions {
		if action.kind == lrShift && !lrActionListHasAction(resolved, action) {
			kept = append(kept, action)
		}
	}
	return kept
}

// lrActionListHasAction reports whether actions holds an action of the same
// kind with the same target state or production as action.
func lrActionListHasAction(actions []lrAction, action lrAction) bool {
	for _, candidate := range actions {
		if candidate.kind == action.kind && candidate.state == action.state && candidate.prodIdx == action.prodIdx {
			return true
		}
	}
	return false
}
