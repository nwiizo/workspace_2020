package repl

import (
	"bufio"
	"fmt"
	"github.com/nwiizo/workspace_2020/go114/waiig_code_1.4/01/monkey/lexer"
	"github.com/nwiizo/workspace_2020/go114/waiig_code_1.4/01/monkey/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
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
