package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"url_go_word_counter/internal/config"
)

type result struct {
	Url        string
	Count      int
	ErrorCause string
}

type processor struct {
	config     *config.Config
	resultChan chan result
	client     *http.Client
	closeChan  chan bool
	stopChan   chan bool
	resByUrls  sync.Map
}

func New(config *config.Config, client *http.Client) *processor {
	if client == nil {
		client = &http.Client{Timeout: time.Duration(config.TimeOutGetSiteContentSec) * time.Second}
	}

	processor := processor{
		config:     config,
		resultChan: make(chan result),
		stopChan:   make(chan bool),
		closeChan:  make(chan bool),
		client:     client,
	}

	return &processor
}

func (t *processor) Init() {
	go func() {
		<-t.stopChan
		t.closeChan <- true
	}()
}

func (t *processor) Processing(urls []string) {
	for _, url := range urls {
		t.resByUrls.Store(url, -2)
	}

	var countGoroutinesStarted int

	var allCount int
	// запускает максимум 5
	t.resByUrls.Range(func(url, resCount interface{}) bool {
		if countGoroutinesStarted < 5 {
			t.resByUrls.Store(url, -1)
			go t.parseSiteByUrl(url.(string))
			countGoroutinesStarted++
		}

		allCount++
		return true
	})

	var countReceive int

L:
	for {
		select {
		case msg := <-t.resultChan:
			// я не стал в мапу класть структуру, чтобы не усложнять, -1 - флаг того, что была ошибка при get
			if msg.ErrorCause != "" {
				t.resByUrls.Store(msg.Url, -1)
			} else {
				t.resByUrls.Store(msg.Url, msg.Count)
			}

			countReceive++
			countGoroutinesStarted--

			if countReceive == allCount {
				t.stopChan <- true
			}

			// горутина освободилась, запускаем
			t.resByUrls.Range(func(url, resCount interface{}) bool {
				if resCount == -2 {
					t.resByUrls.Store(url, -1)
					go t.parseSiteByUrl(url.(string))
					countGoroutinesStarted++

					return false
				}

				return true
			})
		case <-t.closeChan:
			break L
		}
	}

	t.close()

	t.showResult(urls)
}

func (t *processor) parseSiteByUrl(url string) {
	result := result{}

	resp, err := t.client.Get(url)

	if err != nil {
		log.Printf("Error when try http.Get, siteUrl: %s, err: %s\n", url, err)
		result.ErrorCause = err.Error()
		t.resultChan <- result
		return
	}

	defer func() {
		err = resp.Body.Close()

		if err != nil {
			log.Printf("Error when try resp.Body.close, err: %s\n", err)
		}
	}()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error when try ioutil.ReadAll, siteUrl: %s, err: %s\n", url, err)
		return
	}

	countMatched := t.countMatches(string(html))

	result.Count = countMatched
	result.Url = url

	t.resultChan <- result
}

func (t *processor) countMatches(text string) int {
	aORb := regexp.MustCompile("\\bGo\\b")

	matches := aORb.FindAllStringIndex(text, -1)

	return len(matches)
}

func (t *processor) close() {
	close(t.resultChan)
	close(t.stopChan)
	close(t.closeChan)
}

func (t *processor) showResult(urls []string) {
	var total int

	for _, url := range urls {
		if count, ok := t.resByUrls.Load(url); ok {
			if count == -1 {
				log.Printf("While try read data from: %s, something was wrong (http error [timeout, not found, etc ..])", url)
			} else {
				fmt.Printf("Count for %s: %d\n", url, count)
				total += count.(int)
			}
		}
	}

	fmt.Printf("Total: %d\n", total)
}
