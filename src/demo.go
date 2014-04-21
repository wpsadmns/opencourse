package main

import (
	"fmt"
	"bytes"
)

var (
	title_start_upper = []byte("<TITLE>")
	title_end_upper = []byte("</TITLE>")
	title_start_lower = []byte("<title>")
	title_end_lower = []byte("</title>")
)

func main() {
	titleTag := []byte(`<title>xuchenfeng</title>`)
	fmt.Println(string(getTitle(titleTag)))
}

func getTitle(content []byte) []byte {
	if (bytes.HasPrefix(content, title_start_lower) || bytes.HasPrefix(content, title_start_upper)) && 
		(bytes.HasSuffix(content, title_end_lower) || bytes.HasSuffix(content, title_end_upper)) {
		return content[7:len(content)-8]
	}
	return nil
}
