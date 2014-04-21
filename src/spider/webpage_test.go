package spider

import (
	"log"
	"testing"
	"util/htmlutil"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestConnectLink(t *testing.T) {
	page, err := ConnectLink("http://edu.51cto.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("------------------------------------------------------------------------")
	log.Println(string(page.HTML))
	log.Println("------------------------------------------------------------------------")
	html := htmlutil.ConvertText(page.HTML)
	log.Println(string(html))
	log.Println("------------------------------------------------------------------------")
	if 1 == 1 {
		t.Error("")
	}
}
