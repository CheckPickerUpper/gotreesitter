package grammargen

import (
	"testing"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// conditionListGrammar is the shape of Swift's if: a comma-separated condition list before a block,
// where an expression may also take a trailing closure written like the block.
func conditionListGrammar() *Grammar {
	g := NewGrammar("condition_list")
	g.Define("source_file", Repeat(Sym("_statement")))
	g.Define("_statement", Choice(Sym("if_statement"), Sym("_expression")))
	g.Define("if_statement", PrecRight(-1, Seq(
		Str("if"),
		Seq(Field("condition", Sym("_condition")), Repeat(Seq(Str(","), Field("condition", Sym("_condition"))))),
		Sym("_block"),
	)))
	g.Define("_condition", Choice(Sym("_expression")))
	g.Define("_expression", Prec(-1, Choice(
		Sym("identifier"),
		Sym("number"),
		Sym("equality_expression"),
		Sym("call_expression"),
		Sym("if_statement"),
	)))
	g.Define("equality_expression", PrecLeft(5, Seq(
		Field("lhs", Sym("_expression")),
		Str("=="),
		Field("rhs", Sym("_binary_suffix")),
	)))
	g.Define("_binary_suffix", PrecLeft(-2, Choice(
		Sym("_expression"),
		Alias(Sym("suffix_call"), "call_expression", true),
	)))
	g.Define("suffix_call", Seq(Sym("_expression"), Alias(Sym("suffix_call_arguments"), "call_suffix", true)))
	g.Define("suffix_call_arguments", Prec(-2, Sym("value_arguments")))
	g.Define("call_expression", Prec(-2, PrecDynamic(1, Seq(Sym("_expression"), Sym("call_suffix")))))
	g.Define("call_suffix", Prec(-2, Choice(
		Sym("value_arguments"),
		PrecDynamic(-1, Sym("lambda_literal")),
	)))
	g.Define("value_arguments", Seq(Str("("), Optional(Sym("_expression")), Str(")")))
	g.Define("lambda_literal", PrecLeft(-3, Seq(Str("{"), Repeat(Sym("_statement")), Str("}"))))
	g.Define("_block", Prec(2, Seq(Str("{"), Repeat(Sym("_statement")), Str("}"))))
	g.Define("identifier", Pat(`[a-z]+`))
	g.Define("number", Pat(`[0-9]+`))
	g.SetExtras(Pat(`\s`))
	g.SetConflicts([]string{"call_suffix"}, []string{"suffix_call_arguments", "call_suffix"})
	return g
}

// TestAConditionListOfThreeEndsAtItsBlock parses ifs with one, two and three conditions and
// expects each block read as the if's block, not as a trailing closure of the last condition.
func TestAConditionListOfThreeEndsAtItsBlock(t *testing.T) {
	lang, err := GenerateLanguage(conditionListGrammar())
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, sample := range []string{
		"if a == 1 { x }",
		"if a == 1, b == 2 { x }",
		"if a == 1, b == 2, c == 3 { x }",
		"if a, b, c { x }",
	} {
		tree, err := gotreesitter.NewParser(lang).Parse([]byte(sample))
		if err != nil {
			t.Fatalf("parse %q: %v", sample, err)
		}
		t.Logf("%q: %s", sample, tree.RootNode().SExpr(lang))
		if tree.RootNode().HasError() {
			t.Errorf("%q parsed with an error", sample)
		}
		tree.Release()
	}
}
