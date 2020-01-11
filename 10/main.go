package main

import (
	"./compiler"
	"flag"
	"fmt"
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

	fmt.Println(*inPath)
	f, err := os.Stat(*inPath)
	if err != nil {
		path.Dir(*inPath)
		fmt.Println("dadsdse")
		log.Fatal(err)
	}

	// single file
	if f.Mode().IsRegular() {
		outT := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_tokens.xml")
		outP := path.Join(path.Dir(*inPath), removeExt(f.Name()) + "_parsed.xml")

		tokens := tokenizeFile(*inPath)
		tokensXML := "<tokens>\n" + compiler.FormatXML(tokens, 0) + "</tokens>"
		err = ioutil.WriteFile(outT, []byte(tokensXML), 0644)
		if err != nil {
			log.Fatal(err)
		}

		parsed := parseTokens(tokens)
		tokensXML = compiler.FormatXML(parsed, 0)
		err = ioutil.WriteFile(outP, []byte(tokensXML), 0644)
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

				tokens := tokenizeFile(inputPath)
				tokensXML := "<tokens>\n" + compiler.FormatXML(tokens, 0) + "</tokens>"
				err = ioutil.WriteFile(outT, []byte(tokensXML), 0644)
				if err != nil {
					log.Fatal(err)
				}

				parsed := parseTokens(tokens)
				tokensXML = compiler.FormatXML(parsed, 0)
				err = ioutil.WriteFile(outP, []byte(tokensXML), 0644)
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
func tokenizeFile(filepath string) []compiler.Token {
	tokenMachine, err := compiler.NewTokenizer(filepath)
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

func parseTokens(tokens []compiler.Token) []compiler.Token {
	parser := compiler.NewParser(tokens)
	parsedTokens, _ := parser.CompileClass()
	return []compiler.Token{parsedTokens}
}

func removeExt(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
