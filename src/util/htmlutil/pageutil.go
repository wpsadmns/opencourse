package htmlutil

import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"net/http"
	"regexp"
	"strings"
)

var (
	escapeRx  *regexp.Regexp
	scriptRx  *regexp.Regexp
	styleRx   *regexp.Regexp
	nodeRx    *regexp.Regexp
	lineRx    *regexp.Regexp
	charsetRx *regexp.Regexp
	spaceRx   *regexp.Regexp
)

var line []byte = []byte{13, 10}

const (
	CharsetStr = "charset="
)

func init() {
	escapeRx = regexp.MustCompile(`(?i)&quot;|&nbsp;`)
	scriptRx = regexp.MustCompile(`(?i)<script[^>]*?>[^<]*?</script>`)
	styleRx = regexp.MustCompile(`(?i)<style[^>]*?>[^<]*?</style>`)
	nodeRx = regexp.MustCompile(`(?i)<[^>]+?>`)
	// lineRx = regexp.MustCompile(`(?:\r?\n(?:\s*?\r?\n)+)`)
	lineRx = regexp.MustCompile(`(?:\r?\n(?:\s*?\r?\n)+)`)
	spaceRx = regexp.MustCompile(`[ \t]`)
	charsetRx = regexp.MustCompile(`(?i)<meta[^>]*charset\s*=\s*['"]?([^'"]*?)['"]?[^>]*>`)
}

func CharsetFromResponse(response *http.Response) (string, error) {
	if contentType := response.Header.Get("Content-Type"); contentType != "" {
		if idx := strings.Index(contentType, CharsetStr); idx != -1 {
			return strings.Trim(contentType[idx+len(CharsetStr):], " "), nil
		}
	}
	return "", fmt.Errorf("can not find charset from response.")
}

func CharsetFromHTML(html []byte) (string, error) {
	if submatch := charsetRx.FindSubmatch(html); len(submatch) == 2 {
		return strings.Trim(string(submatch[1]), " "), nil
	}
	return "", fmt.Errorf("can't find charset from html code.")
}

func Charset(response *http.Response, html []byte) (string, error) {
	charset, err := CharsetFromResponse(response)
	if err != nil {
		charset, err = CharsetFromHTML(html)
	}
	return charset, err
}

// convert the coding(any coding, such as Big5, ISO-8859-1, GBK) to utf-8
func ConvertToUTF8(text []byte, coding string) []byte {
	output := make([]byte, len(text)*5)
	iconv.Convert(text, output, coding, "UTF-8")
	return output
}

// such as <script type="text/javascript">....</script>
func FilterJavascriptTag(html []byte) []byte {
	return scriptRx.ReplaceAll(html, nil)
}

// such as <style>....</style>
func FilterStyleTag(html []byte) []byte {
	return styleRx.ReplaceAll(html, nil)
}

// e.g. <li>, </li>, <head>, </head>, <body>, </body> ...
func FilterHtmlNode(html []byte) []byte {
	return nodeRx.ReplaceAll(html, nil)
}

func FilterNullLine(html []byte) []byte {
	return lineRx.ReplaceAll(html, line)
}

func FilterSpace(html []byte) []byte {
	return spaceRx.ReplaceAll(html, line)
}

func Filter(html []byte) []byte {
	// attention that the deleted sequences, the html node regular expression must be last removed.
	// return FilterNullLine(FilterSpace(FilterHtmlNode(
	// 	FilterJavascriptTag(
	// 		FilterStyleTag(html)))))
	return FilterNullLine(
		FilterSpace(
			FilterHtmlNode(
				FilterJavascriptTag(
					FilterStyleTag(html)))))
}
