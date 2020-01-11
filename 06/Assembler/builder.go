package Assembler
import (
	"strconv"
	"strings"
)


var compMap = map[string]string{
	"0": "0101010",
	"1": "0111111",
	"-1": "0111010",
	"D": "0001100",
	"A": "0110000",
	"!D": "0001101",
	"!A": "0110001",
	"-D": "0001111",
	"-A": "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",

	"M": "1110000",
	"!M": "1110001",
	"-M": "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

var destMap = map[string]string{
	"": "000",
	"M": "001",
	"D": "010",
	"A": "100",
	"MD": "011",
	"AM": "101",
	"AD": "110",
	"AMD": "111",
}

var jumpMap = map[string]string{
	"": "000",
	"JGT": "001",
	"JEG": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

type Builder struct{
	binary 		string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func splitC (c Command) (string, string, string) {
	var dest, comp, jump = "", "", ""
	str := string(c)

	if idx := strings.Index(str, "="); idx != -1 {
		dest = str[:idx]
		str = str[idx+1:]
	}

	if idx := strings.Index(str, ";"); idx != -1 {
		comp = str[:idx]
		jump = str[idx+1:]
	} else {
		comp = str
	}

	//fmt.Printf("dest=%s, comp=%s, jump=%s", dest, comp, jump)
	return strings.TrimSpace(dest), strings.TrimSpace(comp), strings.TrimSpace(jump)
}


func (b *Builder) BuildC (c Command) {
	dest, comp, jump := splitC(c)
	//fmt.Println(c)
	//fmt.Printf("dest=%s, comp=%s, jump=%s\n", dest, comp, jump)
	//fmt.Printf("comp=%s, dest=%s, jump=%s\n", compMap[comp] ,destMap[dest] ,jumpMap[jump])
	b.binary += "111" + compMap[comp] + destMap[dest] + jumpMap[jump] + "\n"
}




func (b *Builder) BuildA(c Command) {
	str := string(c)
	n, err := strconv.Atoi(str[1:])
	if err != nil {
		panic(err)
	}

	binStr := string(strconv.FormatInt(int64(n), 2))
	for i := len(binStr); i < 15; i++ {
		binStr = "0" + binStr
	}
	b.binary +=  "0" + binStr + "\n"
}

func (b *Builder) Print() string {
	return b.binary
}


