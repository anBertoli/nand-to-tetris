package compiler

const (
	CLASS          TokenType = "class"
	CLASSVARDEC    TokenType = "classVarDec"
	PARAMLIST      TokenType = "parameterList"
	SUBROUTINEDEC  TokenType = "subroutineDec"
	SUBROUTINEBODY TokenType = "subroutineBody"
	VARDEC         TokenType = "varDec"
)

type Parser struct {
	tokens []Token
}
func NewParser(tokens []Token) Parser {
	return Parser{tokens}
}

func (p* Parser) ConsumeNext(dest *[]Token) Token {
	next := p.tokens[0]
	p.tokens = p.tokens[1:]
	if dest != nil {
		*dest = append(*dest, next)
	}
	return next
}
func (p* Parser) LookNext() Token {
	return p.tokens[0]
}
func (p* Parser) LookNextOf(n int) Token {
	return p.tokens[n]
}





func (p* Parser) nextIsClass() bool {
	var firstToken = p.tokens[0]
	if firstToken.Type == KEYWORD && firstToken.Value == "class" {
		return true
	}
	return false
}
func (p* Parser) CompileClass() (Token, error) {
	var toBuild []Token

	p.ConsumeNext(&toBuild)		// "class"
	p.ConsumeNext(&toBuild)		// identifier
	p.ConsumeNext(&toBuild)		// "{"

	// compile class variable declarations, if present
	for p.nextIsClassVarDec() {
		classVarDec, _:= p.compileClassVarDec()
		toBuild = append(toBuild, classVarDec)
	}

	// compile class routines declarations, if present
	for p.nextIsSubroutineDec() {
		subrutDec, _:= p.compileSubroutineDec()
		toBuild = append(toBuild, subrutDec)
	}

	p.ConsumeNext(&toBuild)		// "}"

	return Token{
		Type:     CLASS,
		Children: &toBuild,
	}, nil
}








func (p* Parser) nextIsClassVarDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && (firstToken.Value == "static" || firstToken.Value == "field") {
		return true
	}
	return false
}
func (p* Parser) compileClassVarDec() (Token, error) {
	var toBuild []Token

	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	// compile additional var declarations
	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild)		// ";"

	return Token{
		Type:     CLASSVARDEC,
		Children: &toBuild,
	}, nil
}







func (p* Parser) nextIsSubroutineDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && (firstToken.Value == "constructor" || firstToken.Value == "function" || firstToken.Value == "method") {
		return true
	}
	return false
}
func (p* Parser) compileSubroutineDec() (Token, error) {
	var toBuild []Token

	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	// compile parameter list
	paramList := p.compileParameterList()
	toBuild = append(toBuild, paramList)

	p.ConsumeNext(&toBuild)

	// compile subroutine body
	subRoutineBody, _ := p.compileSubroutineBody()
	toBuild = append(toBuild, subRoutineBody)

	return Token{
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
func (p* Parser) compileSubroutineBody() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)		// "{"

	// compile var declarations, if present
	for p.nextIsVarDec() {
		varDec, _ := p.compileVarDec()
		toBuild = append(toBuild, varDec)
	}

	// compile statements
	stats, _ := p.compileStatements()
	toBuild = append(toBuild, stats)

	p.ConsumeNext(&toBuild)		// "}"

	return Token{
		Type: SUBROUTINEBODY,
		Children: &toBuild,
	}, nil
}






func (p* Parser) compileParameterList() Token {
	var toBuild []Token

	if p.LookNext().Type == KEYWORD || p.LookNext().Type == IDENTIFIER {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	for p.LookNext().Value == ','{
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	return Token{
		Type:     PARAMLIST,
		Children: &toBuild,
	}
}




func (p* Parser) nextIsVarDec() bool {
	var firstToken = p.LookNext()
	if firstToken.Type == KEYWORD && firstToken.Value == "var" {
		return true
	}
	return false
}
func (p* Parser) compileVarDec() (Token, error) {
	var toBuild []Token
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)
	p.ConsumeNext(&toBuild)

	for p.LookNext().Value == ',' {
		p.ConsumeNext(&toBuild)
		p.ConsumeNext(&toBuild)
	}

	p.ConsumeNext(&toBuild)

	return Token{
		Type:     VARDEC,
		Children: &toBuild,
	}, nil
}