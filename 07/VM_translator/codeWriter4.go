package VM_translator

import (
	"fmt"
	"strconv"
)

func (c *CodeWriter) WriteFunction(name string, numLocals int) string {
	cleanLocal := fmt.Sprintf("(%s)\n", name)
	if numLocals > 0 {
		for i := 1; i <= numLocals; i++ {
			cleanLocal += fmt.Sprintf(`@%s
A=M
M=0
@%s
M=M+1
`, STACK_ASSEMBLY, STACK_ASSEMBLY)
		}
	}
	return cleanLocal
}

func (c *CodeWriter) WriteCall(name string, args int) string {
	c.labelCounter++
	returnLabel := fmt.Sprintf(`return-address%s`, strconv.Itoa(c.labelCounter))
	return writeCall(name, returnLabel, args)
}
func writeCall(name string, returnLabel string, numArgs int) string {

	pushRet := fmt.Sprintf(`@%s
D=A
@%s
A=M
M=D
@%s
M=M+1
`, returnLabel, STACK_ASSEMBLY, STACK_ASSEMBLY)

	pushLocal := fmt.Sprintf(`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, vmLabelsToAssemblyLabel[LOCAL], STACK_ASSEMBLY, STACK_ASSEMBLY)

	pushArg := fmt.Sprintf(`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, vmLabelsToAssemblyLabel[ARGS], STACK_ASSEMBLY, STACK_ASSEMBLY)

	pushThis := fmt.Sprintf(`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, vmLabelsToAssemblyLabel[THIS], STACK_ASSEMBLY, STACK_ASSEMBLY)

	pushThat := fmt.Sprintf(`@%s
D=M
@%s
A=M
M=D
@%s
M=M+1
`, vmLabelsToAssemblyLabel[THAT], STACK_ASSEMBLY, STACK_ASSEMBLY)

	// reposition ARG: ARG = SP - 5 - args
	repositionArg := fmt.Sprintf(`@%v
D=A
@5
D=D+A
@%s
A=M
D=A-D
@%s
M=D
`, numArgs, STACK_ASSEMBLY, vmLabelsToAssemblyLabel[ARGS])

	// reposition LCL: LCL = SP
	repositionLcl := fmt.Sprintf(`@%s
D=M
@%s
M=D
`, STACK_ASSEMBLY, vmLabelsToAssemblyLabel[LOCAL])

	jumpAndRet := fmt.Sprintf(`@%s
0;JMP
(%s)
`, name, returnLabel)

	return fmt.Sprintf(
		"%s%s%s%s%s%s%s%s\n",
		pushRet,
		pushLocal,
		pushArg,
		pushThis,
		pushThat,
		repositionArg,
		repositionLcl,
		jumpAndRet,
		)
}

func (c *CodeWriter) WriteReturn() string {

	saveBaseAddress := fmt.Sprintf(`@%s
D=M
@R13
M=D
`, vmLabelsToAssemblyLabel[LOCAL])

	saveReturnAddress := fmt.Sprintf(`@R13
D=M
@5
D=D-A
A=D
D=M
@R14
M=D
`)

	positionReturnValue := fmt.Sprintf(`@%s
A=M-1
D=M
@%s
A=M
M=D
`, STACK_ASSEMBLY, vmLabelsToAssemblyLabel[ARGS])


	restoreStack := fmt.Sprintf(`@%s
D=M+1
@%s
M=D
`, vmLabelsToAssemblyLabel[ARGS], STACK_ASSEMBLY)



	restoreThat := fmt.Sprintf(`@R13
D=M
@1
A=D-A
D=M
@%s
M=D
`, vmLabelsToAssemblyLabel[THAT])

	restoreThis := fmt.Sprintf(`@R13
D=M
@2
A=D-A
D=M
@%s
M=D
`, vmLabelsToAssemblyLabel[THIS])

	restoreArgs := fmt.Sprintf(`@R13
D=M
@3
A=D-A
D=M
@%s
M=D
`, vmLabelsToAssemblyLabel[ARGS])

	restoreLocal := fmt.Sprintf(`@R13
D=M
@4
A=D-A
D=M
@%s
M=D
`, vmLabelsToAssemblyLabel[LOCAL])

	returnToCallee := fmt.Sprintf(`@R14
A=M
0;JMP
`)

	// position return value and restore SP
	return fmt.Sprintf(
		`%s%s%s%s%s%s%s%s%s`,
		saveBaseAddress,
		saveReturnAddress,
		positionReturnValue,
		restoreStack,
		restoreThat,
		restoreThis,
		restoreArgs,
		restoreLocal,
		returnToCallee,
		)
}
