package grammargen

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// separatorGrammar is a list of words whose extras are the given rules.
func separatorGrammar(extras ...*Rule) *Grammar {
	g := NewGrammar("separators")
	g.Define("document", Repeat(Sym("word")))
	g.Define("word", Pat(`[a-z]+`))
	g.Define("comment", Token(Seq(Str("#"), Pat(`[a-z]*`))))
	g.SetExtras(extras...)
	return g
}

// TestLexerSkipsOnlyTheGrammarsSeparators parses inputs under three kinds of extras and expects
// the lexer to skip exactly what tree-sitter's separators match: nothing when the extras only
// name a rule, what a pattern extra matches, and a string extra's text, which leaves no node.
func TestLexerSkipsOnlyTheGrammarsSeparators(t *testing.T) {
	for _, row := range []struct {
		name      string
		extras    []*Rule
		clean     string
		children  int
		erroneous []string
	}{
		{name: "rule only", extras: []*Rule{Sym("comment")}, clean: "a#b", children: 2, erroneous: []string{"a b", "a\nb"}},
		{name: "pattern", extras: []*Rule{Pat(`\n`)}, clean: "a\nb", children: 2, erroneous: []string{"a b", "a\tb"}},
		{name: "string", extras: []*Rule{Str("~")}, clean: "a~b", children: 2, erroneous: []string{"a b"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			lang, err := GenerateLanguage(separatorGrammar(row.extras...))
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			parse := func(sample string) *gotreesitter.Tree {
				tree, err := gotreesitter.NewParser(lang).Parse([]byte(sample))
				if err != nil {
					t.Fatalf("parse %q: %v", sample, err)
				}
				t.Cleanup(tree.Release)
				return tree
			}
			root := parse(row.clean).RootNode()
			if root.HasError() || root.ChildCount() != row.children {
				t.Errorf("%q = %s, want %d children and no error", row.clean, root.SExpr(lang), row.children)
			}
			for _, sample := range row.erroneous {
				if !parse(sample).RootNode().HasError() {
					t.Errorf("%q parsed without an error, want what no separator matches left unskipped", sample)
				}
			}
		})
	}
}
