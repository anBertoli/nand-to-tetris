package VM_translator

import (
	"fmt"
	"strconv"
)

type CodeWriter struct {
	labelCounter int
}

func NewCodeWriter() *CodeWriter {
	return &CodeWriter{labelCounter: 0}
}

func (c *CodeWriter) WriteArithmetic(op string) string {
	switch op {
	case ADD:
		return writeAdd()
	case SUB:
		return writeSub()
	case NEG:
		return writeNeg()
	case EQ:
		c.labelCounter++
		return writeComp(EQ, c.labelCounter)
	case GT:
		c.labelCounter++
		return writeComp(GT, c.labelCounter)
	case LT:
		c.labelCounter++
		return writeComp(LT, c.labelCounter)
	case AND:
		return writeAnd()
	case OR:
		return writeOr()
	case NOT:
		return writeNot()
	default:
		panic(op + ": not a valid operation")
	}
}

func writeAdd() string {
	return fmt.Sprintf(
		`@%s
AM=M-1
D=M
A=A-1
M=D+M
D=A+1
@%s
M=D
`, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writeSub() string {
	return fmt.Sprintf(
		`@%s
AM=M-1
D=M
A=A-1
M=M-D
D=A+1
@%s
M=D
`, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writeNeg() string {
	return fmt.Sprintf(
		`@%s
A=M-1
M=-M
`, STACK_ASSEMBLY)
}

func writeComp(op string, counter int) string {
	var compare string
	switch op {
	case EQ:
		compare = "JEQ"
	case GT:
		compare = "JGT"
	case LT:
		compare = "JLT"
	}
	trueLabel, endLabel := "TRUE"+strconv.Itoa(counter), "END"+strconv.Itoa(counter)

	return fmt.Sprintf(
		`@%s
AM=M-1
D=M
@%s
AM=M-1
D=M-D
@%s
D;%s
D=0
@%s
0;JMP
(%s)
D=-1
(%s)
@%s
A=M
M=D
@%s
M=M+1
`, STACK_ASSEMBLY, STACK_ASSEMBLY, trueLabel, compare, endLabel, trueLabel, endLabel, STACK_ASSEMBLY, STACK_ASSEMBLY)
}

func writeAnd() string {
	return fmt.Sprintf(
		`@%s
AM=M-1
D=M
A=A-1
M=D&M
D=A+1
@%s
M=D
`, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writeNot() string {
	return fmt.Sprintf(
		`@%s
A=M-1
M=!M
`, STACK_ASSEMBLY)
}
func writeOr() string {
	return fmt.Sprintf(
		`@%s
AM=M-1
D=M
A=A-1
M=D|M
D=A+1
@%s
M=D
`, STACK_ASSEMBLY, STACK_ASSEMBLY)
}

func WriteEndLoop() string {
	return `
(ENDLOOP)
@ENDLOOP
0;JMP
`
}
