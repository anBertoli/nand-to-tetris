package main

import (
	"./VM_translator"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Compiled struct {
	assembly string
	filename string
}

func main() {
	inPath := flag.String("in", "", "Path of directory containing input assembly files.")
	outPath := flag.String("out", "", "Path of output binary file.")
	flag.Parse()

	// read files in directory
	files, err := ioutil.ReadDir(*inPath)
	if err != nil {
		log.Fatal(err)
	}

	// parallel compilation
	var compiledFiles []Compiled
	var compWg sync.WaitGroup
	var compChan = make(chan Compiled, 10)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".vm" {
			filePath, _ := filepath.Abs(path.Join(*inPath, f.Name()))
			compWg.Add(1)
			go func() {
				compileFile(filePath, compChan)
				compWg.Done()
			}()
		}
	}
	go func() {
		compWg.Wait()
		fmt.Println("All files compiled to assembly.")
		close(compChan)
	}()

	// wait assembly results
	for compFile := range compChan {
		fmt.Println("Done: " + compFile.filename)
		compiledFiles = append(compiledFiles, compFile)
	}



	// put files together
	var assembly = VM_translator.WriteBootstrap()
	for _, c := range compiledFiles {
		assembly += "// ---------- " + c.filename + " ---------- //\n"
		assembly += c.assembly
	}
	assembly += VM_translator.WriteEndLoop()


	// write to disk
	f, err := os.Create(*outPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(assembly))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Output to: " + *outPath)
}


func compileFile (absFilePath string, resChan chan Compiled) {
	fmt.Println("Working on: " + absFilePath + "...")
	fileName := filepath.Base(absFilePath)
	file := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	parser, err := VM_translator.NewParser(absFilePath)
	if err != nil {
		log.Fatal(err)
	}
	codeWriter := VM_translator.NewCodeWriter()

	// translate
	var assembly string
	for {
		assembly += "\n// " + parser.Actual() + "\n"
		switch parser.CommandType() {
		case VM_translator.C_ARITHMETIC:
			op, _ := parser.Arg1()
			assembly += codeWriter.WriteArithmetic(op)

		case VM_translator.C_PUSH:
			arg1, _ := parser.Arg1()
			arg2, _ := parser.Arg2()
			assembly += codeWriter.WritePushPop(VM_translator.PUSH, arg1, arg2, file)
		case VM_translator.C_POP:
			arg1, _ := parser.Arg1()
			arg2, _ := parser.Arg2()
			assembly += codeWriter.WritePushPop(VM_translator.POP, arg1, arg2, file)

		case VM_translator.C_LABEL:
			label, _ := parser.Arg1()
			assembly += codeWriter.WriteLabel(label, parser.ActualFunc())
		case VM_translator.C_GOTO:
			label, _ := parser.Arg1()
			assembly += codeWriter.WriteGoto(label, parser.ActualFunc())
		case VM_translator.C_IFGOTO:
			label, _ := parser.Arg1()
			assembly += codeWriter.WriteIfGoto(label, parser.ActualFunc())

		case VM_translator.C_FUNCTION:
			funcName, _ := parser.Arg1()
			args, _ := parser.Arg2()
			numArgs, _ := strconv.Atoi(args)
			assembly += codeWriter.WriteFunction(funcName, numArgs)
		case VM_translator.C_CALL:
			funcName, _ := parser.Arg1()
			args, _ := parser.Arg2()
			numArgs, _ := strconv.Atoi(args)
			assembly += codeWriter.WriteCall(funcName, numArgs)
		case VM_translator.C_RETURN:
			assembly += codeWriter.WriteReturn()
		}

		if !parser.Advance() {
			break
		}
	}

	resChan <- Compiled{
		assembly: assembly,
		filename: fileName,
	}
}