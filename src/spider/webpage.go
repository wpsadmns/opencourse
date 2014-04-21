package spider

import (
	// iconv "github.com/djimenez/iconv-go"
	"bufio"
	"bytes"
	"constant"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
	"util/htmlutil"
)

const (
	CHARSET = "UTF-8"
	// 18 sec
	TIMEOUT = time.Second * 18

	METHOD_GET           = "GET"
	ACCEPT_CHARSET_KEY   = "Accept-Charset"
	ACCEPT_CHARSET_VALUE = "UTF-8;q=1, ISO-8859-1;q=0"
)

type Page struct {
	URL     *url.URL
	HTML    []byte
	Charset string
	Timeout int
	Header  http.Header
	Method  string
	Title   string
}

var client *http.Client

func init() {
	client = buildClient()
}

// Add restrictions on request.
// set timeout on request.
func buildClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(TIMEOUT)
				c, err := net.DialTimeout(netw, addr, TIMEOUT)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func PageFromURL(link string) (*Page, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("Error URL: link, error is %v.", link, err)
	}
	request, err := http.NewRequest(METHOD_GET, link, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequest for %s failed, error is %v.", link, err)
	}
	// only support UTF-8 encoding
	request.Header.Add(ACCEPT_CHARSET_KEY, ACCEPT_CHARSET_VALUE)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("connect to %s failed, StatusCode: %d, Status: %s.", link, response.StatusCode, response.Status)
	}
	defer response.Body.Close()

	HTML, title, err := readAllWithNoLF(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read html for %s failed, error is %v.", err)
	}
	page := Page{HTML: HTML, Title: string(title)} // , Header : resp.Header}
	charset, err := htmlutil.Charset(response, HTML)
	if err != nil {
		return nil, err
	}
	page.HTML = htmlutil.ConvertToUTF8(HTML, charset)
	page.URL = URL
	page.Charset = charset
	return &page, nil
}

func readAllWithNoLF(reader io.ReadCloser) ([]byte, []byte, error) {
	bufReader := bufio.NewReader(reader)
	bufData := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	var (
		line, theLine []byte
		title         []byte
		err           error
	)
	for {
		line, err = bufReader.ReadBytes(constant.LF)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		if theLine = bytes.TrimSpace(line); len(theLine) == 0 {
			continue
		} else if title != nil {
			title = getTitle(theLine)
		}
		bufData.Write(line)
	}
	return title, bufData.Bytes(), nil
}

var (
	title_start_upper = []byte("<TITLE>")
	title_end_upper   = []byte("</TITLE>")
	title_start_lower = []byte("<title>")
	title_end_lower   = []byte("</title>")
)

func getTitle(content []byte) []byte {
	if (bytes.HasPrefix(content, title_start_lower) || bytes.HasPrefix(content, title_start_upper)) &&
		(bytes.HasSuffix(content, title_end_lower) || bytes.HasSuffix(content, title_end_upper)) {
		return content[7 : len(content)-8]
	}
	return nil
}
