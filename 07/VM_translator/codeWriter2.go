package VM_translator

import (
	"fmt"
	"strconv"
)

func (c *CodeWriter) WritePushPop(pushPop, arg1, arg2, staticPrefix string) string {
	if pushPop == PUSH {
		return writePush(arg1, arg2, staticPrefix)
	}
	if pushPop == POP {
		return writePop(arg1, arg2, staticPrefix)
	}
	panic(fmt.Sprintf("%s %s %s :not a push/pop command", pushPop, arg1, arg2))
}

func writePush(memSegment, i, staticPrefix string) string {
	switch memSegment {
	case CONSTANT:
		return writePushConstant(i)
	case LOCAL:
		fallthrough
	case ARGS:
		fallthrough
	case THIS:
		fallthrough
	case THAT:
		return writePushIndirectSegment(memSegment, i)
	case POINTER:
		fallthrough
	case TEMP:
		return writePushDirectSegment(memSegment, i)
	case STATIC:
		return writePushStatic(i, staticPrefix)
	default:
		panic(memSegment + ": not a valid segment")
	}
}
func writePushConstant(i string) string {
	return fmt.Sprintf(
		`@%s
D=A
@%s
A=M
M=D
@%s
M=M+1
`, i, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writePushStatic(i, prefix string) string {
	label := fmt.Sprintf("%s.%s", prefix, i)
	return fmt.Sprintf(
		`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, label, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writePushIndirectSegment(memSegment string, i string) string {
	inAssembly := vmLabelsToAssemblyLabel[memSegment]
	return fmt.Sprintf(
		`@%s
D=A
@%s
A=M
A=A+D
D=M
@%s
A=M
M=D
@%s
M=M+1
`, i, inAssembly, STACK_ASSEMBLY, STACK_ASSEMBLY)
}
func writePushDirectSegment(segment string, i string) string {
	n, _ := strconv.Atoi(i)
	addr := strconv.Itoa(vmLabelsToAssemblyAddress[segment] + n)
	return fmt.Sprintf(
		`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, addr, STACK_ASSEMBLY, STACK_ASSEMBLY)
}


func writePop(memSegment, i, staticPrefix string) string {
	switch memSegment {
	case CONSTANT:
		panic(memSegment + ": not a valid segment for POP")
	case LOCAL:
		fallthrough
	case ARGS:
		fallthrough
	case THIS:
		fallthrough
	case THAT:
		return writePopIndirectSegment(memSegment, i)
	case POINTER:
		fallthrough
	case TEMP:
		return writePopDirectSegment(memSegment, i)
	case STATIC:
		return writePopStatic(i, staticPrefix)
	default:
		panic(memSegment + ": not a valid segment")
	}
}
func writePopStatic (i, prefix string) string {
	label := fmt.Sprintf("%s.%s", prefix, i)
	return fmt.Sprintf(
		`@%s
M=M-1
A=M
D=M
@%s
M=D
`, STACK_ASSEMBLY, label)
}
func writePopIndirectSegment(segment, i string) string {
	memAssembly := vmLabelsToAssemblyLabel[segment]
	return fmt.Sprintf(
		`@%s
D=A
@%s
A=M
D=A+D
@R13
M=D
@%s
M=M-1
A=M
D=M
@R13
A=M
M=D
`, i, memAssembly, STACK_ASSEMBLY)
}
func writePopDirectSegment(segment, i string) string {
	n, _ := strconv.Atoi(i)
	addr := strconv.Itoa(vmLabelsToAssemblyAddress[segment] + n)
	return fmt.Sprintf(
		`@%s
M=M-1
A=M
D=M
@%s
M=D
`, STACK_ASSEMBLY, addr)
}
