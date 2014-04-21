package htmlutil

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestCharset(t *testing.T) {
	// charset is Big5
	url := "http://www.google.com"
	resp, err := http.Get(url)
	check(err, t)
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	check(err, t)

	charset, err := Charset(resp, html)
	check(err, t)
	if charset == "" || (charset != "Big5" && charset != "big5") {
		t.Errorf("Charset() return %s, want Big5.", charset)
	}
}

func TestRemoveAll(t *testing.T) {
	file, err := os.Open("./test.html")
	check(err, t)
	html, err := ioutil.ReadAll(file)
	check(err, t)
	theHtml := RemoveAll(html)

	if len(html) == len(theHtml) {
		t.Errorf("Remove relative node failed")
	}
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
