package Assembler
import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Command string
const(
	A_COMMAND_NUM Command = "a_command_num"
	A_COMMAND_VAR Command = "a_command_var"
	C_COMMAND Command = "c_command"
	L_COMMAND Command = "l_command"
)

type Parser struct {
	raw    string
	lines  []string
	rowIdx int
	memIdx int
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

	return &Parser{raw: str, lines: lines, rowIdx: 0, memIdx: 0}, nil
}

func trimComment(s string) string {
	if idx := strings.Index(s, "//"); idx != -1 {
		return s[:idx]
	}
	return s
}




func (p *Parser) hasMoreCommands() bool {
	return p.rowIdx < len(p.lines) - 1
}

func (p *Parser) Advance() bool {
	if p.hasMoreCommands() {
		if p.CommandType() != L_COMMAND {
			p.memIdx++
		}
		p.rowIdx++
		return true
	}
	return false
}

func (p *Parser) Actual() (Command, int) {
	return Command(p.lines[p.rowIdx]), p.memIdx
}

func (p *Parser) Reset() {
	p.rowIdx = 0
	p.memIdx = 0
}

func (p *Parser) CommandType() Command {
	var line = p.lines[p.rowIdx]

	if line[0] == '@' {
		if _, ok := strconv.Atoi(line[1:]); ok == nil {
			return A_COMMAND_NUM
		}
		return A_COMMAND_VAR
	}
	if line[0] == '(' && line[len(line) - 1] == ')' {
		return L_COMMAND
	}
	return C_COMMAND
}

func (p *Parser) Symbol() string {
	var cmd, _ = p.Actual()
	var cmdType = p.CommandType()

	if cmdType == A_COMMAND_VAR || cmdType == A_COMMAND_NUM {
		return string(cmd[1:])
	}
	if cmdType == L_COMMAND {
		return string(cmd[1:len(cmd)-1])
	}
	panic("No symbol in command of type: " + cmdType)
}





