package processor

import (
	"fmt"
	"io/ioutil"
	"log"

	"net/http"

	"url_go_word_counter/config"
)

type Processor struct {
	Config *config.Config
}

func New(config *config.Config) *Processor {
	processor := Processor{
		Config: config,
	}

	return &processor
}

func (t *Processor) Processing(urls []string) (err error) {
	var countGoroutinesStarted int

	for _, url := range urls {
		go t.GetSiteContentByGetRequest(url)
		countGoroutinesStarted++
	}

	// wait all goroutines

	return nil
}

func (t *Processor) GetSiteContentByGetRequest(url string) {
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer func() {
		err = resp.Body.Close()
		log.Printf("Error when try resp.Body.Close")
	}()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", html)
}
