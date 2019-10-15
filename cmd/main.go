package main

import (
	"bufio"
	"log"
	"os"

	"url_go_word_counter/config"
	"url_go_word_counter/processor"
)

func main() {
	cfg := &config.Config{
		MaxCountGoroutines: 5,
	}

	var urls []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if scanner.Err() != nil {
		log.Fatalf("Error when try read from input[main], err: %s", scanner.Err())
	}

	prc := processor.New(cfg)

	err := prc.Processing(urls)

	if err != nil {
		log.Fatal("Error when try prc.GetSitesByUrls[main], err: %s", err)
	}
}
