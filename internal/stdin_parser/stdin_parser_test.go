package stdin_parser

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParseStdIn(t *testing.T) {
	stdInParser := StdInParser{}

	content := []byte("https://yandex.ru\nhttps://google.com")
	tmpfile, err := ioutil.TempFile("", "example")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = os.Remove(tmpfile.Name())
	}()

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	os.Stdin = tmpfile
	stdInParser.ReadStdIn()

	if stdInParser.Urls[0] != "https://yandex.ru" || stdInParser.Urls[1] != "https://google.com" {
		t.Fail()
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}
