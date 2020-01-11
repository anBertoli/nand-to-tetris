package parser

import (
	"../tokenizer"
)



// STATEMENTS
func (p* Parser) nextAreStatements() bool {
	return p.nextIsLetStatement() || p.nextIsIfStatement() || p.nextIsWhileStatement() || p.nextIsDoStatement() || p.nextIsReturnStatement()
}
func (p* Parser) parseStatements() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token

	for {
		if p.nextIsLetStatement() {
			letStat, _ := p.parseLetStatement()
			toBuild = append(toBuild, letStat)
			continue
		}
		if p.nextIsIfStatement() {
			ifStat, _ := p.parseIfStatement()
			toBuild = append(toBuild, ifStat)
			continue
		}
		if p.nextIsWhileStatement() {
			whileStat, _ := p.parseWhileStatement()
			toBuild = append(toBuild, whileStat)
			continue
		}
		if p.nextIsDoStatement() {
			doStat, _ := p.parseDoStatement()
			toBuild = append(toBuild, doStat)
			continue
		}
		if p.nextIsReturnStatement() {
			retStat, _ := p.parseReturnStatement()
			toBuild = append(toBuild, retStat)
			continue
		}
		break
	}

	return tokenizer.Token{
		Type:     STATS,
		Children: &toBuild,
	}, nil
}









// LET statement
func (p* Parser) nextIsLetStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "let" {
		return true
	}
	return false
}
func (p* Parser) parseLetStatement() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	if p.LookNext().Value == '[' {
		p.ConsumeNext(&toBuild)
		expr := p.parseExpression()
		toBuild = append(toBuild, expr)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild) 		// "="

	expr := p.parseExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 		// ";"

	return tokenizer.Token{
		Type:     LETSTAT,
		Children: &toBuild,
	}, nil
}









// IF statement
func (p* Parser) nextIsIfStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "if" {
		return true
	}
	return false
}
func (p* Parser) parseIfStatement() (tokenizer.Token, error)  {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)

	p.ConsumeNext(&toBuild)				// "("
	expr := p.parseExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 			// ")"

	p.ConsumeNext(&toBuild) 			// "{"
	stats, _ := p.parseStatements()
	toBuild = append(toBuild, stats)
	p.ConsumeNext(&toBuild) 			// "}"

	if p.LookNext().Value == "else" {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)			// "{"
		stats, _ := p.parseStatements()
		toBuild = append(toBuild, stats)
		p.ConsumeNext(&toBuild)			// "}"
	}

	return tokenizer.Token{
		Type: IFSTAT,
		Children: &toBuild,
	}, nil
}









// WHILE statement
func (p* Parser) nextIsWhileStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "while" {
		return true
	}
	return false
}
func (p* Parser) parseWhileStatement() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)

	p.ConsumeNext(&toBuild)				// "("
	expr := p.parseExpression()
	toBuild = append(toBuild, expr)
	p.ConsumeNext(&toBuild) 			// ")"

	p.ConsumeNext(&toBuild) 			// "{"
	stats, _ := p.parseStatements()
	toBuild = append(toBuild, stats)
	p.ConsumeNext(&toBuild) 			// "}"

	return tokenizer.Token{
		Type: WHILESTAT,
		Children: &toBuild,
	}, nil
}








// DO statement
func (p* Parser) nextIsDoStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "do" {
		return true
	}
	return false
}
func (p* Parser) parseDoStatement() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)				// "do"
	p.ConsumeNext(&toBuild)

	switch p.LookNext().Value {
	case '(':
		p.ConsumeNext(&toBuild)					// "("
		expList := p.parseExpressionList()
		toBuild = append(toBuild, expList)
		p.ConsumeNext(&toBuild)					// ")"

	case '.':
		p.ConsumeNext(&toBuild)					// "."
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)					// "("
		expList := p.parseExpressionList()
		toBuild = append(toBuild, expList)
		p.ConsumeNext(&toBuild)					// ")"
	}

	p.ConsumeNext(&toBuild) 			// ";"

	return tokenizer.Token{
		Type: DOSTAT,
		Children: &toBuild,
	}, nil
}








// RETURN statement
func (p* Parser) nextIsReturnStatement() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "return" {
		return true
	}
	return false
}
func (p* Parser) parseReturnStatement() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)				// "return"
	if p.LookNext().Value != ';' {
		expr := p.parseExpression()
		toBuild = append(toBuild, expr)
	}
	p.ConsumeNext(&toBuild) 			// ";"

	return tokenizer.Token{
		Type: RETURNSTAT,
		Children: &toBuild,
	}, nil
}
