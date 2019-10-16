package stdin_parser

import (
	"bufio"
	"os"
)

type StdInParser struct {
	Urls []string
}

func (t *StdInParser) ReadStdIn() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t.Urls = append(t.Urls, scanner.Text())
	}
}
