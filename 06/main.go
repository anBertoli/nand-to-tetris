package main
import (
	"./Assembler"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	inPath := flag.String("in", "", "Path of input assembly file.")
	outPath := flag.String("out", "", "Path of output binary file.")
	flag.Parse()
	if *inPath == "" || *outPath == "" {
		log.Fatal("No input or output path.")
	}

	// instantiate parser, symbol table & builder
	parser, err := Assembler.NewParser(*inPath)
	if err != nil {
		log.Fatal(err)
	}
	symTable := Assembler.NewSymbolTable()
	builder := Assembler.NewBuilder()


	// register labels
	for {
		_, memIdx := parser.Actual()
		if parser.CommandType() == Assembler.L_COMMAND {
			label := parser.Symbol()
			if !symTable.Contains(label) {
				symTable.AddEntry(label, memIdx)
			}
		}

		// end reached
		if !parser.Advance() {
			parser.Reset()
			break
		}
	}


	// translate assembly
	fmt.Printf("%+v", symTable)
	for {
		actualCmd, _ := parser.Actual()
		cmdType := parser.CommandType()

		if cmdType == Assembler.C_COMMAND {
			builder.BuildC(actualCmd)
		}

		if cmdType == Assembler.A_COMMAND_NUM {
			builder.BuildA(actualCmd)
		}

		if cmdType == Assembler.A_COMMAND_VAR {
			sym := parser.Symbol()
			var addr int

			if !symTable.Contains(sym) {
				// it's a (new) var
				symTable.AddEntry(sym, symTable.VarMemIdx)
				addr = symTable.VarMemIdx
				symTable.VarMemIdx++
			} else {
				// it's a label
				addr = symTable.GetAddress(sym)
			}
			builder.BuildA(Assembler.Command("@" + strconv.Itoa(addr)))
		}

		// end reached
		if !parser.Advance() {
			parser.Reset()
			break
		}
	}

	// output
	binary := builder.Print()
	fmt.Println(binary)
	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(binary))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Output to:" + *outPath)
}
