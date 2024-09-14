package repl

import (
	"bufio"
	"fmt"
	"os"

	"com.lanuage/monkey/lexer"
	"com.lanuage/monkey/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(os.Stdout, "%+v\n", tok)
		}
	}
}
