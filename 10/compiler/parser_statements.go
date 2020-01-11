package compiler

const (
	STATS 		TokenType = "statements"
	LETSTAT 	TokenType = "letStatement"
	IFSTAT 		TokenType = "ifStatement"
	WHILESTAT 	TokenType = "whileStatement"
	DOSTAT 		TokenType = "doStatement"
	RETURNSTAT 	TokenType = "returnStatement"
)




// STATEMENTS
func (p* Parser) nextAreStatements() bool {
	return p.nextIsLetStatement() || p.nextIsIfStatement() || p.nextIsWhileStatement() || p.nextIsDoStatement() || p.nextIsReturnStatement()
}
func (p* Parser) compileStatements() (Token, error) {
	var toBuild []Token

	for {
		if p.nextIsLetStatement() {
			letStat, _ := p.compileLetStatement()
			toBuild = append(toBuild, letStat)
			continue
		}
		if p.nextIsIfStatement() {
			ifStat, _ := p.compileIfStatement()
			toBuild = append(toBuild, ifStat)
			continue
		}
		if p.nextIsWhileStatement() {
			whileStat, _ := p.compileWhileStatement()
			toBuild = append(toBuild, whileStat)
			continue
		}
		if p.nextIsDoStatement() {
			doStat, _ := p.compileDoStatement()
			toBuild = append(toBuild, doStat)
			continue
		}
		if p.nextIsReturnStatement() {
			retStat, _ := p.compileReturnStatement()
			toBuild = append(toBuild, retStat)
			continue
		}
		break
	}

	return Token{
		Type:     STATS,
		Children: &toBuild,
	}, nil
}









// LET statement
func (p* Parser) nextIsLetStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "let" {
		return true
	}
	return false
}
func (p* Parser) compileLetStatement() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	if p.LookNext().Value == '[' {
		p.ConsumeNext(&toBuild)
		expr := p.compileExpression()
		toBuild = append(toBuild, expr)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild) 		// "="

	expr := p.compileExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 		// ";"

	return Token{
		Type:     LETSTAT,
		Children: &toBuild,
	}, nil
}









// IF statement
func (p* Parser) nextIsIfStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "if" {
		return true
	}
	return false
}
func (p* Parser) compileIfStatement() (Token, error)  {
	var toBuild []Token
	p.ConsumeNext(&toBuild)

	p.ConsumeNext(&toBuild)				// "("
	expr := p.compileExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 			// ")"

	p.ConsumeNext(&toBuild) 			// "{"
	stats, _ := p.compileStatements()
	toBuild = append(toBuild, stats)
	p.ConsumeNext(&toBuild) 			// "}"

	if p.LookNext().Value == "else" {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)			// "{"
		stats, _ := p.compileStatements()
		toBuild = append(toBuild, stats)
		p.ConsumeNext(&toBuild)			// "}"
	}

	return Token{
		Type: IFSTAT,
		Children: &toBuild,
	}, nil
}









// WHILE statement
func (p* Parser) nextIsWhileStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "while" {
		return true
	}
	return false
}
func (p* Parser) compileWhileStatement() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)

	p.ConsumeNext(&toBuild)				// "("
	expr := p.compileExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 			// ")"

	p.ConsumeNext(&toBuild) 			// "{"
	stats, _ := p.compileStatements()
	toBuild = append(toBuild, stats)
	p.ConsumeNext(&toBuild) 			// "}"

	return Token{
		Type: WHILESTAT,
		Children: &toBuild,
	}, nil
}








// DO statement
func (p* Parser) nextIsDoStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "do" {
		return true
	}
	return false
}
func (p* Parser) compileDoStatement() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)				// "do"
	p.ConsumeNext(&toBuild)

	switch p.LookNext().Value {
	case '(':
		p.ConsumeNext(&toBuild)					// "("
		expList := p.compileExpressionList()
		toBuild = append(toBuild, expList)
		p.ConsumeNext(&toBuild)					// ")"

	case '.':
		p.ConsumeNext(&toBuild)					// "."
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)					// "("
		expList := p.compileExpressionList()
		toBuild = append(toBuild, expList)
		p.ConsumeNext(&toBuild)					// ")"
	}

	p.ConsumeNext(&toBuild) 			// ";"

	return Token{
		Type: DOSTAT,
		Children: &toBuild,
	}, nil
}








// RETURN statement
func (p* Parser) nextIsReturnStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "return" {
		return true
	}
	return false
}
func (p* Parser) compileReturnStatement() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)				// "return"
	if p.LookNext().Value != ';' {
		expr := p.compileExpression()
		toBuild = append(toBuild, expr)
	}
	p.ConsumeNext(&toBuild) 			// ";"

	return Token{
		Type: RETURNSTAT,
		Children: &toBuild,
	}, nil
}
