package writer

import (
	"../parser"
	"../tokenizer"
	"fmt"
	"strconv"
	"strings"
)

var opMap = map[string]string{
	"+": "add",
	"-": "sub",
	"=": "eq",
	">": "gt",
	"<": "lt",
	"&": "and",
	"|": "or",
	"*": "call Math.multiply 2",
	"/": "call Math.divide 2",
}
var unaryOpMap = map[string]string{
	"-": "neg",
	"~": "not",
}
func writeOp (op string) string {
	op = strings.TrimSpace(op)
	if vmOp, ok := opMap[op]; ok {
		return vmOp
	}
	panic("OOOOOP")
}
func writeUnaryOp (op string) string {
	op = strings.TrimSpace(op)
	if vmOp, ok := unaryOpMap[op]; ok {
		return vmOp
	}
	panic("UNARY OOOOOP")
}


func (w* Writer) writeExpression(token tokenizer.Token) string {
	var children = token.GetChildren()

	switch len(children) {

	// term
	case 1:
		term := children[0]
		code := w.writeTerm(term)
		return code + "\n"

	// unaryOp term
	case 2:
		op := writeUnaryOp(children[0].Raw)
		term := children[1]
		code := w.writeTerm(term)
		return code + "\n" + op

	// term op term
	case 3:
		term1 := children[0]
		op := writeOp(children[1].Raw)
		term2 := children[2]
		code1 := w.writeTerm(term1)
		code2 := w.writeTerm(term2)
		return code1 + "\n" + code2 + "\n" + op + "\n"
	}

	panic("invalid expression")
}



func (w* Writer) writeTerm(token tokenizer.Token) string {
	var children = token.GetChildren()

	// int, string, keyword, varName
	if len(children) == 1 {
		return w.writeUnaryTerm(children[0])
	}

	// unaryOp term
	if len(children) == 2 {
		op := writeUnaryOp(children[0].Raw)
		term := children[1]
		return w.writeTerm(term) + "\n" + op
	}

	// varName[expr]
	if children[0].Type == tokenizer.IDENTIFIER && children[1].Raw == "[" {
		varName, _ := children[0].Value.(string)
		kind := w.symbolTable.kindOf(varName)

		index := w.symbolTable.indexOf(varName)
		expr := w.writeExpression(children[2])
		return fmt.Sprintf(`push %s %v
			%s
			add
			pop pointer 1
			push that 0
		`, kind, index, expr)
	}

	// (expr)
	if children[0].Raw == "(" {
		expr := w.writeExpression(children[1])
		return expr
	}

	// subroutine call
	if children[0].Type == tokenizer.IDENTIFIER && (children[1].Raw == "(" || children[3].Raw == "(" ) {
		return w.writeSubroutineCall(children)
	}

	panic("invalid term")
}







func (w* Writer) writeSubroutineCall(tokens []tokenizer.Token) string {
	var args string
	var nArgs int


	// method of actual instance (itself)
	if tokens[1].Raw == "(" {
		class := w.class
		methodName := tokens[0].Raw
		args, nArgs = w.writeExpressionList(tokens[2])

		return fmt.Sprintf(`// instance call
			push pointer 0
			%s
			call %s.%s %v
		`, args, class, methodName, nArgs+1)
	}

	// method of other instance or function of any class
	if tokens[3].Raw == "(" {

		classOrVar := tokens[0].Raw
		routineName := tokens[2].Raw
		args, nArgs = w.writeExpressionList(tokens[4])

		if w.class == "PongGame"  && tokens[2].Raw == "new" {
			fmt.Printf("%+v\n", tokens[0])
			fmt.Printf("%+v\n", tokens[2])
			fmt.Println(classOrVar)
			fmt.Printf("%+v\n", w.symbolTable)
			fmt.Println()
		}

		if w.symbolTable.exists(classOrVar) {
			// instance method
			kind := w.symbolTable.kindOf(classOrVar)
			index := w.symbolTable.indexOf(classOrVar)
			typeClass := w.symbolTable.typeOf(classOrVar)
			return fmt.Sprintf(`// instance call
				push %s %v
				%s
				call %s.%s %v
			`, kind, index, args, typeClass, routineName, nArgs+1)

		} else {
			// class function

			return fmt.Sprintf(`// function call
				%s
				call %s.%s %v
			`, args, classOrVar, routineName, nArgs)
		}

	}

	panic("not valid subroutine call")
}








func (w* Writer) writeExpressionList(expressionList tokenizer.Token) (string, int) {
	var args []tokenizer.Token
	var writtenArgs string
	var nArgs int

	for _, t := range expressionList.GetChildren() {
		if t.Type == parser.EXPRESSION {
			args = append(args, t)
		}
	}

	for _, arg := range args {
		writtenArgs += w.writeExpression(arg)
		nArgs++
	}
	return writtenArgs, nArgs
}







func (w* Writer) writeUnaryTerm(token tokenizer.Token) string {

	if token.Type == tokenizer.INTCONST {
		return "push constant " + token.Raw
	}

	if token.Type == tokenizer.STRINGCONST {
		str := token.Value.(string)
		build := ""
		for _, r := range str {
			build += "push constant " + strconv.Itoa(int(r)) + "\n"
			build += "call String.appendChar 2 \n"
		}
		return fmt.Sprintf(`push constant %v
			call String.new 1
			%s
		`, len(str), build)
	}

	if token.Type == tokenizer.KEYWORD {
		str := token.Value.(string)
		switch str {
		case "true":
			return "push constant 1 \n neg"
		case "false":
			fallthrough
		case "null":
			return "push constant 0"
		case "this":
			return "push pointer 0"
		}
	}

	if token.Type == tokenizer.IDENTIFIER {
		s, _ := token.Value.(string)
		kind := string(w.symbolTable.kindOf(s))
		index := strconv.Itoa(w.symbolTable.indexOf(s))
		return "push " + kind + " " + index + "\n"
	}

	panic(fmt.Sprintf("invalid term: %+v", token))
}





func (w* Writer) printContext (token []tokenizer.Token) {
	for _, t := range token {
		fmt.Printf("%+v\n", t)
		_printContext(t, 0)
		fmt.Println()
	}
}
func _printContext (token tokenizer.Token, spaces int) {
	padding := strings.Repeat("\t", spaces)
	if token.Value != nil {
		fmt.Printf(padding + "Val: %+v \tRaw : %+v \tType: %+v\n", token.Value, token.Raw, token.Type)
	} else {
		for _, t := range token.GetChildren() {
			_printContext(t, spaces+2)
		}
	}
}