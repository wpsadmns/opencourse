/*
 *
 */
package spider

import (
	"bytes"
	"constant"
	"net/url"
	"os"
	"strconv"
	"util/dateutil"
	"util/htmlutil"
)

func doAnalyzer(writer *os.File, page *Page) []string {
	links := htmlutil.URLDecoder(page.URL, page.HTML)
	saveDoc(writer, page)
	return links
}

func saveDoc(dest *os.File, page *Page) {
	buf := bytes.NewBuffer(nil)
	// begin
	writeHeadInfo(page.URL, len(page.HTML), buf)
	buf.WriteByte(constant.LF)
	buf.Write(page.HTML)
	buf.WriteByte(constant.LF)
	// write html to dest file.
	buf.WriteTo(dest)
}

func writeHeadInfo(url *url.URL, length int, buf *bytes.Buffer) {
	buf.WriteString(constant.HeadVersion)
	buf.WriteString("1.0")
	buf.WriteByte(constant.LF)

	buf.WriteString(constant.HeadURL)
	buf.WriteString(url.String())
	buf.WriteByte(constant.LF)

	buf.WriteString(constant.HeadIP)
	buf.WriteString(url.Host)
	buf.WriteByte(constant.LF)
	// TODO -
	/*
		// ipAddr, err := net.LookupIP(url.Host)
		// _, err := net.LookupIP(url.Host)
		if err == nil {
			buf.WriteString(ipAddr.String())
			// buf.WriteString(url.Host)
		} else {
			logx.Logger.Println("[ERROR] Parse IP error, host: %s, error is: %v.", url.Host, err)
			buf.WriteString(url.Host)
		}
	*/

	buf.WriteString(constant.HeadDate)
	buf.WriteString(dateutil.NowDateString())
	buf.WriteByte(constant.LF)

	buf.WriteString(constant.HeadLength)
	buf.WriteString(strconv.Itoa(length))
	buf.WriteByte(constant.LF)
}

func writeHeaderInfo(header http.Header, buf *byte.Buffer) {
	buf.writeByte(constant.LF)
	for k, vs := range header {
		buf.WriteString(k)
		buf.WriteByte(':')
		for _, v := range vs {
			buf.WriteString(v)
			buf.Write(' ')
		}
		buf.WriteByte(constant.LF)
	}
	buf.WriteByte(constant.LF)
}
