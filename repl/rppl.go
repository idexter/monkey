package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/idexter/monkey/lexer"
	"github.com/idexter/monkey/parser"
)

// StartRPPL implements Read-Parse-Print-Loop.
func StartRPPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
