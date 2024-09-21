package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"com.lanuage/monkey/lexer"
	"com.lanuage/monkey/parser"
)

func Repl() {
	const PROMPT = ">> "

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		if len(p.Errors()) > 0 {
			printParseErrors(os.Stderr, p.Errors())
		}

		io.WriteString(os.Stdout, p.ParserProgram().String())
		io.WriteString(os.Stdout, "\n")

	}
}

func printParseErrors(out io.Writer, errors []string) {

	for _, err := range errors {
		io.WriteString(out, err+"\n")
	}
}
