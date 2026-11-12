package gotreesitter

// RepairNoLookaheadLexModes marks parser states as no-lookahead when they only
// need EOF-triggered reductions plus trivia handling. Tree-sitter's C runtime
// uses these states to reduce before lexing the next real token. A
// no-lookahead state reduces by its EOF actions without lexing, so a state is
// marked only when every other token, internal or external, is trivia there:
// a real action on an external token needs the external scanner to run.
func RepairNoLookaheadLexModes(lang *Language) {
	if lang == nil || len(lang.LexModes) == 0 || len(lang.ParseActions) == 0 {
		return
	}

	for state := 0; state < len(lang.LexModes) && state < int(lang.StateCount); state++ {
		if lang.LexModes[state].LexStateIndex() == noLookaheadLexState {
			continue
		}

		eofIdx := lookupRepairActionIndex(lang, StateID(state), 0)
		if eofIdx == 0 || int(eofIdx) >= len(lang.ParseActions) {
			continue
		}
		eofEntry := lang.ParseActions[eofIdx]
		if len(eofEntry.Actions) == 0 {
			continue
		}

		allReduce := true
		for _, act := range eofEntry.Actions {
			if act.Type != ParseActionReduce {
				allReduce = false
				break
			}
		}
		if !allReduce {
			continue
		}

		needsRealLookahead := false
		for sym := Symbol(1); uint32(sym) < lang.TokenCount; sym++ {
			idx := lookupRepairActionIndex(lang, StateID(state), sym)
			if idx == 0 || int(idx) >= len(lang.ParseActions) {
				continue
			}
			if repairEntryIsTriviaOnly(lang.ParseActions[idx]) {
				continue
			}
			needsRealLookahead = true
			break
		}

		if !needsRealLookahead {
			lang.LexModes[state].SetLexStateIndex(noLookaheadLexState)
		}
	}
}

func repairEntryIsTriviaOnly(entry ParseActionEntry) bool {
	if len(entry.Actions) == 0 {
		return false
	}
	for _, act := range entry.Actions {
		if !act.Extra && !act.ExtraChain {
			return false
		}
	}
	return true
}

func lookupRepairActionIndex(lang *Language, state StateID, sym Symbol) uint16 {
	if lang == nil {
		return 0
	}

	denseLimit := languageDenseLimit(lang)
	if int(state) < denseLimit {
		if int(state) >= len(lang.ParseTable) {
			return 0
		}
		row := lang.ParseTable[state]
		if int(sym) >= len(row) {
			return 0
		}
		return row[sym]
	}

	smallIdx := int(state) - int(lang.LargeStateCount)
	if smallIdx < 0 || smallIdx >= len(lang.SmallParseTableMap) {
		return 0
	}
	table := lang.SmallParseTable
	offset := lang.SmallParseTableMap[smallIdx]
	if int(offset) >= len(table) {
		return 0
	}
	groupCount := table[offset]
	pos := int(offset) + 1
	for i := uint16(0); i < groupCount; i++ {
		if pos+1 >= len(table) {
			break
		}
		sectionValue := table[pos]
		symbolCount := table[pos+1]
		pos += 2
		for j := uint16(0); j < symbolCount; j++ {
			if pos >= len(table) {
				break
			}
			if Symbol(table[pos]) == sym {
				return sectionValue
			}
			pos++
		}
	}
	return 0
}
