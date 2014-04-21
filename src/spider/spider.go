package spider

import (
	"logx"
	"strconv"
)

type Spider struct {
	urls       []string
	done       chan struct{}
	dispatcher *Dispatcher
}

func BuildSpider(urls []string, num int) *Spider {
	return &Spider{urls: urls, done: make(chan struct{}, num)}
}

func (s *Spider) Start() {
	s.dispatcher = NewDispatcher(s.urls, s.Cap())
	for i := 0; i < s.Cap(); i++ {
		logx.Logger.Println("new gather...")
		go NewGather(s.dispatcher, strconv.Itoa(i+1), s.done).Do()
	}
}

func (s *Spider) Wait2Process(f func()) {
	logx.Logger.Println("wait function...")
	for i := 0; i < s.Cap(); i++ {
		<-s.done
	}
	for url, finish := range s.dispatcher.Destory() {
		logx.Logger.Println(url, "->", finish)
	}
	f()
}

func (s *Spider) Cap() int {
	return cap(s.done)
}
