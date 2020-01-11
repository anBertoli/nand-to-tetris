package VM_translator

import "fmt"

func (c *CodeWriter) WriteLabel(label string, actualFunc string) string {
	scopedLabel := fmt.Sprintf("%s$%s", actualFunc, label)
	return fmt.Sprintf("(%s)\n", scopedLabel)
}

func (c *CodeWriter) WriteGoto(label string, actualFunc string) string {
	scopedLabel := fmt.Sprintf("%s$%s", actualFunc, label)
	return fmt.Sprintf(`@%s
0;JMP
`, scopedLabel)
}

func (c *CodeWriter) WriteIfGoto(label string, actualFunc string) string {
	scopedLabel := fmt.Sprintf("%s$%s", actualFunc, label)
	return fmt.Sprintf(`@%s
AM=M-1
D=M
@%s
D;JGT
D;JLT
`, STACK_ASSEMBLY, scopedLabel)
}
