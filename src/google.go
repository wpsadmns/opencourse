package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"regexp"
	"strconv"
	iconv "github.com/djimenez/iconv-go"
)

var hrefRx *regexp.Regexp = regexp.MustCompile(`<cite>([^<]+)</cite>`)

func main() {
	urls := URLfromGOOGLE()
	curl := make(chan string, 30)
	go func() {
		for _, url := range urls {
			curl <- url
		}
		fmt.Println("[INFO] close(curl)..")
		close(curl)
	}()

	workers := 8
	done := make(chan struct{}, workers)
	ccurl := make(chan string, 30)
	for i := 0; i < workers; i++ {
		go func() {
			for {
				if u, ok := <-curl; ok {
					URL, _ := url.Parse(u)
					if URL.Scheme == "" {
						u = "http://" + u
					}
					_, err := http.Get(u)
					if err == nil {
						ccurl <- u
					}
				} else {
					break
				}
			}
			done<-struct{}{}
		}()
	}
	go func() {
		for i := 0; i < workers; i++ {
			<-done
		}
		close(ccurl)
		fmt.Println("[INFO] close(ccurl)...")
	}()

	for link := range ccurl {
		if link[len(link)-1] == '/' {
			fmt.Println(link[:len(link)-1])
		} else {
			fmt.Println(link)
		}
	}
}

func URLfromGOOGLE() (urls []string) {
	html1 := HTMLfromGOOGLE(0)
	html2 := HTMLfromGOOGLE(100)
	for _, submatch := range hrefRx.FindAllSubmatch(html1, -1) {
		urls = append(urls, string(submatch[1]))
	}
	for _, submatch := range hrefRx.FindAllSubmatch(html2, -1) {
		urls = append(urls, string(submatch[1]))
	}
	return
}

func HTMLfromGOOGLE(pos int) []byte {
	url := "https://www.google.com.hk/search?start="+ strconv.Itoa(pos) +"&num=100&newwindow=1&safe=strict&q=intitle%3A%E5%85%AC%E5%BC%80%E8%AF%BE&oq=intitle%3A%E5%85%AC%E5%BC%80%E8%AF%BE&gs_l=serp.12...0.0.0.1628503.0.0.0.0.0.0.0.0..0.0....0...1c..38.serp..0.0.0._NuhQqY0huk"
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	input, err := ioutil.ReadAll(resp.Body)
	check(err)
	output := make([]byte, len(input) * 20)
	iconv.Convert(input, output, "Big5", "UTF-8")
	return output
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
