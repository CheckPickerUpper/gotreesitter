//go:build !grammar_subset || grammar_subset_kotlin

package grammarruntime

import (
	"unicode"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// External token indexes for the kotlin grammar, in grammar.json order.
//
// fwcd/tree-sitter-kotlin@d6af96d2e23e added ten externals after
// _by_delegation_hint. Five of them are hints the scanner reads only as
// context flags in valid_symbols (_by_delegation_hint, _try_clause_hint,
// _backing_field_hint, _expression_end_hint and _call_lambda_hint); the
// others are emitted tokens.
const (
	kotlinTokAutoSemicolon          = 0  // "_automatic_semicolon"
	kotlinTokMultilineComment       = 1  // "multiline_comment"
	kotlinTokStringStart            = 2  // "_string_start"
	kotlinTokStringEnd              = 3  // "_string_end"
	kotlinTokStringContent          = 4  // "string_content"
	kotlinTokPrimaryCtorKW          = 5  // "_primary_constructor_keyword"
	kotlinTokImportDot              = 6  // "_import_dot"
	kotlinTokInterpExprStart        = 7  // "_interpolation_expression_start"
	kotlinTokInterpIdentStart       = 8  // "_interpolation_identifier_start"
	kotlinTokByDelegationHint       = 9  // "_by_delegation_hint"
	kotlinTokPrimaryCtorHint        = 10 // "_primary_constructor_hint"
	kotlinTokCallableReferenceQuest = 11 // "_callable_reference_quest"
	kotlinTokCommentContinuation    = 12 // "_comment_continuation_hint"
	kotlinTokTryClauseHint          = 13 // "_try_clause_hint"
	kotlinTokBackingFieldHint       = 14 // "_backing_field_hint"
	kotlinTokExpressionEndHint      = 15 // "_expression_end_hint"
	kotlinTokDelegateStartHint      = 16 // "_delegate_start_hint"
	kotlinTokDelegateEndHint        = 17 // "_delegate_end_hint"
	kotlinTokQuotedIdentifier       = 18 // "_quoted_identifier"
	kotlinTokCallLambdaHint         = 19 // "_call_lambda_hint"
	kotlinTokenCount                = 20
)

// kotlinDefaultSymTable holds the concrete symbol IDs of the generated
// grammar's ts_external_scanner_symbol_map, indexed by token.
var kotlinDefaultSymTable = [kotlinTokenCount]gotreesitter.Symbol{
	146, // _automatic_semicolon
	125, // multiline_comment
	147, // _string_start
	148, // _string_end
	149, // string_content
	150, // _primary_constructor_keyword
	151, // _import_dot
	152, // _interpolation_expression_start
	153, // _interpolation_identifier_start
	154, // _by_delegation_hint
	155, // _primary_constructor_hint
	156, // _callable_reference_quest
	157, // _comment_continuation_hint
	158, // _try_clause_hint
	159, // _backing_field_hint
	160, // _expression_end_hint
	161, // _delegate_start_hint
	162, // _delegate_end_hint
	163, // _quoted_identifier
	164, // _call_lambda_hint
}

var kotlinExternalScannerSpec = ExternalScannerSpec{
	Language:       "kotlin",
	UpstreamRepo:   "https://github.com/fwcd/tree-sitter-kotlin",
	UpstreamCommit: "d6af96d2e23ef0a8ff5ab512864378cf08fba311",
	SourceFiles: []ExternalScannerSourceFile{
		{Path: "src/grammar.json", SHA256: "b3b03a895464780a0a946b0a4baa7397599c4711ba8073357a78b69916f67ee1"},
		{Path: "src/scanner.c", SHA256: "a3c22a33953365af2b6a9408fd7c3828d4c1e83fb95a4ce29a854410ee16e99b"},
	},
	Externals: []string{
		"_automatic_semicolon",
		"multiline_comment",
		"_string_start",
		"_string_end",
		"string_content",
		"_primary_constructor_keyword",
		"_import_dot",
		"_interpolation_expression_start",
		"_interpolation_identifier_start",
		"_by_delegation_hint",
		"_primary_constructor_hint",
		"_callable_reference_quest",
		"_comment_continuation_hint",
		"_try_clause_hint",
		"_backing_field_hint",
		"_expression_end_hint",
		"_delegate_start_hint",
		"_delegate_end_hint",
		"_quoted_identifier",
		"_call_lambda_hint",
	},
}

func init() {
	RegisterExternalScannerSpec(kotlinExternalScannerSpec)
}

// String delimiter stack encoding, byte-for-byte the C scanner's: two bytes
// per open string, [delimiter | prefix_hi<<6, prefix_lo]. The delimiter is
// '"' for a single-quoted string and '"'+1 for a triple-quoted one; the
// 10-bit prefix is the number of '$' that trigger interpolation, or
// kotlinPrefixUnrepresentable for a prefix longer than kotlinPrefixLenMax.
const (
	kotlinDelimiterLength       = 3
	kotlinDelimiterMask         = 0x3F
	kotlinPrefixHiShift         = 6
	kotlinPrefixLenMax          = 1023
	kotlinPrefixUnrepresentable = 0
	kotlinSerializationSize     = 1024
	kotlinPendingHeaderBytes    = 11
	kotlinPendingFrameBits      = 11
	kotlinMaxStringFrames       = kotlinSerializationSize / 2
)

// kotlinPendingComment records a block comment whose ASI decision was made
// before the comment token itself was emitted. Offsets count UTF-8 bytes from
// the start of the scan that produced the decision.
type kotlinPendingComment struct {
	pending   bool
	startByte uint32
	endByte   uint32
	semiOwed  bool
}

// kotlinScannerState is the C scanner's Scanner struct.
type kotlinScannerState struct {
	strings            []byte
	constructorNewline bool
	comment            kotlinPendingComment
	inDelegate         bool
}

// kotlinScan carries one Scan call: the lexer, the valid token set indexed by
// kotlinTok*, and the token index the scan will emit.
type kotlinScan struct {
	lexer  *gotreesitter.ExternalLexer
	valid  *[kotlinTokenCount]bool
	result int
}

// KotlinExternalScanner implements gotreesitter.ExternalScanner for
// tree-sitter-kotlin. It is a line-for-line port of src/scanner.c at the
// commit named in kotlinExternalScannerSpec.
//
// Go's ExternalLexer reports end of input as a 0 lookahead, so a literal NUL
// byte in the source reads as end of input here, where the C scanner keeps
// consuming it as content.
type KotlinExternalScanner struct {
	symbols         [kotlinTokenCount]gotreesitter.Symbol
	externalToToken []int
}

// ExternalScannerForLanguage binds the scanner's tokens to lang's externals.
func (KotlinExternalScanner) ExternalScannerForLanguage(lang *gotreesitter.Language) gotreesitter.ExternalScanner {
	s := KotlinExternalScanner{symbols: kotlinDefaultSymTable}
	s.externalToToken = bindExternalScannerSpec(lang, kotlinExternalScannerSpec, func(tokenIdx int, sym gotreesitter.Symbol) {
		s.symbols[tokenIdx] = sym
	})
	return s
}

// Create returns an empty scanner state.
func (KotlinExternalScanner) Create() any {
	return &kotlinScannerState{}
}

// Destroy releases nothing; the state is garbage collected.
func (KotlinExternalScanner) Destroy(payload any) {}

// Serialize mirrors tree_sitter_kotlin_external_scanner_serialize.
func (KotlinExternalScanner) Serialize(payload any, buf []byte) int {
	s := payload.(*kotlinScannerState)
	n := len(s.strings)
	if n > kotlinSerializationSize || n%2 != 0 {
		return 0
	}
	if s.comment.pending {
		frames := n / 2
		length := kotlinPendingHeaderBytes + (frames*kotlinPendingFrameBits+7)/8
		if length > len(buf) {
			return 0
		}
		clear(buf[:length])
		var flags byte = 8
		if s.constructorNewline {
			flags |= 4
		}
		if s.comment.semiOwed {
			flags |= 16
		}
		if s.inDelegate {
			flags |= 32
		}
		buf[0] = flags
		buf[1] = byte(frames)
		buf[2] = byte(frames >> 8)
		kotlinWriteU32(buf[3:], s.comment.startByte)
		kotlinWriteU32(buf[7:], s.comment.endByte)
		for i := range frames {
			delimiter := s.strings[2*i]
			value := uint(delimiter>>kotlinPrefixHiShift)<<8 | uint(s.strings[2*i+1]) | uint(delimiter&1)<<10
			for bit := range kotlinPendingFrameBits {
				position := i*kotlinPendingFrameBits + bit
				if value&(1<<bit) != 0 {
					buf[kotlinPendingHeaderBytes+position/8] |= 1 << (position % 8)
				}
			}
		}
		return length
	}
	var flags byte
	if s.constructorNewline {
		flags |= 4
	}
	if s.inDelegate {
		flags |= 16
	}
	if n > 0 {
		if n > len(buf) {
			return 0
		}
		copy(buf, s.strings)
		buf[0] |= flags
		return n
	}
	if flags != 0 && len(buf) > 0 {
		buf[0] = flags
		return 1
	}
	return 0
}

// Deserialize mirrors tree_sitter_kotlin_external_scanner_deserialize.
func (KotlinExternalScanner) Deserialize(payload any, buf []byte) {
	s := payload.(*kotlinScannerState)
	s.strings = s.strings[:0]
	s.constructorNewline = false
	s.comment = kotlinPendingComment{}
	s.inDelegate = false
	length := len(buf)
	if length == 0 || length > kotlinSerializationSize {
		return
	}
	flags := buf[0]
	if flags&3 == 0 && flags&8 != 0 {
		if length < kotlinPendingHeaderBytes || uint(flags)&^60 != 0 {
			return
		}
		frames := int(buf[1]) | int(buf[2])<<8
		if frames > kotlinMaxStringFrames ||
			length != kotlinPendingHeaderBytes+(frames*kotlinPendingFrameBits+7)/8 {
			return
		}
		start, end := kotlinReadU32(buf[3:]), kotlinReadU32(buf[7:])
		if end < start || end-start < 4 {
			return
		}
		for i := range frames {
			var value uint
			for bit := range kotlinPendingFrameBits {
				position := i*kotlinPendingFrameBits + bit
				if buf[kotlinPendingHeaderBytes+position/8]&(1<<(position%8)) != 0 {
					value |= 1 << bit
				}
			}
			s.strings = kotlinStackPush(s.strings, '"', value&1024 != 0, uint16(value&1023))
		}
		s.constructorNewline = flags&4 != 0
		s.inDelegate = flags&32 != 0
		s.comment = kotlinPendingComment{pending: true, startByte: start, endByte: end, semiOwed: flags&16 != 0}
		return
	}
	if length == 1 && uint(flags)&^20 == 0 {
		s.constructorNewline = flags&4 != 0
		s.inDelegate = flags&16 != 0
		return
	}
	if length%2 != 0 {
		return
	}
	s.strings = append(s.strings, buf...)
	s.strings[0] &^= 20
	s.constructorNewline = flags&4 != 0
	s.inDelegate = flags&16 != 0
}

// Scan mirrors tree_sitter_kotlin_external_scanner_scan.
func (k KotlinExternalScanner) Scan(payload any, lexer *gotreesitter.ExternalLexer, validSymbols []bool) bool {
	s := payload.(*kotlinScannerState)
	var valid [kotlinTokenCount]bool
	if len(k.externalToToken) > 0 {
		for externalIdx, isValid := range validSymbols {
			if !isValid || externalIdx >= len(k.externalToToken) {
				continue
			}
			tokenIdx := k.externalToToken[externalIdx]
			if tokenIdx >= 0 && tokenIdx < kotlinTokenCount {
				valid[tokenIdx] = true
			}
		}
	} else {
		copy(valid[:], validSymbols)
	}
	scan := &kotlinScan{lexer: lexer, valid: &valid}
	found := scan.token(s)
	if !found {
		return false
	}
	if scan.result == kotlinTokAutoSemicolon {
		s.inDelegate = false
	}
	symbols := &k.symbols
	if k.symbols == ([kotlinTokenCount]gotreesitter.Symbol{}) {
		symbols = &kotlinDefaultSymTable
	}
	lexer.SetResultSymbol(symbols[scan.result])
	return true
}

// ---------------------------------------------------------------------------
// Lexer primitives
// ---------------------------------------------------------------------------

// skip consumes the lookahead outside the token.
func (c *kotlinScan) skip() { c.lexer.Advance(true) }

// advance consumes the lookahead into the token.
func (c *kotlinScan) advance() { c.lexer.Advance(false) }

// la returns the lookahead rune.
func (c *kotlinScan) la() rune { return c.lexer.Lookahead() }

// eof reports end of input (see the NUL note on KotlinExternalScanner).
func (c *kotlinScan) eof() bool { return c.lexer.Lookahead() == 0 }

// markEnd marks the token end at the cursor.
func (c *kotlinScan) markEnd() { c.lexer.MarkEnd() }

// advanceCounted advances, adding the lookahead's UTF-8 width to bytes.
func (c *kotlinScan) advanceCounted(bytes *uint32, whitespace bool) {
	if bytes != nil {
		width := kotlinUTF8Width(c.la())
		if *bytes <= ^uint32(0)-width {
			*bytes += width
		} else {
			*bytes = ^uint32(0)
		}
	}
	c.lexer.Advance(whitespace)
}

// kotlinUTF8Width is the encoded width of r, as C's utf8_width.
func kotlinUTF8Width(r rune) uint32 {
	switch {
	case r < 0x80:
		return 1
	case r < 0x800:
		return 2
	case r < 0x10000:
		return 3
	}
	return 4
}

// kotlinStackPush appends one packed delimiter entry.
func kotlinStackPush(stack []byte, chr byte, triple bool, prefixLen uint16) []byte {
	delimiter := chr
	if triple {
		delimiter = chr + 1
	}
	hi := byte(prefixLen>>8) & 0x3
	return append(stack, delimiter|hi<<kotlinPrefixHiShift, byte(prefixLen))
}

// kotlinStackPop drops the top entry; false when the stack is empty.
func kotlinStackPop(s *kotlinScannerState) bool {
	if len(s.strings) < 2 {
		return false
	}
	s.strings = s.strings[:len(s.strings)-2]
	return true
}

// kotlinWriteU32 writes value little-endian.
func kotlinWriteU32(buf []byte, value uint32) {
	for i := range 4 {
		buf[i] = byte(value >> (8 * i))
	}
}

// kotlinReadU32 reads a little-endian value.
func kotlinReadU32(buf []byte) uint32 {
	var value uint32
	for i := range 4 {
		value |= uint32(buf[i]) << (8 * i)
	}
	return value
}

// kotlinIsSpace is C's iswspace.
func kotlinIsSpace(r rune) bool { return unicode.IsSpace(r) }

// kotlinIsWordChar is C's is_word_char: iswalnum or '_'.
func kotlinIsWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// kotlinIsAlpha is C's iswalpha.
func kotlinIsAlpha(r rune) bool { return unicode.IsLetter(r) }

// ---------------------------------------------------------------------------
// Strings and comments
// ---------------------------------------------------------------------------

// scanStringStart mirrors scan_string_start.
func (c *kotlinScan) scanStringStart(s *kotlinScannerState) bool {
	if len(s.strings)+2 > kotlinSerializationSize {
		return false
	}
	var dollars uint32
	for c.la() == '$' {
		c.advance()
		if dollars <= kotlinPrefixLenMax {
			dollars++
		}
	}
	var prefixLen uint16
	switch {
	case dollars == 0:
		prefixLen = 1
	case dollars > kotlinPrefixLenMax:
		prefixLen = kotlinPrefixUnrepresentable
	default:
		prefixLen = uint16(dollars)
	}
	if c.la() != '"' {
		return false
	}
	c.advance()
	c.markEnd()
	for count := 1; count < kotlinDelimiterLength; count++ {
		if c.la() != '"' {
			s.strings = kotlinStackPush(s.strings, '"', false, prefixLen)
			return true
		}
		c.advance()
	}
	c.markEnd()
	s.strings = kotlinStackPush(s.strings, '"', true, prefixLen)
	return true
}

// scanStringContent mirrors scan_string_content.
func (c *kotlinScan) scanStringContent(s *kotlinScannerState) bool {
	if len(s.strings) < 2 {
		return false
	}
	packed := s.strings[len(s.strings)-2]
	prefixLen := uint16(packed>>kotlinPrefixHiShift)<<8 | uint16(s.strings[len(s.strings)-1])
	triggerLen := uint32(prefixLen)
	if prefixLen == kotlinPrefixUnrepresentable {
		triggerLen = kotlinPrefixLenMax + 1
	}
	rawDelimiter := packed & kotlinDelimiterMask
	isTriple := rawDelimiter&1 != 0
	endChar := rune(rawDelimiter)
	if isTriple {
		endChar = rune(rawDelimiter - 1)
	}
	hasContent := false
	for {
		if c.eof() {
			break
		}
		switch la := c.la(); {
		case la == '$':
			if hasContent {
				c.result = kotlinTokStringContent
				return true
			}
			c.advance()
			c.markEnd()
			var additionalDollars uint32
			for c.la() == '$' {
				c.advance()
				if additionalDollars <= kotlinPrefixLenMax {
					additionalDollars++
				}
			}
			totalDollars := 1 + additionalDollars
			if totalDollars >= triggerLen && (kotlinIsAlpha(c.la()) || c.la() == '_' || c.la() == '{') {
				if totalDollars > triggerLen {
					c.result = kotlinTokStringContent
					return true
				}
				if additionalDollars > 0 {
					c.markEnd()
				}
				if c.valid[kotlinTokInterpExprStart] && c.la() == '{' {
					c.advance()
					if c.la() == '}' {
						return false
					}
					c.markEnd()
					c.result = kotlinTokInterpExprStart
					return true
				}
				if c.valid[kotlinTokInterpIdentStart] && (kotlinIsAlpha(c.la()) || c.la() == '_') {
					c.result = kotlinTokInterpIdentStart
					return true
				}
				return false
			}
			if additionalDollars > 0 {
				c.markEnd()
			}
			c.result = kotlinTokStringContent
			return true
		case la == '\\':
			c.advance()
			if c.la() == '$' {
				c.advance()
				if c.la() == endChar {
					if !kotlinStackPop(s) {
						return false
					}
					c.advance()
					c.markEnd()
					c.result = kotlinTokStringEnd
					return true
				}
			} else if isTriple && c.la() == endChar {
				hasContent = true
				continue
			}
		case la == endChar:
			if isTriple {
				c.markEnd()
				for count := 1; count < kotlinDelimiterLength; count++ {
					c.advance()
					if c.la() != endChar {
						c.markEnd()
						c.result = kotlinTokStringContent
						return true
					}
				}
				if hasContent && c.la() == endChar {
					c.result = kotlinTokStringContent
					return true
				}
				c.result = kotlinTokStringEnd
				c.markEnd()
				for c.la() == endChar {
					c.advance()
					c.markEnd()
				}
				return kotlinStackPop(s)
			}
			if hasContent {
				c.markEnd()
				c.result = kotlinTokStringContent
				return true
			}
			if !kotlinStackPop(s) {
				return false
			}
			c.advance()
			c.markEnd()
			c.result = kotlinTokStringEnd
			return true
		}
		c.advance()
		hasContent = true
	}
	return false
}

// scanMultilineComment mirrors scan_multiline_comment.
func (c *kotlinScan) scanMultilineComment(newline *bool, bytes *uint32) bool {
	if c.la() != '/' {
		return false
	}
	c.advanceCounted(bytes, false)
	if c.la() != '*' {
		return false
	}
	c.advanceCounted(bytes, false)
	afterStar := false
	nestingDepth := 1
	for {
		if newline != nil && (c.la() == '\n' || c.la() == '\r') {
			*newline = true
		}
		switch c.la() {
		case '*':
			c.advanceCounted(bytes, false)
			afterStar = true
		case '/':
			c.advanceCounted(bytes, false)
			if afterStar {
				afterStar = false
				nestingDepth--
				if nestingDepth == 0 {
					c.result = kotlinTokMultilineComment
					c.markEnd()
					return true
				}
			} else {
				afterStar = false
				if c.la() == '*' {
					nestingDepth++
					c.advanceCounted(bytes, false)
				}
			}
		case 0:
			c.result = kotlinTokMultilineComment
			c.markEnd()
			return true
		default:
			c.advanceCounted(bytes, false)
			afterStar = false
		}
	}
}

// ---------------------------------------------------------------------------
// Word and trivia probes
// ---------------------------------------------------------------------------

// scanForWord mirrors scan_for_word: skip the already-checked first
// character, then match word and a word boundary.
func (c *kotlinScan) scanForWord(word string) bool {
	c.skip()
	for i := 0; i < len(word); i++ {
		if c.la() != rune(word[i]) {
			return false
		}
		c.skip()
	}
	return !kotlinIsWordChar(c.la())
}

// checkWord mirrors check_word.
func (c *kotlinScan) checkWord(word string) bool {
	for i := 0; i < len(word); i++ {
		if c.la() != rune(word[i]) {
			return false
		}
		c.skip()
	}
	return !kotlinIsWordChar(c.la())
}

// skipWhitespaceAndComments mirrors skip_whitespace_and_comments.
func (c *kotlinScan) skipWhitespaceAndComments() bool {
	for {
		for kotlinIsSpace(c.la()) {
			c.skip()
		}
		if c.la() != '/' {
			return true
		}
		c.skip()
		switch c.la() {
		case '/':
			c.skip()
			for c.la() != '\n' && c.la() != '\r' && !c.eof() {
				c.skip()
			}
		case '*':
			c.skip()
			depth := 1
			for depth > 0 && !c.eof() {
				switch c.la() {
				case '*':
					c.skip()
					if c.la() == '/' {
						c.skip()
						depth--
					}
				case '/':
					c.skip()
					if c.la() == '*' {
						c.skip()
						depth++
					}
				default:
					c.skip()
				}
			}
		default:
			return false
		}
	}
}

// followedByArrow mirrors followed_by_arrow.
func (c *kotlinScan) followedByArrow() bool {
	if !c.skipWhitespaceAndComments() {
		return false
	}
	if c.la() != '-' {
		return false
	}
	c.skip()
	return c.la() == '>'
}

// kotlinVisibilityModifiers are the words check_modifier_then_* accept.
var kotlinVisibilityModifiers = map[string]bool{
	"internal": true, "private": true, "protected": true, "public": true,
}

// checkModifierThen mirrors check_modifier_then_constructor and
// check_modifier_then_field: a visibility modifier, horizontal whitespace,
// then keyword.
func (c *kotlinScan) checkModifierThen(keyword string) bool {
	word := make([]rune, 0, 19)
	for kotlinIsWordChar(c.la()) && len(word) < 19 {
		word = append(word, c.la())
		c.skip()
	}
	if !kotlinVisibilityModifiers[string(word)] {
		return false
	}
	for c.la() == ' ' || c.la() == '\t' {
		c.skip()
	}
	return c.checkWord(keyword)
}

// checkAnnotationThen mirrors check_annotation_then_constructor and
// check_annotation_then_field: annotations, an optional visibility modifier,
// then keyword.
func (c *kotlinScan) checkAnnotationThen(keyword string) bool {
	for c.la() == '@' {
		c.skip()
		if !kotlinIsWordChar(c.la()) {
			return false
		}
		for kotlinIsWordChar(c.la()) {
			c.skip()
		}
		for c.la() == '.' {
			c.skip()
			if !kotlinIsWordChar(c.la()) {
				break
			}
			for kotlinIsWordChar(c.la()) {
				c.skip()
			}
		}
		if c.la() == '(' {
			depth := 1
			c.skip()
			for depth > 0 && c.la() != 0 && !c.eof() {
				if c.la() == '"' {
					c.skip()
					for c.la() != '"' && c.la() != 0 && !c.eof() {
						if c.la() == '\\' {
							c.skip()
						}
						c.skip()
					}
					if c.la() == '"' {
						c.skip()
					}
					continue
				}
				if c.la() == '(' {
					depth++
				} else if c.la() == ')' {
					depth--
				}
				c.skip()
			}
		}
		for kotlinIsSpace(c.la()) {
			c.skip()
		}
	}
	if kotlinIsWordChar(c.la()) && c.la() != rune(keyword[0]) {
		return c.checkModifierThen(keyword)
	}
	return c.checkWord(keyword)
}

// kotlinConstructorSpace mirrors constructor_space.
func kotlinConstructorSpace(r rune) bool {
	return r == ' ' || (r >= '\t' && r <= '\r') || r == 0x85 || r == 0xa0 ||
		r == 0x1680 || (r >= 0x2000 && r <= 0x200a) || r == 0x2028 ||
		r == 0x2029 || r == 0x202f || r == 0x205f || r == 0x3000
}

// constructorTrivia mirrors constructor_trivia: advance over whitespace and
// comments, noting line breaks; false on a bare '/' or an unterminated block
// comment.
func (c *kotlinScan) constructorTrivia(newline *bool) bool {
	for {
		for kotlinConstructorSpace(c.la()) {
			if c.la() == '\n' || c.la() == '\r' {
				*newline = true
			}
			c.advance()
		}
		if c.la() != '/' {
			return true
		}
		c.advance()
		switch c.la() {
		case '/':
			for !c.eof() && c.la() != '\n' && c.la() != '\r' {
				c.advance()
			}
		case '*':
			c.advance()
			depth := 1
			for depth > 0 && !c.eof() {
				r := c.la()
				if r == '\n' || r == '\r' {
					*newline = true
				}
				c.advance()
				if r == '/' && c.la() == '*' {
					c.advance()
					depth++
				} else if r == '*' && c.la() == '/' {
					c.advance()
					depth--
				}
			}
			if depth > 0 {
				return false
			}
		default:
			return false
		}
	}
}

// kotlinConstructorWordChar mirrors constructor_word_char.
func kotlinConstructorWordChar(r rune) bool {
	if r < 128 {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	}
	return !kotlinConstructorSpace(r)
}

// scanTryClauseKeyword mirrors scan_try_clause_keyword.
func (c *kotlinScan) scanTryClauseKeyword(suffix string) bool {
	return c.scanForWord(suffix) && !kotlinConstructorWordChar(c.la())
}

// kotlinWordCapacity is the C probes' char buffer size, terminator included.
const kotlinWordCapacity = 32

// constructorWord mirrors constructor_word with a 32-byte buffer: it
// returns the (possibly truncated) word and whether it was a whole ASCII
// word that fit.
func (c *kotlinScan) constructorWord() (string, bool) {
	word := make([]byte, 0, kotlinWordCapacity)
	ascii := true
	for !c.eof() && kotlinConstructorWordChar(c.la()) {
		r := c.la()
		if r >= 128 || len(word)+1 >= kotlinWordCapacity {
			ascii = false
		}
		if len(word)+1 < kotlinWordCapacity {
			word = append(word, byte(r))
		}
		c.advance()
	}
	return string(word), ascii && len(word) > 0
}

// kotlinConstructorModifiers are the words constructor_modifier accepts.
var kotlinConstructorModifiers = map[string]bool{
	"abstract": true, "actual": true, "annotation": true, "const": true,
	"crossinline": true, "data": true, "expect": true, "external": true,
	"final": true, "infix": true, "inline": true, "inner": true,
	"internal": true, "lateinit": true, "noinline": true, "open": true,
	"operator": true, "override": true, "private": true, "protected": true,
	"public": true, "sealed": true, "suspend": true, "tailrec": true,
	"value": true, "vararg": true,
}

// kotlinProbeFrameKind says what a constructorBalanced frame is inside.
type kotlinProbeFrameKind uint8

const (
	kotlinProbeBracket kotlinProbeFrameKind = iota
	kotlinProbeQuoted
	kotlinProbeTripleString
	kotlinProbeBacktick
)

// kotlinProbeFrame is C's ConstructorProbeFrame.
type kotlinProbeFrame struct {
	end     rune
	kind    kotlinProbeFrameKind
	dollars uint
}

// constructorBalanced mirrors constructor_balanced.
func (c *kotlinScan) constructorBalanced(close rune) bool {
	stack := []kotlinProbeFrame{{end: close}}
	c.advance()
	for len(stack) > 0 && !c.eof() {
		frame := stack[len(stack)-1]
		r := c.la()
		if frame.kind != kotlinProbeBracket {
			if frame.kind == kotlinProbeQuoted && r == '\\' {
				c.advance()
				if !c.eof() {
					c.advance()
				}
				continue
			}
			if r == frame.end {
				c.advance()
				if frame.kind != kotlinProbeTripleString {
					stack = stack[:len(stack)-1]
					continue
				}
				if c.la() == '"' {
					c.advance()
					if c.la() == '"' {
						c.advance()
						stack = stack[:len(stack)-1]
					}
				}
				continue
			}
			if frame.end == '"' && r == '$' {
				var run uint
				for c.la() == '$' {
					c.advance()
					run++
				}
				if run >= frame.dollars && c.la() == '{' {
					c.advance()
					stack = append(stack, kotlinProbeFrame{end: '}'})
				}
				continue
			}
			c.advance()
			continue
		}
		if r == frame.end {
			c.advance()
			stack = stack[:len(stack)-1]
			continue
		}
		if r == '(' || r == '[' || r == '{' || (r == '<' && frame.end == '>') {
			var end rune
			switch r {
			case '(':
				end = ')'
			case '[':
				end = ']'
			case '{':
				end = '}'
			default:
				end = '>'
			}
			stack = append(stack, kotlinProbeFrame{end: end})
			c.advance()
			continue
		}
		if r == '/' {
			newline := false
			c.constructorTrivia(&newline)
			continue
		}
		var dollars uint = 1
		if r == '$' {
			dollars = 0
			for c.la() == '$' {
				c.advance()
				dollars++
			}
			r = c.la()
		}
		if r == '"' || r == '\'' || r == '`' {
			kind := kotlinProbeQuoted
			if r == '`' {
				kind = kotlinProbeBacktick
			}
			c.advance()
			if r == '"' && c.la() == '"' {
				c.advance()
				if c.la() != '"' {
					continue // Empty ordinary string.
				}
				c.advance()
				kind = kotlinProbeTripleString
			}
			stack = append(stack, kotlinProbeFrame{end: r, kind: kind, dollars: dollars})
			continue
		}
		c.advance()
	}
	return len(stack) == 0
}

// constructorAnnotation mirrors constructor_annotation.
func (c *kotlinScan) constructorAnnotation() bool {
	c.advance() // @
	newline := false
	if !c.constructorTrivia(&newline) {
		return false
	}
	if c.la() == '[' {
		return c.constructorBalanced(']')
	}
	for {
		if !c.constructorTrivia(&newline) {
			return false
		}
		if c.la() == '`' {
			c.advance()
			for !c.eof() && c.la() != '`' {
				c.advance()
			}
			if c.eof() {
				return false
			}
			c.advance()
		} else {
			if !kotlinConstructorWordChar(c.la()) {
				return false
			}
			c.constructorWord()
		}
		if !c.constructorTrivia(&newline) {
			return false
		}
		if c.la() == ':' {
			c.advance() // use-site target
			if !c.constructorTrivia(&newline) {
				return false
			}
			if c.la() == '[' {
				return c.constructorBalanced(']')
			}
			continue
		}
		if c.la() == '.' {
			c.advance()
			continue
		}
		break
	}
	if !c.constructorTrivia(&newline) {
		return false
	}
	if c.la() == '<' && !c.constructorBalanced('>') {
		return false
	}
	if !c.constructorTrivia(&newline) {
		return false
	}
	if c.la() == '(' && !c.constructorBalanced(')') {
		return false
	}
	return true
}

// whereTypeConstraintFollows mirrors where_type_constraint_follows.
func (c *kotlinScan) whereTypeConstraintFollows(newline *bool) bool {
	if !c.constructorTrivia(newline) {
		return false
	}
	for c.la() == '@' {
		if !c.constructorAnnotation() {
			return false
		}
		if !c.constructorTrivia(newline) {
			return false
		}
	}
	if c.la() == '`' {
		c.advance()
		for !c.eof() && c.la() != '`' && c.la() != '\n' && c.la() != '\r' {
			c.advance()
		}
		if c.la() != '`' {
			return false
		}
		c.advance()
	} else {
		if !kotlinConstructorWordChar(c.la()) {
			return false
		}
		for kotlinConstructorWordChar(c.la()) {
			c.advance()
		}
	}
	if !c.constructorTrivia(newline) {
		return false
	}
	if c.la() != ':' {
		return false
	}
	c.advance()
	return c.la() != ':'
}

// braceStartsStatement mirrors brace_starts_statement.
func (c *kotlinScan) braceStartsStatement(inDelegate bool) bool {
	return c.valid[kotlinTokExpressionEndHint] && !inDelegate && !c.valid[kotlinTokStringContent]
}

// scanPrimaryConstructorHint mirrors scan_primary_constructor_hint.
func (c *kotlinScan) scanPrimaryConstructorHint(s *kotlinScannerState) bool {
	newline := s.constructorNewline
	s.constructorNewline = false
	for kotlinConstructorSpace(c.la()) {
		if c.la() == '\n' || c.la() == '\r' {
			newline = true
		}
		c.skip()
	}
	if c.la() == '/' {
		if c.valid[kotlinTokMultilineComment] && c.scanMultilineComment(&newline, nil) {
			s.constructorNewline = newline
			return true
		}
		return false
	}
	c.markEnd()
	if c.la() == ';' {
		c.advance()
		c.markEnd()
		c.result = kotlinTokAutoSemicolon
		return c.valid[kotlinTokAutoSemicolon]
	}
	if c.eof() {
		c.result = kotlinTokAutoSemicolon
		return c.valid[kotlinTokAutoSemicolon]
	}
	headerNewline := newline
	first := c.la()
	firstWord := ""
	firstIdentifier := true
	for {
		if c.la() == '@' {
			if !c.constructorAnnotation() {
				break
			}
		} else {
			word, ascii := c.constructorWord()
			if firstIdentifier {
				firstWord = word
				firstIdentifier = false
			}
			if !ascii {
				break
			}
			if word == "constructor" {
				if c.constructorTrivia(&newline) && c.la() == '(' {
					c.result = kotlinTokPrimaryCtorHint
					return true
				}
				break
			}
			if !kotlinConstructorModifiers[word] {
				break
			}
		}
		if !c.constructorTrivia(&newline) {
			break
		}
	}
	continuation := false
	switch first {
	case '(', '{', '<', '.', '?', ',', '*', '%', '>', '=', '|', '&':
		continuation = true
	}
	if firstWord == "as" || (c.valid[kotlinTokByDelegationHint] && firstWord == "by") {
		continuation = true
	}
	if firstWord == "where" {
		continuation = c.whereTypeConstraintFollows(&newline)
	}
	if first == ':' {
		c.advance()
		continuation = c.la() != ':'
	}
	if firstWord == "finally" {
		continuation = c.valid[kotlinTokTryClauseHint]
	}
	if first == '!' {
		c.advance()
		continuation = c.la() == '='
	}
	if firstWord == "else" {
		continuation = true
		if c.constructorTrivia(&newline) && c.la() == '-' {
			c.advance()
			continuation = c.la() != '>'
		}
	}
	c.result = kotlinTokAutoSemicolon
	ends := firstWord == "import"
	if headerNewline {
		ends = !continuation
	}
	return ends && c.valid[kotlinTokAutoSemicolon]
}

// trailingLambdaFollows mirrors trailing_lambda_follows; it returns the
// leading word the probe read alongside its verdict.
func (c *kotlinScan) trailingLambdaFollows() (bool, string) {
	leadingWord := ""
	annotated := false
	newline := false
	if !c.constructorTrivia(&newline) {
		return false, leadingWord
	}
	for c.la() == '@' {
		annotated = true
		if !c.constructorAnnotation() || !c.constructorTrivia(&newline) {
			return false, leadingWord
		}
	}
	if c.la() == '{' {
		return true, leadingWord
	}
	if c.la() == '`' {
		c.advance()
		for !c.eof() && c.la() != '`' {
			c.advance()
		}
		if c.eof() {
			return false, leadingWord
		}
		c.advance()
	} else {
		if !kotlinConstructorWordChar(c.la()) {
			return false, leadingWord
		}
		word, ascii := c.constructorWord()
		if !annotated && ascii {
			leadingWord = word
		}
	}
	if !c.constructorTrivia(&newline) || c.la() != '@' {
		return false, leadingWord
	}
	c.advance()
	return c.constructorTrivia(&newline) && c.la() == '{', leadingWord
}

// callPrefixContinues mirrors call_prefix_continues.
func (c *kotlinScan) callPrefixContinues() bool {
	follows, word := c.trailingLambdaFollows()
	if follows {
		return true
	}
	switch word {
	case "as":
		return true
	case "by":
		return c.valid[kotlinTokByDelegationHint]
	case "where":
		newline := false
		return c.whereTypeConstraintFollows(&newline)
	case "catch", "finally":
		return c.valid[kotlinTokTryClauseHint]
	case "field":
		return c.valid[kotlinTokBackingFieldHint]
	case "else":
		return !c.followedByArrow()
	}
	return false
}

// callLambdaProbeApplies is the shared guard of the three call_prefix_continues
// call sites.
func (c *kotlinScan) callLambdaProbeApplies() bool {
	return c.valid[kotlinTokCallLambdaHint] &&
		(c.la() == '@' || kotlinConstructorWordChar(c.la()) || c.la() == '`')
}

// asiAfterComment mirrors asi_after_comment.
func (c *kotlinScan) asiAfterComment(inDelegate bool) bool {
	if c.callLambdaProbeApplies() {
		return !c.callPrefixContinues()
	}
	switch c.la() {
	case '{':
		return c.braceStartsStatement(inDelegate)
	case '.', ',', '*', '%', '>', '<', '=', '?', '|', '&', '/':
		return false
	case ':':
		c.skip()
		return c.la() == ':'
	case '!':
		c.skip()
		return c.la() != '='
	case 'e':
		if c.scanForWord("lse") {
			return c.followedByArrow()
		}
		return true
	case 'a':
		return !c.scanForWord("s")
	case 'w':
		if c.scanForWord("here") {
			newline := false
			return !c.whereTypeConstraintFollows(&newline)
		}
		return true
	case 'c':
		return !(c.valid[kotlinTokTryClauseHint] && c.scanTryClauseKeyword("atch"))
	case 'b':
		return !(c.valid[kotlinTokByDelegationHint] && c.scanForWord("y"))
	case 'f':
		if c.valid[kotlinTokBackingFieldHint] && c.scanTryClauseKeyword("ield") {
			return false
		}
		return !(c.valid[kotlinTokTryClauseHint] && c.scanTryClauseKeyword("inally"))
	}
	return true
}

// continuationAfterBlockComment mirrors continuation_after_block_comment.
func (c *kotlinScan) continuationAfterBlockComment(inDelegate bool) bool {
	newline := false
	if !c.constructorTrivia(&newline) {
		return !c.eof()
	}
	if c.callLambdaProbeApplies() {
		return c.callPrefixContinues()
	}
	switch c.la() {
	case '{':
		return !c.braceStartsStatement(inDelegate)
	case '.', ',', '*', '%', '>', '<', '=', '?', '|', '&':
		return true
	case ':':
		c.advance()
		return c.la() != ':'
	case '!':
		c.advance()
		return c.la() == '='
	}
	word, ok := c.constructorWord()
	if !ok {
		return false
	}
	switch word {
	case "as":
		return true
	case "by":
		return c.valid[kotlinTokByDelegationHint]
	case "where":
		return c.whereTypeConstraintFollows(&newline)
	case "catch", "finally":
		return c.valid[kotlinTokTryClauseHint]
	case "field":
		return c.valid[kotlinTokBackingFieldHint]
	case "else":
		if !c.constructorTrivia(&newline) || c.la() != '-' {
			return true
		}
		c.advance()
		return c.la() != '>'
	}
	return false
}

// resumePendingComment mirrors resume_pending_comment.
func (c *kotlinScan) resumePendingComment(s *kotlinScannerState) bool {
	if s.comment.semiOwed {
		if !c.valid[kotlinTokAutoSemicolon] {
			return false
		}
		c.markEnd()
		c.result = kotlinTokAutoSemicolon
		s.comment.semiOwed = false
		return true
	}
	if !c.valid[kotlinTokMultilineComment] {
		return false
	}
	length := s.comment.endByte - s.comment.startByte
	var consumed uint32
	for kotlinIsSpace(c.la()) {
		c.skip()
	}
	matched := c.scanMultilineComment(nil, &consumed)
	s.comment = kotlinPendingComment{}
	return matched && consumed == length
}

// scanAutomaticSemicolon mirrors scan_automatic_semicolon. cursorClean is
// set only by returns that leave the cursor at the first non-whitespace
// character, so the caller may still lex a token there.
func (c *kotlinScan) scanAutomaticSemicolon(s *kotlinScannerState, cursorClean *bool) bool {
	c.result = kotlinTokAutoSemicolon
	c.markEnd()

	var prefixBytes uint32
	sameline := true
	for {
		if c.eof() {
			return true
		}
		if c.la() == ';' {
			c.advance()
			c.markEnd()
			return true
		}
		if !kotlinIsSpace(c.la()) {
			break
		}
		if c.la() == '\n' {
			c.advanceCounted(&prefixBytes, true)
			sameline = false
			break
		}
		if c.la() == '\r' {
			c.advanceCounted(&prefixBytes, true)
			if c.la() == '\n' {
				c.advanceCounted(&prefixBytes, true)
			}
			sameline = false
			break
		}
		c.advanceCounted(&prefixBytes, true)
	}

	for kotlinIsSpace(c.la()) {
		c.advanceCounted(&prefixBytes, true)
	}

	if sameline {
		switch c.la() {
		case 'i':
			return c.scanForWord("mport")
		case ';':
			c.advance()
			c.markEnd()
			return true
		}
		*cursorClean = true
		return false
	}

	if c.callLambdaProbeApplies() {
		return !c.callPrefixContinues()
	}

	switch c.la() {
	case '{':
		if c.braceStartsStatement(s.inDelegate) {
			return true
		}
		*cursorClean = true
		return false
	case ',', '.', '*', '%', '>', '<', '=', '?', '|', '&':
		*cursorClean = true
		return false
	case ':':
		c.skip()
		return c.la() == ':'
	case '/':
		return c.automaticSemicolonBeforeSlash(s, prefixBytes)
	case '+', '-':
		return true
	case '!':
		c.skip()
		return c.la() != '='
	case 'b':
		return !(c.valid[kotlinTokByDelegationHint] && c.scanForWord("y"))
	case 'e':
		if !c.scanForWord("lse") {
			return true
		}
		return c.followedByArrow()
	case 'a':
		return !c.scanForWord("s")
	case 'w':
		if !c.scanForWord("here") {
			return true
		}
		newline := false
		return !c.whereTypeConstraintFollows(&newline)
	case 'i', 'p':
		if c.valid[kotlinTokPrimaryCtorKW] && !c.valid[kotlinTokStringContent] &&
			c.checkModifierThen("constructor") {
			return false
		}
		if c.valid[kotlinTokBackingFieldHint] && !c.valid[kotlinTokStringContent] &&
			c.checkModifierThen("field") {
			return false
		}
		return true
	case 'c':
		if c.valid[kotlinTokPrimaryCtorKW] && !c.valid[kotlinTokStringContent] {
			if c.matchConstructorKeyword() {
				c.result = kotlinTokPrimaryCtorKW
				c.markEnd()
			}
			return true
		}
		return !(c.valid[kotlinTokTryClauseHint] && c.scanTryClauseKeyword("atch"))
	case 'f':
		if c.valid[kotlinTokBackingFieldHint] && c.scanTryClauseKeyword("ield") {
			return false
		}
		return !(c.valid[kotlinTokTryClauseHint] && c.scanTryClauseKeyword("inally"))
	case '@':
		if c.valid[kotlinTokPrimaryCtorKW] && !c.valid[kotlinTokStringContent] &&
			c.checkAnnotationThen("constructor") {
			return false
		}
		if c.valid[kotlinTokBackingFieldHint] && !c.valid[kotlinTokStringContent] &&
			c.checkAnnotationThen("field") {
			return false
		}
		return true
	case ';':
		c.advance()
		c.markEnd()
		return true
	}
	return true
}

// automaticSemicolonBeforeSlash is the '/' arm of scan_automatic_semicolon:
// division, a line comment, or a block comment after a line break.
func (c *kotlinScan) automaticSemicolonBeforeSlash(s *kotlinScannerState, prefixBytes uint32) bool {
	c.advance()
	switch c.la() {
	case '/':
		c.skip()
		for c.la() != '\n' && c.la() != '\r' && c.la() != 0 && !c.eof() {
			c.skip()
		}
		if !c.skipWhitespaceAndComments() {
			return false
		}
		return c.asiAfterComment(s.inDelegate)
	case '*':
		c.advance()
		var commentBytes uint32 = 2
		nestingDepth := 1
		nested := false
		hasNUL := false
		for nestingDepth > 0 && !c.eof() {
			switch c.la() {
			case '*':
				c.advanceCounted(&commentBytes, false)
				if c.la() == '/' {
					c.advanceCounted(&commentBytes, false)
					nestingDepth--
				}
			case '/':
				c.advanceCounted(&commentBytes, false)
				if c.la() == '*' {
					c.advanceCounted(&commentBytes, false)
					nestingDepth++
					nested = true
				}
			default:
				if c.la() == 0 {
					hasNUL = true
				}
				c.advanceCounted(&commentBytes, false)
			}
		}
		if nestingDepth > 0 {
			return true
		}
		if nested || hasNUL {
			if c.continuationAfterBlockComment(s.inDelegate) &&
				c.valid[kotlinTokCommentContinuation] &&
				!c.valid[kotlinTokStringContent] &&
				prefixBytes <= ^uint32(0)-commentBytes {
				s.comment = kotlinPendingComment{pending: true, startByte: prefixBytes, endByte: prefixBytes + commentBytes}
				c.result = kotlinTokCommentContinuation
			}
			return true
		}
		if !c.skipWhitespaceAndComments() {
			return false
		}
		return c.asiAfterComment(s.inDelegate)
	}
	return false
}

// matchConstructorKeyword advances over "constructor" and reports whether it
// matched as a whole word, as the C scanner's inline loops do.
func (c *kotlinScan) matchConstructorKeyword() bool {
	const keyword = "constructor"
	for i := 0; i < len(keyword); i++ {
		if c.la() != rune(keyword[i]) {
			return false
		}
		c.advance()
	}
	return !kotlinConstructorWordChar(c.la())
}

// scanImportDot mirrors scan_import_dot.
func (c *kotlinScan) scanImportDot() bool {
	if c.la() != '.' {
		return false
	}
	c.markEnd()
	c.advance()
	foundNewline := false
	for kotlinIsSpace(c.la()) {
		if c.la() == '\n' || c.la() == '\r' {
			foundNewline = true
		}
		c.skip()
	}
	if foundNewline && c.la() == 'i' && c.scanForWord("mport") {
		c.result = kotlinTokAutoSemicolon
		return true
	}
	c.result = kotlinTokImportDot
	c.markEnd()
	return true
}

// token mirrors scan_token.
func (c *kotlinScan) token(s *kotlinScannerState) bool {
	if s.comment.pending {
		return c.resumePendingComment(s)
	}
	if c.valid[kotlinTokDelegateStartHint] && !c.valid[kotlinTokStringContent] {
		c.markEnd()
		s.inDelegate = true
		c.result = kotlinTokDelegateStartHint
		return true
	}
	if c.valid[kotlinTokDelegateEndHint] && !c.valid[kotlinTokStringContent] {
		c.markEnd()
		s.inDelegate = false
		c.result = kotlinTokDelegateEndHint
		return true
	}
	if c.valid[kotlinTokPrimaryCtorHint] && !c.valid[kotlinTokStringContent] {
		return c.scanPrimaryConstructorHint(s)
	}
	s.constructorNewline = false
	if c.valid[kotlinTokAutoSemicolon] {
		cursorClean := false
		if c.scanAutomaticSemicolon(s, &cursorClean) {
			return true
		}
		if !cursorClean && !c.valid[kotlinTokStringContent] {
			return false
		}
	}

	if c.valid[kotlinTokQuotedIdentifier] && !c.valid[kotlinTokStringContent] {
		for kotlinIsSpace(c.la()) {
			c.skip()
		}
		if c.la() == '`' {
			c.advance()
			content := false
			for !c.eof() && c.la() != '`' && c.la() != '\n' && c.la() != '\r' {
				content = true
				c.advance()
			}
			if content && c.la() == '`' {
				c.advance()
				c.markEnd()
				c.result = kotlinTokQuotedIdentifier
				return true
			}
			return false
		}
	}

	if c.valid[kotlinTokCallableReferenceQuest] && c.la() == '?' {
		c.advance()
		c.markEnd()
		for c.la() == ' ' || c.la() == '\t' {
			c.skip()
		}
		if c.la() == '?' {
			c.result = kotlinTokCallableReferenceQuest
			return true
		}
		if c.la() == ':' {
			c.advance()
			if c.la() == ':' {
				c.result = kotlinTokCallableReferenceQuest
				return true
			}
		}
		return false
	}

	if c.valid[kotlinTokImportDot] && c.scanImportDot() {
		return true
	}

	if c.valid[kotlinTokPrimaryCtorKW] && !c.valid[kotlinTokStringContent] {
		for kotlinIsSpace(c.la()) {
			c.skip()
		}
		if c.la() == 'c' && c.matchConstructorKeyword() {
			c.result = kotlinTokPrimaryCtorKW
			c.markEnd()
			return true
		}
	}

	if c.valid[kotlinTokStringContent] || c.valid[kotlinTokInterpExprStart] ||
		c.valid[kotlinTokInterpIdentStart] {
		if c.scanStringContent(s) {
			return true
		}
	}

	for kotlinIsSpace(c.la()) {
		c.skip()
	}

	if c.valid[kotlinTokStringStart] && c.scanStringStart(s) {
		c.result = kotlinTokStringStart
		return true
	}

	if c.valid[kotlinTokMultilineComment] && c.scanMultilineComment(nil, nil) {
		return true
	}

	return false
}
