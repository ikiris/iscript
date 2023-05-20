package repl

import (
	"bufio"
	"fmt"
	"io"
	"iscript/compiler"
	"iscript/lexer"
	"iscript/object"
	"iscript/parser"
	"iscript/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symTable := compiler.NewSymTable()
	for i, v := range object.Builtins {
		symTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Fprint(out, PROMPT)
		scanned := s.Scan()
		if !scanned {
			return
		}

		line := s.Text()
		l := lexer.New(line)
		p := parser.New(l)

		prog, err := p.ParseProgram()
		if err != nil {
			fmt.Fprintf(out, "err: %v", err)
			continue
		}

		c := compiler.NewWithState(symTable, constants)
		err = c.Compile(prog)
		if err != nil {
			fmt.Fprintf(out, "Whoops!: Compile failed:\n %s\n", err)
			continue
		}

		code := c.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops!: Bytecode excecution failed:\n %s\n", err)
			continue
		}

		lastPop := machine.LastPoppedStackElem()
		io.WriteString(out, lastPop.Inspect())
		io.WriteString(out, "\n")
	}
}
