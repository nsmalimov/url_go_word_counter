package processor

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"url_go_word_counter/internal/config"
)

type TestCaseGetAndParseUrl struct {
	url            string
	minCountResult int
	maxCountResult int
}

type TestCaseCountMatchGoWord struct {
	text     string
	resCount int
}

type TestCasesParallelWork struct {
	urls        []string
	result      []int
	timeOutCase bool
}

func TestGetAndParseUrl(t *testing.T) {
	cfg := config.Config{}

	processor := New(&cfg, nil)

	testCaseGetAndParseUrl := []TestCaseGetAndParseUrl{
		{
			url:            "https://golang.org",
			minCountResult: 1,
			maxCountResult: 30,
		},
		{
			url:            "https://google.com",
			minCountResult: 0,
			maxCountResult: 0,
		},
	}

	for _, testCase := range testCaseGetAndParseUrl {
		go processor.parseSiteByUrl(testCase.url)

		result := <-processor.resultChan

		if result.Count < testCase.minCountResult || result.Count > testCase.maxCountResult {
			t.Fail()
		}
	}

	processor.close()
}

func TestGetCountMatches(t *testing.T) {
	testCaseCountMatchGoWord := []TestCaseCountMatchGoWord{
		{
			text:     "Go world only one exist",
			resCount: 1,
		},
		{
			text:     "wGo world go only one exist",
			resCount: 0,
		},
		{
			text:     "wGo world go only one exist\nfhfhf hfhfjf Go fjfjfj",
			resCount: 1,
		},
	}

	cfg := config.Config{}

	processor := New(&cfg, nil)

	for _, elem := range testCaseCountMatchGoWord {
		res := processor.countMatches(elem.text)

		if res != elem.resCount {
			t.Fail()
		}
	}
}

func TestParallelWork(t *testing.T) {
	testCases := []TestCasesParallelWork{
		{
			urls: []string{
				"https://golang.org",
				"https://golang.org",
			},
			result: []int{15, 15},
		},
		{
			urls: []string{
				"https://golang.org",
				"https://google.ru",
			},
			result: []int{15, 0},
		},
		{
			urls: []string{
				"https://google.ru",
				"https://golang.org",
				"https://mail.ru",
				"https://gobyexample.com",
				"https://stackoverflow.com",
			},
			result: []int{0, 15, 0, 5, 0},
		},
		{
			urls: []string{
				"https://google.ru",
				"https://golang.org",
				"https://mail.ru",
				"https://gobyexample.com",
				"https://stackoverflow.com",
				"https://vk.com",
				"https://habr.com",
				"https://hh.ru",
			},
			result: []int{0, 15, 0, 5, 0, 0, 0, 0},
		},
		// timeout
		{
			urls: []string{
				"https://google.ru",
			},
			timeOutCase: true,
			result:      []int{-1},
		},
		// not found
		{
			urls: []string{
				"https://go333ogle.ru",
			},
			result: []int{-1},
		},
	}

	cfg := config.Config{
		MaxCountGoroutines:       5,
		TimeOutGetSiteContentSec: 2,
	}

	var wg sync.WaitGroup

	for index, testCase := range testCases {
		//  не запускать новый тест, пока не завершился предыдущий
		wg.Add(1)

		fmt.Printf("Start test case num: %d\n", index)

		go func() {
			processor := New(&cfg, nil)
			processor.Init()

			// timeout check case
			if testCase.timeOutCase {
				processor.client.Timeout = time.Duration(5) * time.Millisecond
			}

			processor.Processing(testCase.urls)

			processor.resultChan = make(chan result)

			for index, url := range testCase.urls {
				if val, ok := processor.resByUrls.Load(url); ok {
					if val.(int) != testCase.result[index] {
						// можно было мы использовать assert, но по требованию без сторонних либ
						fmt.Printf("Get: %d, but need: %d\n", val.(int), testCase.result[index])
						defer wg.Done()
						t.Fail()
						return
					}
				}
			}

			defer wg.Done()

			fmt.Printf("Final test case num: %d\n\n", index)
		}()

		wg.Wait()
	}
}
