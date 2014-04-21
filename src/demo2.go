package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"util/htmlutil"
)

func main() {
	resp, _ := http.Get("http://search.51job.com/job/59126997,c.html")
	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	html1 := htmlutil.ConvertToUTF8(html, "gb2312")
	fmt.Println("size:", size(html), size(html1))
	fmt.Println("len:", len(html), len(html1))
	fmt.Println("cap:", cap(html), cap(html1))
	fmt.Println(string(html1))
	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Println(string(htmlutil.Filter(html1)))
}

func size(b []byte) int {
	sum := 0
	for _, _ = range b {
		sum++
	}
	return sum
}
