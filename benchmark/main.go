package main

import (
	"flag"
	"fmt"
	"iscript/compiler"
	"iscript/evaluator"
	"iscript/lexer"
	"iscript/object"
	"iscript/parser"
	"iscript/vm"
	"time"

	"log"
)

var engine = flag.String("engine", "vm", "use `vm` or `eval`")

var input = `
let cache = {};
let memo = fn(f, x) {
	if (cache[x] != null) {
		return cache[x];
	};
	let c = f(x);
	updateHash(cache, x, c);
	return c;
};
let fib = fn(x) {
	if (x == 0) {
		return 0;
	};
	if (x == 1) {
		return 1;
	};
	memo(fib, x - 1) + memo(fib, x - 2);
};
memo(fib, 92);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		log.Fatalf("failed to parse: %s", err)
	}

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(prog)
		if err != nil {
			log.Fatalf("failed to compile: %s", err)
		}

		machine := vm.New(comp.Bytecode())

		start := time.Now()

		err = machine.Run()
		if err != nil {
			log.Fatalf("vm error: %s", err)
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()

		result = evaluator.Eval(prog, env)
		duration = time.Since(start)
	}

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration,
	)
}
