package main

import (
	"strconv"
	"strings"
	"./tokenizer"
)

func FormatXML(tokens []tokenizer.Token, indentSpaces int) string {
	var xml string

	for _, t := range tokens {

		if t.Children != nil && len(*t.Children) > 0 {
			childXml := FormatXML(*t.Children, indentSpaces + 2)
			xml += strings.Repeat(" ", indentSpaces) + "<" + string(t.Type) + "> \n" + childXml + strings.Repeat(" ", indentSpaces) + "</" + string(t.Type) + ">\n"

		} else {
			var val string
			switch v := t.Value.(type) {
			case int:
				val = strconv.Itoa(v)
			case string:
				val = v
			case rune:
				switch v {
				case '<':
					val = "&lt;"
				case '>':
					val = "&gt;"
				case '&':
					val = "&amp;"
				default:
					val = string(v)
				}
			case nil:
				val = "\n" + strings.Repeat(" ", indentSpaces)
			default:
				panic("Type not known")
			}

			xml += strings.Repeat(" ", indentSpaces) + "<" + string(t.Type) + "> " + val + " </" + string(t.Type) + ">\n"
		}

	}

	return xml
}
