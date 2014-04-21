/*
 *
 */
package spider

import (
	"logx"
	"os"
	"util/stringutil"
)

type Gather struct {
	dispatcher *Dispatcher
	offset     int64
	Id         string
	done       chan<- struct{}
	RawFile    *os.File
}

/* build new gather */
func NewGather(dispatcher *Dispatcher, id string, done chan<- struct{}) *Gather {
	filepath := "./raws/OpenCourse_Raw." + id
	rawFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logx.Logger.Fatalf("[ERROR] open file %s failed, error is %s.", filepath, err.Error())
	}
	return &Gather{dispatcher: dispatcher, Id: id, RawFile: rawFile, done: done}
}

func (this *Gather) Do() {
	pop := this.dispatcher.Pop
	defer this.RawFile.Close()
	var (
		page  *Page
		err   error
		links []string
	)
	for url := pop(); url != ""; url = pop() {
		page, err = PageFromURL(url)
		if err != nil {
			logx.Logger.Println(err)
			continue
		}
		// analysis page to decode url and save document.
		links = doAnalyzer(this.RawFile, page)
		for _, link := range links {
			if stringutil.IsNotEmpty(link) {
				this.dispatcher.Push(link)
			}
		}
	}
	this.done <- struct{}{}
	logx.Logger.Printf("[INFO] a goroutinue disposed...\n")
}
