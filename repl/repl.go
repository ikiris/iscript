package repl

import (
	"bufio"
	"fmt"
	"io"
	"iscript/compiler"
	"iscript/lexer"
	"iscript/parser"
	"iscript/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)

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

		c := compiler.New()
		err = c.Compile(prog)
		if err != nil {
			fmt.Fprintf(out, "Whoops!: Compile failed:\n %s\n", err)
			continue
		}

		machine := vm.New(c.Bytecode())
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
