package htmlutil

import (
	"logx"
	"net/url"
	"regexp"
	"strings"
	"util/stringutil"

	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var hrefRx *regexp.Regexp = regexp.MustCompile(`(?i)<a[^>]*?href\s*?=\s*?['"]?([^'"]+)['"]?[^>]*?>`)

func URLDecoder(prefixURL *url.URL, html []byte) []string {
	complexLinks := hrefRx.FindAllSubmatch(html, -1)
	if n := len(complexLinks); n > 0 {
		links := make([]string, n)
		var (
			link string
			err  error
		)
		for i := 0; i < n; i++ {
			link = string(complexLinks[i][1])
			link, err = constructURL(prefixURL, link)
			// match algrithm
			if MatchURL(prefixURL, link) {
				links[i] = link
				logx.Logger.Printf("[INFO] Matched Link: %s", link)
			}
			if err != nil {
				logx.Logger.Println("[ERROR] Parse URL failed! error is:", err)
			}
		}
		return links
	}
	return nil
}

func MatchURL(url *url.URL, link string) bool {
	return strings.Contains(link, url.Host)
}

func constructURL(prefixURL *url.URL, link string) (string, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return stringutil.EmptyString, err
	}
	if stringutil.IsEmpty(URL.Scheme) {
		if stringutil.IsNotEmpty(URL.Path) && URL.Path[0] == '/' {
			return prefixURL.Scheme + "://" + prefixURL.Host + URL.Path, nil
		}
		link = prefixURL.String() + URL.Path
	}
	return link, nil
}
