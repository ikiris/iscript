package repl

import (
	"bufio"
	"fmt"
	"io"
	"iscript/evaluator"
	"iscript/lexer"
	"iscript/object"
	"iscript/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	env := object.NewEnvironment()

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

		evald := evaluator.Eval(prog, env)
		if evald != nil {
			io.WriteString(out, evald.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
