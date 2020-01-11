package VM_translator

import (
	"errors"
	"io/ioutil"
	"strings"
)

type Parser struct {
	lines      []string
	rowIdx     int
	actualFunc string
}

func NewParser(path string) (*Parser, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	str := string(in)

	// trim spaces and comments
	var lines []string
	for _, l := range strings.Split(str, "\n") {
		l = strings.TrimSpace(l)
		l = trimComment(l)
		if l != "" {
			lines = append(lines, l)
		}
	}
	return &Parser{lines: lines, rowIdx: 0, actualFunc: "UNDEFINED"}, nil
}
func trimComment(s string) string {
	if idx := strings.Index(s, "//"); idx != -1 {
		return s[:idx]
	}
	return s
}

func (p *Parser) Actual() string {
	return p.lines[p.rowIdx]
}
func (p *Parser) ActualFunc() string {
	return p.actualFunc
}
func (p *Parser) Advance() bool {
	if p.rowIdx < len(p.lines)-1 {
		p.rowIdx++
		return true
	}
	return false
}

func (p *Parser) CommandType() VMCommand {
	var fields = strings.Split(p.Actual(), " ")
	switch fields[0] {
	case ADD:
		fallthrough
	case SUB:
		fallthrough
	case NEG:
		fallthrough
	case EQ:
		fallthrough
	case GT:
		fallthrough
	case LT:
		fallthrough
	case AND:
		fallthrough
	case OR:
		fallthrough
	case NOT:
		return C_ARITHMETIC
	case PUSH:
		return C_PUSH
	case POP:
		return C_POP
	case LABEL:
		return C_LABEL
	case GOTO:
		return C_GOTO
	case IFGOTO:
		return C_IFGOTO
	case FUNCTION:
		p.actualFunc = fields[1]
		return C_FUNCTION
	case CALL:
		return C_CALL
	case RETURN:
		return C_RETURN
	}
	panic("Unknown command type.")
}

func (p *Parser) Arg1() (string, error) {
	var fields = strings.Split(p.Actual(), " ")
	var cmdType = p.CommandType()
	switch cmdType {
	case C_ARITHMETIC:
		return fields[0], nil
	case C_POP:
		fallthrough
	case C_PUSH:
		fallthrough
	case C_LABEL:
		fallthrough
	case C_GOTO:
		fallthrough
	case C_IFGOTO:
		fallthrough
	case C_FUNCTION:
		fallthrough
	case C_CALL:
		return fields[1], nil
	}
	return "", errors.New("invalid command type (" + p.lines[p.rowIdx] + ") in Arg1 method.")
}

func (p *Parser) Arg2() (string, error) {
	var fields = strings.Split(p.Actual(), " ")
	var cmdType = p.CommandType()
	switch cmdType {
	case C_POP:
		fallthrough
	case C_PUSH:
		fallthrough
	case C_FUNCTION:
		fallthrough
	case C_CALL:
		n := fields[2]
		return n, nil
	}
	return "", errors.New("invalid command type (" + p.lines[p.rowIdx] + ") in Arg2 method")
}
