package compiler



const (
	KEYWORD     TokenType = "keyword"
	SYMBOL      TokenType = "symbol"
	IDENTIFIER  TokenType = "identifier"
	INTCONST    TokenType = "integerConstant"
	STRINGCONST TokenType = "stringConstant"
	INVALID     TokenType = "invalid"
)

type TokenType string
type TokenValue interface{}
type Token struct {
	Raw      string
	Type     TokenType
	Value    TokenValue
	Children *[]Token
}



var empty struct{}

var symbols = map[rune]struct{}{
	'{': empty, '}': empty, '(': empty, ')': empty,
	'[': empty, ']': empty, '.': empty, ',': empty,
	';': empty, '+': empty, '-': empty, '*': empty,
	'/': empty, '&': empty, '|': empty, '<': empty,
	'>': empty, '=': empty, '~': empty,
}
var symbolsOp = map[rune]struct{}{
	'+': empty, '-': empty, '*': empty,
	'/': empty, '&': empty, '|': empty, '<': empty,
	'>': empty, '=': empty,
}
var symbolsUnaryOp = map[rune]struct{}{
	'-': empty, '~': empty,
}

func isSymbol(char rune) bool {
	_, ok := symbols[char]
	return ok
}
func isSymbolOp(char rune) bool {
	_, ok := symbolsOp[char]
	return ok
}
func isSymbolUnaryOp(char rune) bool {
	_, ok := symbolsUnaryOp[char]
	return ok
}



var keywords = map[string]struct{}{
	"class": empty, "constructor": empty, "function": empty,
	"method": empty, "field": empty, "static": empty,
	"var": empty, "int": empty, "char": empty, "boolean": empty,
	"void": empty, "true": empty, "false": empty, "null": empty,
	"this": empty, "let": empty, "do": empty, "if": empty,
	"else": empty, "while": empty, "return": empty,
}
var keywordsConst = map[string]struct{}{
	"this": empty, "true": empty, "false": empty, "null": empty,
}

