package repl

import (
	"bufio"
	"fmt"
	"io"
	"iscript/lexer"
	"iscript/parser"
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

		io.WriteString(out, prog.String())
		io.WriteString(out, "\n")
	}
}
