package parser

import (
	"../tokenizer"
)

const (
	CLASS          		tokenizer.TokenType = "class"
	CLASSVARDEC    		tokenizer.TokenType = "classVarDec"
	PARAMLIST      		tokenizer.TokenType = "parameterList"
	SUBROUTINEDEC  		tokenizer.TokenType = "subroutineDec"
	SUBROUTINEBODY 		tokenizer.TokenType = "subroutineBody"
	VARDEC         		tokenizer.TokenType = "varDec"

	EXPRESSION_LIST 	tokenizer.TokenType = "expressionList"
	EXPRESSION 			tokenizer.TokenType = "expression"
	TERM 				tokenizer.TokenType = "term"

	STATS 				tokenizer.TokenType = "statements"
	LETSTAT 			tokenizer.TokenType = "letStatement"
	IFSTAT 				tokenizer.TokenType = "ifStatement"
	WHILESTAT 			tokenizer.TokenType = "whileStatement"
	DOSTAT 				tokenizer.TokenType = "doStatement"
	RETURNSTAT 			tokenizer.TokenType = "returnStatement"
)

type Parser struct {
	tokens []tokenizer.Token
}
func NewParser(tokens []tokenizer.Token) Parser {
	return Parser{tokens}
}

func (p* Parser) ConsumeNext(dest *[]tokenizer.Token) tokenizer.Token {
	next := p.tokens[0]
	p.tokens = p.tokens[1:]
	if dest != nil {
		*dest = append(*dest, next)
	}
	return next
}
func (p* Parser) LookNext() tokenizer.Token {
	return p.tokens[0]
}
func (p* Parser) LookNextOf(n int) tokenizer.Token {
	return p.tokens[n]
}





func (p* Parser) nextIsClass() bool {
	var firstToken = p.tokens[0]
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "class" {
		return true
	}
	return false
}
func (p* Parser) ParseClass() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token

	p.ConsumeNext(&toBuild)		// "class"
	p.ConsumeNext(&toBuild)		// identifier
	p.ConsumeNext(&toBuild)		// "{"

	// compile class variable declarations, if present
	for p.nextIsClassVarDec() {
		classVarDec, _:= p.parseClassVarDec()
		toBuild = append(toBuild, classVarDec)
	}

	// compile class routines declarations, if present
	for p.nextIsSubroutineDec() {
		subrutDec, _:= p.parseSubroutineDec()
		toBuild = append(toBuild, subrutDec)
	}

	p.ConsumeNext(&toBuild)		// "}"

	return tokenizer.Token{
		Type:     CLASS,
		Children: &toBuild,
	}, nil
}








func (p* Parser) nextIsClassVarDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && (firstToken.Value == "static" || firstToken.Value == "field") {
		return true
	}
	return false
}
func (p* Parser) parseClassVarDec() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token

	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	// compile additional var declarations
	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild)		// ";"

	return tokenizer.Token{
		Type:     CLASSVARDEC,
		Children: &toBuild,
	}, nil
}







func (p* Parser) nextIsSubroutineDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && (firstToken.Value == "constructor" || firstToken.Value == "function" || firstToken.Value == "method") {
		return true
	}
	return false
}
func (p* Parser) parseSubroutineDec() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token

	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	// compile parameter list
	paramList := p.parseParameterList()
	toBuild = append(toBuild, paramList)

	p.ConsumeNext(&toBuild)

	// compile subroutine body
	subRoutineBody, _ := p.parseSubroutineBody()
	toBuild = append(toBuild, subRoutineBody)

	return tokenizer.Token{
		Type:     SUBROUTINEDEC,
		Children: &toBuild,
	}, nil
}







func (p* Parser) nextIsSubroutineBody() bool {
	var first = p.tokens[0]
	if first.Value == '{' {
		return true
	}
	return false
}
func (p* Parser) parseSubroutineBody() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)		// "{"

	// compile var declarations, if present
	for p.nextIsVarDec() {
		varDec, _ := p.parseVarDec()
		toBuild = append(toBuild, varDec)
	}

	// compile statements
	stats, _ := p.parseStatements()
	toBuild = append(toBuild, stats)

	p.ConsumeNext(&toBuild)		// "}"

	return tokenizer.Token{
		Type: SUBROUTINEBODY,
		Children: &toBuild,
	}, nil
}






func (p* Parser) parseParameterList() tokenizer.Token {
	var toBuild []tokenizer.Token

	if p.LookNext().Type == tokenizer.KEYWORD || p.LookNext().Type == tokenizer.IDENTIFIER {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	return tokenizer.Token{
		Type:     PARAMLIST,
		Children: &toBuild,
	}
}




func (p* Parser) nextIsVarDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == tokenizer.KEYWORD && firstToken.Value == "var" {
		return true
	}
	return false
}
func (p* Parser) parseVarDec() (tokenizer.Token, error) {
	var toBuild []tokenizer.Token
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild)

	return tokenizer.Token{
		Type:     VARDEC,
		Children: &toBuild,
	}, nil
}