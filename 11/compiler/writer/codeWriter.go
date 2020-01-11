package writer

import (
	"../parser"
	"../tokenizer"
	"fmt"
	"strings"
)

type Writer struct {
	class          string
	subRoutine	   string
	symbolTable    SymbolTable
	counter		   int
}
func NewWriter() Writer {
	return Writer{}
}

func (w* Writer) cleanSubroutine() {
	w.subRoutine = ""
}


func (w* Writer) CompileClass(classToken tokenizer.Token) string {
	var classTokens = classToken.GetChildren()
	var compiled string

	w.class = classTokens[1].Raw

	// compile fields and statics
	for _, token := range classTokens {
		if token.Type == parser.CLASSVARDEC {
			classVarDecTokens := token.GetChildren()

			if classVarDecTokens[0].Value.(string) == "field" {
				w.defineField(classVarDecTokens)
			}
			if classVarDecTokens[0].Value.(string) == "static" {
				w.defineStatic(classVarDecTokens)
			}
		}
	}

	// compile subroutine declaration
	for _, token := range classTokens {
		if token.Type == parser.SUBROUTINEDEC {
			w.symbolTable.cleanSubroutineScope()
			subRoutineTokens := token.GetChildren()

			switch subRoutineTokens[0].Value.(string) {
			case "constructor":
				compiled += w.writeConstructor(subRoutineTokens)
			case "function":
				compiled += w.writeFunction(subRoutineTokens)
			case "method":
				compiled += w.writeMethod(subRoutineTokens)
			}
		}
	}

	pretty := strings.ReplaceAll(compiled, "\t", "")
	pretty = strings.ReplaceAll(pretty, "\n\n", "\n")
	pretty = strings.ReplaceAll(pretty, "\n\n", "\n")
	return pretty
}




func (w* Writer) defineStatic(classVarDecTokens []tokenizer.Token) {
	typ := symbolType(classVarDecTokens[1].Type)
	for _, token := range classVarDecTokens[2:] {
		if token.Type == tokenizer.IDENTIFIER {
			name := token.Value.(string)
			w.symbolTable.define(name, typ, STATIC, CLASS_SCOPE)
		}
	}
}

func (w* Writer) defineField(classVarDecTokens []tokenizer.Token) {
	typ := symbolType(classVarDecTokens[1].Raw)
	for _, token := range classVarDecTokens[2:] {
		if token.Type == tokenizer.IDENTIFIER {
			name := token.Value.(string)
			w.symbolTable.define(name, typ, FIELD, CLASS_SCOPE)
		}
	}
}

func (w* Writer) defineParameterList (paramList []tokenizer.Token) {
	for i := 0; i < len(paramList); i += 3 {
		typ := symbolType(paramList[i].Value.(string))
		name := paramList[i+1].Value.(string)
		w.symbolTable.define(name, typ, ARGUMENT, SUBROUTINE_SCOPE)
	}
}

func (w* Writer) defineLocalList (varDecs []tokenizer.Token) {
	for _, varDec := range varDecs {
		subVarDec := varDec.GetChildren()
		subVarDec = subVarDec[1:len(subVarDec)-1]
		varType := symbolType(subVarDec[0].Raw)

		for _, t := range subVarDec[1:] {
			if t.Type == tokenizer.IDENTIFIER {
				w.symbolTable.define(t.Raw, varType, LOCAL, SUBROUTINE_SCOPE)
			}
		}
	}
}




func (w* Writer) writeConstructor(constructorTokens []tokenizer.Token) string {
	var className = constructorTokens[1].Raw
	var functionName = constructorTokens[2].Raw

	var params = constructorTokens[4].GetChildren()
	var bodyStats = getTokenOfType(constructorTokens[6].GetChildren(), parser.STATS).GetChildren()
	var varDecs = getTokensOfType(constructorTokens[6].GetChildren(), parser.VARDEC)

	// register parameters and locals
	w.defineParameterList(params)
	w.defineLocalList(varDecs)

	// alloc memory & compile body
	numFields := w.symbolTable.lastOfKind(FIELD) + 1	// bytes needed
	numLocals := w.symbolTable.lastOfKind(LOCAL) + 1
	compiledStats := w.writeStatements(bodyStats)

	return fmt.Sprintf(`// constructor
		function %s.%s %v
		push constant %v
		call Memory.alloc 1
		pop pointer 0
		%s
		return
	`, className, functionName, numLocals, numFields, compiledStats)
}


func (w* Writer) writeFunction(constructorTokens []tokenizer.Token) string {
	var className = w.class
	var functionName = constructorTokens[2].Raw

	var params = constructorTokens[4].GetChildren()
	var bodyStats = getTokenOfType(constructorTokens[6].GetChildren(), parser.STATS).GetChildren()
	var varDecs = getTokensOfType(constructorTokens[6].GetChildren(), parser.VARDEC)

	// register parameters and locals
	w.defineParameterList(params)
	w.defineLocalList(varDecs)

	// compile body
	numLocals := w.symbolTable.lastOfKind(LOCAL) + 1
	compiledStats := w.writeStatements(bodyStats)

	return fmt.Sprintf(`// function
		function %s.%s %v
		%s
		return
	`, className, functionName, numLocals, compiledStats)
}


func (w* Writer) writeMethod(constructorTokens []tokenizer.Token) string {
	var className = w.class
	var functionName = constructorTokens[2].Raw

	var params = constructorTokens[4].GetChildren()
	var bodyStats = getTokenOfType(constructorTokens[6].GetChildren(), parser.STATS).GetChildren()
	var varDecs = getTokensOfType(constructorTokens[6].GetChildren(), parser.VARDEC)

	// register parameters and locals
	w.defineParameterList(params)
	w.defineLocalList(varDecs)

	// compile body
	numLocals := w.symbolTable.lastOfKind(LOCAL) + 1
	compiledStats := w.writeStatements(bodyStats)

	return fmt.Sprintf(`// method
		function %s.%s %v
		push argument 0
		pop pointer 0
		%s
		return
	`, className, functionName, numLocals, compiledStats)
}





func getTokenOfType(tokens []tokenizer.Token, tokenType tokenizer.TokenType) tokenizer.Token {
	for _, t := range tokens {
		if t.Type == tokenType {
			return t
		}
	}
	panic("not found token" + tokenType)
}
func getTokensOfType(tokens []tokenizer.Token, tokenType tokenizer.TokenType) []tokenizer.Token {
	var tt []tokenizer.Token
	for _, t := range tokens {
		if t.Type == tokenType {
			tt = append(tt, t)
		}
	}
	return tt
}