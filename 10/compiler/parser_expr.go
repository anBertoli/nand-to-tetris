package compiler

const (
	EXPRESSION_LIST 	TokenType = "expressionList"
	EXPRESSION 			TokenType = "expression"
	TERM 				TokenType = "term"
)


func (p* Parser) nextIsExpression () bool {
	return p.nextIsTerm()
}
func (p* Parser) compileExpression () Token {
	var toBuild []Token

	term := p.compileTerm()
	toBuild = append(toBuild, term)

	r, isRune := p.LookNext().Value.(rune)
	if isRune && isSymbolOp(r) {
		p.ConsumeNext(&toBuild)
		term := p.compileTerm()
		toBuild = append(toBuild, term)
	}

	return Token{
		Type: EXPRESSION,
		Children: &toBuild,
	}
}



func (p* Parser) nextIsTerm() bool {
	var firstToken = p.LookNext()

	if firstToken.Type == INTCONST || firstToken.Type == STRINGCONST {
		return true
	}

	s, isString := firstToken.Value.(string)
	if _, ok := keywordsConst[s]; isString && ok {
		return true
	}

	if firstToken.Type == IDENTIFIER {
		return true
	}

	if firstToken.Value == '(' {
		return true
	}

	r, isRune := p.LookNext().Value.(rune)
	if isRune && isSymbolUnaryOp(r) {
		return true
	}

	return false
}
func (p* Parser) compileTerm () Token {
	var toBuild []Token
	var isString, isRune bool
	var r rune
	var s string

	// int const, string const
	if p.LookNext().Type == INTCONST || p.LookNext().Type == STRINGCONST {
		p.ConsumeNext(&toBuild)
		goto end
	}

	// keyword const
	s, isString = p.LookNext().Value.(string)
	if _, ok := keywordsConst[s]; isString && ok {
		p.ConsumeNext(&toBuild)
		goto end
	}

	// varName, array entry or subrotuine call
	if p.LookNext().Type == IDENTIFIER {

		switch p.LookNextOf(1).Value {
		case '[':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "["
			exp := p.compileExpression()
			toBuild = append(toBuild, exp)
			p.ConsumeNext(&toBuild)					// "]"

		case '(':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "("
			expList := p.compileExpressionList()
			toBuild = append(toBuild, expList)
			p.ConsumeNext(&toBuild)					// ")"

		case '.':
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "."
			p.ConsumeNext(&toBuild)
			p.ConsumeNext(&toBuild)					// "("
			expList := p.compileExpressionList()
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
		exp := p.compileExpression()
		toBuild = append(toBuild, exp)
		p.ConsumeNext(&toBuild)					// ")"
		goto end
	}

	// unaryOp term
	r, isRune = p.LookNext().Value.(rune)
	if isRune && isSymbolUnaryOp(r) {
		p.ConsumeNext(&toBuild)
		term := p.compileTerm()
		toBuild = append(toBuild, term)
		goto end
	}

	// if here, is not a term

	end:
		return Token{
			Type: TERM,
			Children: &toBuild,
		}
}





func (p* Parser) compileExpressionList () Token {
	var toBuild []Token

	if p.nextIsExpression() {
		exp := p.compileExpression()
		toBuild = append(toBuild, exp)
	}

	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		exp := p.compileExpression()
		toBuild = append(toBuild, exp)
	}

	return Token{
		Type:EXPRESSION_LIST,
		Children: &toBuild,
	}
}