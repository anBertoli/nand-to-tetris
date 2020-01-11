package VM_translator

const (
	// implicit
	STACK_ASSEMBLY string = "SP"
	// in vm
	LOCAL    string = "local"
	ARGS     string = "argument"
	THIS     string = "this"
	THAT     string = "that"
	CONSTANT string = "constant"
	POINTER  string = "pointer"
	TEMP     string = "temp"
	STATIC   string = "static"
)

// vm-to-assembly mapping
var vmLabelsToAssemblyLabel = map[string]string{
	LOCAL: "LCL",
	ARGS:  "ARG",
	THIS:  "THIS",
	THAT:  "THAT",
}

// vm direct segment addresses mapping
var vmLabelsToAssemblyAddress = map[string]int{
	POINTER: 3,  // 3 this, 4 that
	TEMP:    5,  // 5 to 12
	STATIC:  16, // 16 to 255
}


type VMCommand string

const (
	C_ARITHMETIC VMCommand = "c_arithmetic"
	C_PUSH       VMCommand = "c_push"
	C_POP        VMCommand = "c_pop"
	C_LABEL      VMCommand = "c_label"
	C_GOTO       VMCommand = "c_goto"
	C_IFGOTO     VMCommand = "c_if"
	C_FUNCTION   VMCommand = "c_function"
	C_RETURN     VMCommand = "c_return"
	C_CALL       VMCommand = "c_call"
)
const (
	ADD      string = "add"
	SUB      string = "sub"
	NEG      string = "neg"
	EQ       string = "eq"
	GT       string = "gt"
	LT       string = "lt"
	AND      string = "and"
	OR       string = "or"
	NOT      string = "not"
	PUSH     string = "push"
	POP      string = "pop"
	LABEL    string = "label"
	GOTO     string = "goto"
	IFGOTO   string = "if-goto"
	FUNCTION string = "function"
	CALL     string = "call"
	RETURN   string = "return"
)

