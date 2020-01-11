package VM_translator

import (
	"fmt"
)

func WriteBootstrap() string {

	initializeSegments := fmt.Sprintf(`// initialize memory segments
@256
D=A
@%s
M=D

@2000
D=A
@%s
M=D

@2500
D=A
@%s
M=D

@3000
D=A
@%s
M=D

@3500
D=A
@%s
M=D
`,
		STACK_ASSEMBLY,
		vmLabelsToAssemblyLabel[LOCAL],
		vmLabelsToAssemblyLabel[ARGS],
		vmLabelsToAssemblyLabel[THIS],
		vmLabelsToAssemblyLabel[THAT])

	callSysInit := "// call Sys.init\n" + writeCall("Sys.init", "return-address___sys.init", 0)

	return "// -------- bootstrap -------- //\n" + initializeSegments + callSysInit
}
