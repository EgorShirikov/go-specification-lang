package front

import (
	acsl "gpma/framework/gosl/front/generated/grammar/acsl"
	grammar "gpma/framework/gosl/front/generated/grammar/acsl"
	specir "gpma/framework/gosl/ir"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	log "github.com/sirupsen/logrus"
)

func (b *ACSLBuilder) VisitAcsl_comment(ctx *acsl.Acsl_commentContext) interface{} {
	switch {
	case ctx.Acsl_comment_contract() != nil:
		return b.Visit(ctx.Acsl_comment_contract())

	case ctx.Acsl_comment_loop_annot() != nil:
		return b.Visit(ctx.Acsl_comment_loop_annot())

	case ctx.Acsl_comment_assertion() != nil:
		return b.Visit(ctx.Acsl_comment_assertion())

	case ctx.Acsl_comment_statement_contract() != nil:
		return b.Visit(ctx.Acsl_comment_statement_contract())

	default:
		return nil
	}
}

func (b *ACSLBuilder) VisitAcsl_comment_contract(ctx *acsl.Acsl_comment_contractContext) interface{} {
	fc := &specir.FuncContract{
		FuncName:     "",
		RecvTypeName: "",
		Span:         b.blockSpan(ctx),
	}

	b.Mod.FuncContracts = append(b.Mod.FuncContracts, fc)
	// opening function contract context
	b.currFunc = fc
	// closing function contract context at VisitAcsl_comment_contract ending
	defer func() { b.currFunc = nil }()

	fcCtx := ctx.Function_contract()
	if fcCtx != nil {
		b.Visit(fcCtx)
	}

	fc.StartPos, fc.EndPos = b.blockPositions(fcCtx)

	log.Debug("[ACSLBuilder] added FuncContract. Text:", ctx.GetText())
	// log.Debug("[ACSLBuilder] added FuncContract:", fc.FuncName, "Text:", ctx.GetText(), "comment_span:", b.blockSpan(ctx), "contract_span:", b.blockSpan(fcCtx))
	// log.Debug("contract_start:", fc.StartPos, "contract_end:", fc.EndPos)
	// log.Debug("block_start:", b.Block.Start, "block_end:", b.Block.End)
	return nil
}

func (b *ACSLBuilder) VisitAcsl_comment_loop_annot(ctx *acsl.Acsl_comment_loop_annotContext) interface{} {
	ls := &specir.LoopSpec{
		ParentFuncName: "",
		Span:           b.blockSpan(ctx),
	}

	b.Mod.LoopSpecs = append(b.Mod.LoopSpecs, ls)

	log.Debug("[ACSLBuilder] added LoopSpec in:", ls.ParentFuncName, "Text:", ctx.GetText())
	return nil
}

func (b *ACSLBuilder) VisitAcsl_comment_assertion(ctx *acsl.Acsl_comment_assertionContext) interface{} {
	log.Debug("[ACSLBuilder] assertion (stub)", "Text:", ctx.GetText())
	return nil
}

func (b *ACSLBuilder) VisitAcsl_comment_statement_contract(ctx *acsl.Acsl_comment_statement_contractContext) interface{} {
	log.Debug("[ACSLBuilder] statement contract (stub)", "Text:", ctx.GetText())
	return nil
}

// Function contract: requires/ensures/assigns
func (b *ACSLBuilder) VisitFunction_contract(ctx *acsl.Function_contractContext) interface{} {
	log.Debug("VisitFunction_contract text:", ctx.GetText())
	if b.currFunc == nil || ctx == nil {
		return nil
	}

	// requires
	for _, rc := range ctx.AllRequires_clause() {
		np, ok := b.Visit(rc).(specir.NamedPredicate)
		if ok {
			b.currFunc.Requires = append(b.currFunc.Requires, np)
		}
	}

	// simple_clause usually contains ensures/assigns ...
	for _, sc := range ctx.AllSimple_clause() {
		b.Visit(sc)
	}

	return nil
}

func (b *ACSLBuilder) VisitSimple_clause(ctx *acsl.Simple_clauseContext) interface{} {
	log.Debug("VisitSimple_clause text:", ctx.GetText())
	if b.currFunc == nil || ctx == nil {
		return nil
	}
	if ec := ctx.Ensures_clause(); ec != nil {
		np, ok := b.Visit(ec).(specir.NamedPredicate)
		if ok {
			b.currFunc.Ensures = append(b.currFunc.Ensures, np)
		}
		return nil
	}
	if ac := ctx.Assigns_clause(); ac != nil {
		ls, ok := b.Visit(ac).(specir.LocationSet)
		if ok {
			b.currFunc.Assigns = ls
		}
		return nil
	}
	return nil
}

// requires_clause / ensures_clause -> NamedPredicate
func (b *ACSLBuilder) VisitRequires_clause(ctx *acsl.Requires_clauseContext) interface{} {
	log.Debug("VisitRequires_clause text:", ctx.GetText())
	if ctx == nil {
		return nil
	}
	return b.parseNamedPred(ctx.Pred())
}

func (b *ACSLBuilder) VisitEnsures_clause(ctx *acsl.Ensures_clauseContext) interface{} {
	log.Debug("VisitEnsures_clause text:", ctx.GetText())
	if ctx == nil {
		return nil
	}
	return b.parseNamedPred(ctx.Pred())
}

func (b *ACSLBuilder) parseNamedPred(p acsl.IPredContext) specir.NamedPredicate {
	log.Debug("parseNamedPred text:", p.GetText())
	// default without label
	np := specir.NamedPredicate{
		Label: "",
		Expr:  nil,
		Span:  b.blockSpan(p), // TODO
	}
	if p == nil {
		return np
	}

	// syntactic naming: "bound: <pred>"
	if sn, ok := p.(*acsl.Syntactic_naming_predContext); ok {
		// label can be String_ or Id(). Grammar dependent
		if sn.String_() != nil {
			np.Label = sn.String_().GetText()
		} else if sn.Id() != nil {
			np.Label = sn.Id().GetText()
		}
		np.Expr = b.predToExpr(sn.Pred())
		return np
	}

	// Simple pred
	np.Expr = b.predToExpr(p)
	return np
}

// assigns_clause -> LocationSet
func (b *ACSLBuilder) VisitAssigns_clause(ctx *acsl.Assigns_clauseContext) interface{} {
	log.Debug("VisitAssigns_clause text:", ctx.GetText())
	if ctx == nil {
		return specir.LocationSet{Kind: specir.LocSetNothing}
	}
	txt := ctx.GetText()

	// \nothing
	if strings.Contains(txt, `\nothing`) {
		return specir.LocationSet{Kind: specir.LocSetNothing}
	}

	// collecting identifiers. TODO
	vars := collectIdentifiers(ctx)
	if len(vars) == 0 {
		// fallback
		return specir.LocationSet{Kind: specir.LocSetVars, Vars: []string{txt}}
	}
	return specir.LocationSet{Kind: specir.LocSetVars, Vars: unique(vars)}
}

// subtree walk and collecting Identifier tokens TODO
func collectIdentifiers(t antlr.Tree) []string {
	var out []string
	var walk func(n antlr.Tree)
	walk = func(n antlr.Tree) {
		if n == nil {
			return
		}
		if tn, ok := n.(antlr.TerminalNode); ok {
			sym := tn.GetSymbol()
			if sym != nil && sym.GetTokenType() == acsl.ACSLParserIdentifier {
				out = append(out, sym.GetText())
			}
		}
		for i := 0; i < n.GetChildCount(); i++ {
			if c, ok := n.GetChild(i).(antlr.Tree); ok {
				walk(c)
			}
		}
	}
	walk(t)
	return out
}

func unique(xs []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	return out
}

// specir.Expr
func (b *ACSLBuilder) predToExpr(p acsl.IPredContext) specir.Expr {
	log.Debugf("predToExpr text:%s, type:%T\n", p.GetText(), p)
	if p == nil {
		return nil
	}

	v := b.Visit(p)
	if e, ok := v.(specir.Expr); ok {
		return e
	}
	if np, ok := v.(specir.NamedPredicate); ok {
		// returning NamedPredicate
		return np.Expr
	}
	// fallback
	log.Debug("predToExpr fallback")
	return &specir.ExprVar{BaseExpr: specir.BaseExpr{Span: b.blockSpan(p)}, Name: p.GetText()}
}

func (b *ACSLBuilder) VisitLogical_true_pred(ctx *acsl.Logical_true_predContext) interface{} {
	log.Debug("VisitLogical_true_pred text:", ctx.GetText())
	return &specir.ExprLiteral{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Kind:     specir.LitBool,
		Bool:     true,
	}
}

func (b *ACSLBuilder) VisitLogical_false_pred(ctx *acsl.Logical_false_predContext) interface{} {
	log.Debug("VisitLogical_false_pred text:", ctx.GetText())
	return &specir.ExprLiteral{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Kind:     specir.LitBool,
		Bool:     false,
	}
}

// binary op
func (b *ACSLBuilder) VisitBinary_op_term(ctx *acsl.Binary_op_termContext) interface{} {
	log.Debug("VisitBinary_op_term text:", ctx.GetText())
	left := b.termToExpr(ctx.Term(0))
	right := b.termToExpr(ctx.Term(1))

	opText := ctx.Bin_op().GetText()
	op := relOpToBinary(opText)

	return &specir.ExprBinary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// a <= b, a == b, ...
func (b *ACSLBuilder) VisitComparison_pred(ctx *acsl.Comparison_predContext) interface{} {
	log.Debug("VisitComparison_pred text:", ctx.GetText())
	//Term Rel_op Term
	terms := ctx.AllTerm()
	ops := ctx.AllRel_op()
	if len(terms) < 2 || len(ops) < 1 {
		return &specir.ExprVar{BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)}, Name: ctx.GetText()}
	}

	left := b.termToExpr(terms[0])
	right := b.termToExpr(terms[1])
	log.Debugf("VisitComparison_pred l: %T r: %T\n", terms[0], terms[1])
	op := relOpToBinary(ops[0].GetText())
	return &specir.ExprBinary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

func relOpToBinary(s string) specir.BinaryOp {
	switch s {
	case "==":
		return specir.OpEq
	case "!=":
		return specir.OpNeq
	case "<=":
		return specir.OpLe
	case ">=":
		return specir.OpGe
	case "<":
		return specir.OpLt
	case ">":
		return specir.OpGt
	default:
		return specir.OpEq
	}
}

// && / ||
func (b *ACSLBuilder) VisitConjunction_pred(ctx *acsl.Conjunction_predContext) interface{} {
	log.Debug("VisitConjunction_pred text:", ctx.GetText())
	ps := ctx.AllPred()
	if len(ps) < 2 {
		return b.predToExpr(ctx.Pred(0))
	}
	return &specir.ExprBinary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       specir.OpAnd,
		Left:     b.predToExpr(ps[0]),
		Right:    b.predToExpr(ps[1]),
	}
}

func (b *ACSLBuilder) VisitDisjunction_pred(ctx *acsl.Disjunction_predContext) interface{} {
	log.Debug("VisitDisjunction_pred text:", ctx.GetText())
	ps := ctx.AllPred()
	if len(ps) < 2 {
		return b.predToExpr(ctx.Pred(0))
	}
	return &specir.ExprBinary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       specir.OpOr,
		Left:     b.predToExpr(ps[0]),
		Right:    b.predToExpr(ps[1]),
	}
}

// pred: cond ? p1 : p2
func (b *ACSLBuilder) VisitTernary_condition_pred(ctx *acsl.Ternary_condition_predContext) interface{} {
	log.Debug("VisitTernary_condition_pred text:", ctx.GetText())
	ps := ctx.AllPred()
	if len(ps) != 3 {
		return &specir.ExprVar{BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)}, Name: ctx.GetText()}
	}

	cond := b.predToExpr(ps[0])
	thenE := b.predToExpr(ps[1])
	elseE := b.predToExpr(ps[2])
	return &specir.ExprConditional{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Cond:     cond,
		Then:     thenE,
		Else:     elseE,
	}
}

// \old(pred)
func (b *ACSLBuilder) VisitOld_pred(ctx *acsl.Old_predContext) interface{} {
	log.Debug("VisitOld_pred text:", ctx.GetText())
	inner := b.predToExpr(ctx.Pred())
	return &specir.ExprOld{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Inner:    inner,
	}
}

// terms
func (b *ACSLBuilder) termToExpr(t acsl.ITermContext) specir.Expr {
	if t == nil {
		return nil
	}
	v := b.Visit(t)
	if e, ok := v.(specir.Expr); ok {
		return e
	}
	return &specir.ExprVar{BaseExpr: specir.BaseExpr{Span: b.blockSpan(t)}, Name: t.GetText()}
}

func (b *ACSLBuilder) VisitInternal_term(ctx *acsl.Internal_termContext) interface{} {
	log.Debug("VisitInternal_term text:", ctx.GetText())
	// "id" / literal / ...
	return b.VisitChildren(ctx)
}

func (b *ACSLBuilder) VisitFunc_application_term(ctx *acsl.Func_application_termContext) interface{} {
	log.Debug("VisitFunc_application_term text:", ctx.GetText())
	// term ::= ident '(' term (',' term)* ')'
	name := ""
	// if ctx.Id() != nil {
	//     name = ctx.Id().GetText()
	// }
	argsCtx := ctx.AllTerm()
	args := make([]specir.Expr, 0, len(argsCtx))
	for _, a := range argsCtx {
		args = append(args, b.termToExpr(a))
	}
	return &specir.ExprCall{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Name:     name,
		Args:     args,
	}
}

func (b *ACSLBuilder) VisitOld_term(ctx *acsl.Old_termContext) interface{} {
	log.Debug("VisitOld_term text:", ctx.GetText())
	inner := b.termToExpr(ctx.Term())
	return &specir.ExprOld{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Inner:    inner,
	}
}

func (b *ACSLBuilder) VisitResult_term(ctx *acsl.Result_termContext) interface{} {
	log.Debug("VisitResult_term text:", ctx.GetText())
	return &specir.ExprResult{BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)}}
}

func (b *ACSLBuilder) VisitParentheses_pred(ctx *acsl.Parentheses_predContext) interface{} {
	log.Debug("VisitParentheses_pred text:", ctx.GetText())
	return b.predToExpr(ctx.Pred())
}

func (b *ACSLBuilder) VisitUnary_op_term(ctx *acsl.Unary_op_termContext) interface{} {
	log.Debugf("VisitUnary_op_term text:%s, type:%T\n", ctx.GetText(), ctx.Term())
	opCtx := ctx.Unary_op()
	inner := b.termToExpr(ctx.Term())

	var op specir.UnaryOp

	switch opCtx.GetText() {
	case "!":
		op = specir.UnaryNot

	case "-":
		op = specir.UnaryMinus

	case "+":
		op = specir.UnaryPlus

	default:
		// fallback
		return &specir.ExprVar{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Name:     ctx.GetText(),
		}
	}

	return &specir.ExprUnary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       op,
		Expr:     inner,
	}
}

func (b *ACSLBuilder) VisitLiteral_term(ctx *grammar.Literal_termContext) interface{} {
	if ctx == nil {
		return nil
	}

	lit := ctx.Literal()

	log.Debugf("VisitLiteral_term text:%s, type:%T\n", ctx.GetText(), lit)
	if lit == nil {
		// fallback
		return &specir.ExprVar{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Name:     ctx.GetText(),
		}
	}

	switch x := lit.(type) {

	case *grammar.True_constantContext:
		return &specir.ExprLiteral{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Kind:     specir.LitBool,
			Bool:     true,
		}

	case *grammar.False_constantContext:
		return &specir.ExprLiteral{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Kind:     specir.LitBool,
			Bool:     false,
		}

	case *grammar.String_constantContext:
		raw := ""
		if s := x.String_(); s != nil {
			raw = s.GetText()
		} else {
			raw = x.GetText()
		}

		// removing ".."
		unq, err := strconv.Unquote(raw)
		if err != nil {
			unq = raw // fallback
		}

		return &specir.ExprLiteral{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Kind:     specir.LitString,
			Str:      unq,
		}

	case *grammar.Trivial_constantContext:
		// Constant()
		raw := ""
		if c := x.Constant(); c != nil {
			raw = c.GetText()
		} else {
			raw = x.GetText()
		}

		// int (base=0)
		if i, err := strconv.ParseInt(raw, 0, 64); err == nil {
			return &specir.ExprLiteral{
				BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
				Kind:     specir.LitInteger,
				Int:      i,
			}
		}

		// float
		if f, err := strconv.ParseFloat(raw, 64); err == nil {
			return &specir.ExprLiteral{
				BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
				Kind:     specir.LitFloat,
				Float:    f,
			}
		}

		// fallback
		return &specir.ExprVar{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Name:     raw,
		}

	default:
		// TODO: extra literals
		return &specir.ExprVar{
			BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
			Name:     lit.(antlr.ParserRuleContext).GetText(),
		}
	}
}

func (b *ACSLBuilder) VisitNegation_pred(ctx *acsl.Negation_predContext) interface{} {
	return &specir.ExprUnary{
		BaseExpr: specir.BaseExpr{Span: b.blockSpan(ctx)},
		Op:       specir.UnaryNot,
		Expr:     b.predToExpr(ctx.Pred()),
	}
}
