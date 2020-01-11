package main

import (
	"./parser"
	"./tokenizer"
	"./writer"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	inPath := flag.String("in", "", "Path of directory or file containing Jack code.")
	flag.Parse()

	f, err := os.Stat(*inPath)
	if err != nil {
		path.Dir(*inPath)
		log.Fatal(err)
	}

	// single file
	if f.Mode().IsRegular() {
		outT := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_tokens.xml")
		outP := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_parsed.xml")
		outVM := path.Join(path.Dir(*inPath), removeExt(f.Name()) + ".vm")

		tokens := tokenizeFile(*inPath)
		tokensXML := "<tokens>\n" + FormatXML(tokens, 0) + "</tokens>"
		err = ioutil.WriteFile(outT, []byte(tokensXML), 0644)
		if err != nil {
			log.Fatal(err)
		}

		parsed := parseTokens(tokens)
		tokensXML = FormatXML(parsed, 0)
		err = ioutil.WriteFile(outP, []byte(tokensXML), 0644)
		if err != nil {
			log.Fatal(err)
		}

		compiled := compileParsedTokens(parsed)
		err = ioutil.WriteFile(outVM, []byte(compiled), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// directory
	if f.Mode().IsDir() {
		files, _ := ioutil.ReadDir(*inPath)

		for _, f := range files {
			if filepath.Ext(f.Name()) == ".jack" {
				inputPath := path.Join(*inPath, f.Name())

				outT := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_tokens.xml")
				outP := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_parsed.xml")
				outVM := path.Join(path.Dir(*inPath), removeExt(f.Name()) + ".vm")

				tokens := tokenizeFile(inputPath)
				tokensXML := "<tokens>\n" + FormatXML(tokens, 0) + "</tokens>"
				err = ioutil.WriteFile(outT, []byte(tokensXML), 0644)
				if err != nil {
					log.Fatal(err)
				}

				parsed := parseTokens(tokens)
				tokensXML = FormatXML(parsed, 0)
				err = ioutil.WriteFile(outP, []byte(tokensXML), 0644)
				if err != nil {
					log.Fatal(err)
				}

				compiled := compileParsedTokens(parsed)
				err = ioutil.WriteFile(outVM, []byte(compiled), 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}



/*
 short-hand function to tokenize a file
*/
func tokenizeFile(filepath string) []tokenizer.Token {
	tokenMachine, err := tokenizer.NewTokenizer(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for {
		if !tokenMachine.HasMoreTokens() {
			break
		}
		tokenMachine.Next()
	}
	return tokenMachine.Tokens()
}

func parseTokens(tokens []tokenizer.Token) []tokenizer.Token {
	newParser := parser.NewParser(tokens)
	parsedTokens, _ := newParser.ParseClass()
	return []tokenizer.Token{parsedTokens}
}

func compileParsedTokens(tokens []tokenizer.Token) string {
	newWriter := writer.NewWriter()
	compiledVm := newWriter.CompileClass(tokens[0])
	return compiledVm
}

func removeExt(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
