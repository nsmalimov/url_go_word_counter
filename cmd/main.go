package main

import (
	"url_go_word_counter/internal/config"
	"url_go_word_counter/internal/processor"
	"url_go_word_counter/internal/stdin_parser"
)

func main() {
	cfg := &config.Config{
		MaxCountGoroutines:       5,
		TimeOutGetSiteContentSec: 2,
	}

	stdInParser := stdin_parser.StdInParser{}
	stdInParser.ReadStdIn()

	prc := processor.New(cfg, nil)
	prc.Init()

	prc.Processing(stdInParser.Urls)
}
