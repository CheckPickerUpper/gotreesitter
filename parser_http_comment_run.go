package gotreesitter

import "bytes"

// parseHTTPCommentRun reproduces the locked C parser's section boundaries
// for a document made only of ordinary hash comments. Its grammar permits
// ambiguous boundaries; C retains at most eight comments per section.
// Empty comments have different boundaries and stay on the DFA route.
func (p *Parser) parseHTTPCommentRun(source []byte) (*Tree, bool) {
	if p == nil || p.language == nil || p.language.Name != "http" || len(p.included) != 0 ||
		p.parseWorkLimits.configured() || len(source) == 0 {
		return nil, false
	}
	lines := 0
	for start := 0; start < len(source); {
		if start+2 > len(source) || source[start] != '#' || source[start+1] != ' ' {
			return nil, false
		}
		end := bytes.IndexByte(source[start:], '\n')
		if end < 0 || len(bytes.TrimSpace(source[start+2:start+end])) == 0 {
			return nil, false
		}
		start += end + 1
		lines++
	}
	if lines <= 8 {
		return nil, false
	}
	commentSymbol, commentOK := p.language.visibleSymbolByNameAndNamed("comment", true)
	sectionSymbol, sectionOK := p.language.visibleSymbolByNameAndNamed("section", true)
	documentSymbol, documentOK := p.language.visibleSymbolByNameAndNamed("document", true)
	if !commentOK || !sectionOK || !documentOK {
		return nil, false
	}
	sectionCount := (lines + 7) / 8
	arena := acquireNodeArena(arenaClassFull)
	budget := parseMemoryBudgetForParser(p, len(source))
	sections := make([]*Node, 0, sectionCount)
	groupSize := lines % 8
	if groupSize == 0 {
		groupSize = 8
	}
	offset := 0
	row := uint32(0)
	for len(sections) < sectionCount {
		if reason := p.parseStopReasonNow(); parseStopReasonIsTerminal(reason) {
			arena.Release()
			return nil, false
		}
		comments := make([]*Node, 0, groupSize)
		for i := 0; i < groupSize; i++ {
			end := offset + bytes.IndexByte(source[offset:], '\n') + 1
			comments = append(comments, newLeafNodeInArena(arena, commentSymbol, true,
				uint32(offset), uint32(end), Point{Row: row}, Point{Row: row + 1}))
			offset = end
			row++
		}
		sections = append(sections, newParentNodeInArena(arena, sectionSymbol, true, comments, nil, 0))
		groupSize = 8
		if budget > 0 && arena.allocatedBytes > budget {
			arena.Release()
			return nil, false
		}
	}
	root := newParentNodeInArena(arena, documentSymbol, true, sections, nil, 0)
	if budget > 0 && arena.allocatedBytes > budget {
		arena.Release()
		return nil, false
	}
	tree := newTreeWithArenas(root, source, p.language, arena, nil)
	runtime := forestAcceptedRuntime(root, source)
	runtime.NodesAllocated = lines + sectionCount + 1
	runtime.TokensConsumed = uint64(lines + 1)
	runtime.Iterations = lines + 1
	runtime.MaxStacksSeen = 1
	runtime.ArenaBytesAllocated = arena.allocatedBytes
	runtime.MemoryBudgetBytes = budget
	tree.setParseRuntime(runtime)
	tree.forestFastPath = true
	tree.incrementalReuseDisabled = true
	return tree, true
}
