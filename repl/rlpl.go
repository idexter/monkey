package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/idexter/monkey/lexer"
	"github.com/idexter/monkey/token"
)

// StartRLPL implements Read-Lex-Print-Loop.
func StartRLPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
