package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"com.language/monkey/evaluator"
	"com.language/monkey/lexer"
	"com.language/monkey/object"
	"com.language/monkey/parser"
)

func Repl() {
	const PROMPT = ">> "

	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(os.Stdout, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParserProgram()

		if len(p.Errors()) > 0 {
			printParseErrors(os.Stderr, p.Errors())
			continue
		}
		result := evaluator.Eval(program, env)

		if result != nil {
			io.WriteString(os.Stdout, result.Inspect())
			io.WriteString(os.Stdout, "\n")
			io.WriteString(os.Stdout, program.String())
			io.WriteString(os.Stdout, "\n")
		}

	}
}

func printParseErrors(out io.Writer, errors []string) {

	for _, err := range errors {
		io.WriteString(out, err+"\n")
	}
}
