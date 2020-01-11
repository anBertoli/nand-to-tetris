package writer

import (
	"../tokenizer"
	"../parser"
	"fmt"
)

func (w* Writer) writeLetStatement(letStatTokens []tokenizer.Token) string {
	var isArray = letStatTokens[2].Value.(rune) == '['
	if isArray {
		return w.writeLetArrayStatement(letStatTokens)
	}

	var identifier = letStatTokens[1].Value.(string)
	if !w.symbolTable.exists(identifier) {
		panic(identifier + " not recorded")
	}

	rightHandExpr := letStatTokens[3]
	compiledExpr := w.writeExpression(rightHandExpr)

	kind := w.symbolTable.kindOf(identifier)
	index := w.symbolTable.indexOf(identifier)

	switch kind {
	case FIELD:
		return fmt.Sprintf(`// let statement
			%s
			pop this %v
		`, compiledExpr, index)
	case STATIC:
		return fmt.Sprintf(`// let statement
			%s
			pop static %v
		`, compiledExpr, index)
	case LOCAL:
		return fmt.Sprintf(`// let statement
			%s
			pop local %v
		`, compiledExpr, index)
	case ARGUMENT:
		return fmt.Sprintf(`// let statement
			%s
			pop argument %v
		`, compiledExpr, index)
	default:
		panic("writeLetStatement")
	}

}


func (w* Writer) writeLetArrayStatement(letStatTokens []tokenizer.Token) string {
	var identifier = letStatTokens[1].Value.(string)
	var compiled string
	if !w.symbolTable.exists(identifier) {
		panic(identifier + " not recorded")
	}

	kind := w.symbolTable.kindOf(identifier)
	index := w.symbolTable.indexOf(identifier)

	leftHandExpr := letStatTokens[3]
	compiledLeftHandExpr := w.writeExpression(leftHandExpr)
	rightHandExpr := letStatTokens[6]
	compiledRightHandExpr := w.writeExpression(rightHandExpr)


	switch kind {
	case FIELD:
		compiled += fmt.Sprintf("// let statement \n push this %v", index)
	case STATIC:
		compiled += fmt.Sprintf("// let statement \n push static %v", index)
	case LOCAL:
		compiled += fmt.Sprintf("// let statement \n push local %v", index)
	case ARGUMENT:
		compiled += fmt.Sprintf("// let statement \n push argument %v", index)
	default:
		panic("writeLetArrayStatement")
	}

	return fmt.Sprintf(`%s
		%s
		add
		%s
		pop temp 0
		pop pointer 1
		push temp 0
		pop that 0
	`, compiled, compiledLeftHandExpr, compiledRightHandExpr)
}



func (w* Writer) writeDoStatement(doStatTokens []tokenizer.Token) string {
	doStatTokens = doStatTokens[1:]

	// discard fake returned value
	return "// do statement\n" + w.writeSubroutineCall(doStatTokens) + "pop temp 0 \n"
}


func (w* Writer) writeReturnStatement(returnStatTokens []tokenizer.Token) string {
	returnStatTokens = returnStatTokens[1:]
	if returnStatTokens[0].Raw == ";" {
		return "// return statement \npush constant 0 \n"
	}
	return "// return statement \n" + w.writeExpression(returnStatTokens[0])
}



func (w* Writer) writeIfStatement(ifStatTokens []tokenizer.Token) string {
	var elseIndex = -1
	for i, t := range ifStatTokens {
		if t.Raw == "else" {
			elseIndex = i
		}
	}

	var cond = ifStatTokens[2]
	var trueStats = ifStatTokens[5].GetChildren()

	compiledCond := w.writeExpression(cond)
	compiledTrueStats := w.writeStatements(trueStats)

	w.counter++
	labelElse := fmt.Sprintf("IF_TRUE_%v", w.counter)
	w.counter++
	labelEnd := fmt.Sprintf("IF_FALSE_%v", w.counter)

	if elseIndex == -1 {
		// doesn't have else
		return fmt.Sprintf(`// if statement
			%s
			not
			if-goto %s
			%s
			label %s
		`, compiledCond, labelEnd, compiledTrueStats, labelEnd)
	} else {
		// does have else
		var elseStats = ifStatTokens[9]
		compiledElseStats := w.writeStatements(elseStats.GetChildren())
		return fmt.Sprintf(`// if statement
			%s
			not
			if-goto %s
			%s
			goto %s
			label %s
			%s
			label %s
		`, compiledCond, labelElse, compiledTrueStats, labelEnd, labelElse, compiledElseStats, labelEnd)
	}
}


func (w* Writer) writeWhileStatement(whileStatTokens []tokenizer.Token) string {
	var cond = whileStatTokens[2]
	var stats = whileStatTokens[5].GetChildren()

	compiledCond := w.writeExpression(cond)
	compiledStats := w.writeStatements(stats)

	w.counter++
	labelStart := fmt.Sprintf("WHILE_START_%v", w.counter)
	w.counter++
	labelEnd := fmt.Sprintf("WHILE_END_%v", w.counter)

	return fmt.Sprintf(`// while statement
		label %s
		%s
		not
		if-goto %s
		%s
		goto %s
		label %s
	`, labelStart, compiledCond, labelEnd, compiledStats, labelStart, labelEnd)
}



func (w* Writer) writeStatements(statsTokens []tokenizer.Token) string {
	var compiled string
	for _, stat := range statsTokens {
		children := stat.GetChildren()

		switch stat.Type {
		case parser.LETSTAT:
			compiled += w.writeLetStatement(children)
		case parser.DOSTAT:
			compiled += w.writeDoStatement(children)
		case parser.RETURNSTAT:
			compiled += w.writeReturnStatement(children)
		case parser.IFSTAT:
			compiled += w.writeIfStatement(children)
		case parser.WHILESTAT:
			compiled += w.writeWhileStatement(children)
		default:
			panic("invalid statement: " + stat.Type)
		}
	}
	return compiled
}