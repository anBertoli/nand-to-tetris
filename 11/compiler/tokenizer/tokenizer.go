package tokenizer

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type Tokenizer struct {
	source string
	raws   []string
	tokens []Token
}

func NewTokenizer(path string) (*Tokenizer, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	str := string(in)
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\r", "\n")

	// clean inline comments: "//"
	var lines []string
	for _, l := range strings.Split(str, "\n") {
		l = strings.TrimSpace(trimInlineComment(l))
		if l != "" {
			lines = append(lines, l)
		}
	}
	str = strings.Join(lines, "\n")

	// clean multiline comments ("/* */") and put everything in one line
	str = trimMultilineComment(str)
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.TrimSpace(str)

	return &Tokenizer{
		source: str,
	}, nil
}

func trimInlineComment(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}

func trimMultilineComment(s string) string {
	var startIdx = strings.Index(s, "/*")
	for startIdx != -1 {
		endIdx := strings.Index(s, "*/")
		if endIdx > -1 && endIdx > startIdx {
			s = s[0:startIdx] + s[endIdx+2:]
		} else {
			s = s[0:startIdx]
		}
		startIdx = strings.Index(s, "/*")
	}
	return s
}

func buildString(str string) string {
	var tmp = "\""
	for _, c := range str[1:] {
		tmp += string(c)
		if c == '"' {
			break
		}
	}
	return tmp
}

func (t *Tokenizer) HasMoreTokens() bool {
	if len(t.source) == 0 {
		return false
	}
	for _, c := range t.source {
		// newlines have already been removed
		if !unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func (t *Tokenizer) Next() Token {
	var rawToken string
	var token Token

	// build raw token
	for i, c := range t.source {
		if i == 0 && c == '"' {
			rawToken = buildString(t.source)
			break
		}
		if i == 0 && IsSymbol(c) {
			rawToken = string(c)
			break
		}
		// else identifier or int/str constant
		if c == ' ' || IsSymbol(c) {
			break
		}
		rawToken += string(c)
	}

	// build token
	tokenType, tokenValue, err := t.tokenTypeAndVal(rawToken)
	if err != nil {
		log.Fatal(err)
	}
	token = Token{
		Raw:      rawToken,
		Type:     tokenType,
		Value:    tokenValue,
		Children: &[]Token{},
	}

	// update tokenizer
	t.source = strings.TrimSpace(t.source[len(rawToken):])
	t.raws = append(t.raws, rawToken)
	t.tokens = append(t.tokens, token)
	return token
}

func (t *Tokenizer) Tokens() []Token {
	return t.tokens
}

func (t *Tokenizer) tokenTypeAndVal(raw string) (TokenType, TokenValue, error) {
	firstRune := rune(raw[0])
	lastRune := rune(raw[len(raw)-1])

	if _, ok := symbols[firstRune]; ok && len(raw) == 1 {
		return SYMBOL, firstRune, nil
	}
	if _, ok := keywords[raw]; ok {
		return KEYWORD, raw, nil
	}
	if firstRune == '"' && lastRune == '"' {
		return STRINGCONST, raw[1:len(raw)-1], nil
	}
	if num, err := strconv.Atoi(raw); err == nil {
		return INTCONST, num, nil
	}
	if unicode.IsLetter(firstRune) {
		return IDENTIFIER, raw, nil
	}
	return INVALID, nil, errors.New("not a valid token type")
}
