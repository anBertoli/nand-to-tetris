package parser

import (
	"../tokenizer"
)

func (p* Parser) nextIsExpression () bool {
	return p.nextIsTerm()
}
func (p* Parser) parseExpression() tokenizer.Token {
	var toBuild []tokenizer.Token

	term := p.parseTerm()
	toBuild = append(toBuild, term)

	r, isRune := p.LookNext().Value.(rune)
	if isRune && tokenizer.IsSymbolOp(r) {
		p.ConsumeNext(&toBuild)
		term := p.parseTerm()
		toBuild = append(toBuild, term)
	}

	return tokenizer.Token{
		Type: EXPRESSION,
		Children: &toBuild,
	}
}



func (p* Parser) nextIsTerm() bool {
	var firstToken = p.LookNext()

	if firstToken.Type == tokenizer.INTCONST || firstToken.Type == tokenizer.STRINGCONST {
		return true
	}
	s, isString := firstToken.Value.(string)
	if ok := tokenizer.IsKeywordConst(s); isString && ok {
		return true
	}
	if firstToken.Type == tokenizer.IDENTIFIER {
		return true
	}
	if firstToken.Value == '(' {
		return true
	}
	r, isRune := p.LookNext().Value.(rune)
	if isRune && tokenizer.IsSymbolUnaryOp(r) {
		return true
	}
	return false
}
func (p* Parser) parseTerm() tokenizer.Token {
	var toBuild []tokenizer.Token
	var isString, isRune bool
	var r rune
	var s string

	// int const, string const
	if p.LookNext().Type == tokenizer.INTCONST || p.LookNext().Type == tokenizer.STRINGCONST {
		p.ConsumeNext(&toBuild)
		goto end
	}

	// keyword const
	s, isString = p.LookNext().Value.(string)
	if ok := tokenizer.IsKeywordConst(s); isString && ok {
		p.ConsumeNext(&toBuild)
		goto end
	}

	// varName, array entry or subroutuine call
	if p.LookNext().Type == tokenizer.IDENTIFIER {

		switch p.LookNextOf(1).Value {
		case '[':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "["
			exp := p.parseExpression()
			toBuild = append(toBuild, exp)
			p.ConsumeNext(&toBuild)					// "]"

		case '(':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "("
			expList := p.parseExpressionList()
			toBuild = append(toBuild, expList)
			p.ConsumeNext(&toBuild)					// ")"

		case '.':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "."
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "("
			expList := p.parseExpressionList()
			toBuild = append(toBuild, expList)
			p.ConsumeNext(&toBuild)					// ")"

		default:
			p.ConsumeNext(&toBuild)
		}
		goto end
	}

	// expression inside parens
	if p.LookNext().Value == '(' {
		p.ConsumeNext(&toBuild)					// "("
		exp := p.parseExpression()
		toBuild = append(toBuild, exp)
		p.ConsumeNext(&toBuild)					// ")"
		goto end
	}

	// unaryOp term
	r, isRune = p.LookNext().Value.(rune)
	if isRune && tokenizer.IsSymbolUnaryOp(r) {
		p.ConsumeNext(&toBuild)
		term := p.parseTerm()
		toBuild = append(toBuild, term)
		goto end
	}

	// if here, is not a term

	end:
		return tokenizer.Token{
			Type: TERM,
			Children: &toBuild,
		}
}





func (p* Parser) parseExpressionList() tokenizer.Token {
	var toBuild []tokenizer.Token

	if p.nextIsExpression() {
		exp := p.parseExpression()
		toBuild = append(toBuild, exp)
	}

	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		exp := p.parseExpression()
		toBuild = append(toBuild, exp)
	}

	return tokenizer.Token{
		Type: EXPRESSION_LIST,
		Children: &toBuild,
	}
}