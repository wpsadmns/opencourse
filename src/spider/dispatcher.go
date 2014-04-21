/**
 *
 */
package spider

import (
	"time"
)

var workers int

type Dispatcher struct {
	urls ConcurrentURL
}

// Build New Dispatcher
func NewDispatcher(urls []string, theWorkers int) *Dispatcher {
	workers = theWorkers
	concurrentUrl := NewConcurrentURL()
	d := &Dispatcher{urls: concurrentUrl}
	for _, url := range urls {
		d.Push(url)
	}
	return d
}

func (this *Dispatcher) Push(url string) {
	this.urls.Add(url)
}

func (this *Dispatcher) PushAll(urls map[string]bool) {
	for url, _ := range urls {
		this.urls.Add(url)
	}
}

func (this *Dispatcher) Pop() string {
	return this.urls.Get()
}

func (this *Dispatcher) Destory() map[string]bool {
	return this.urls.Over()
}

func (this *Dispatcher) Len() int {
	return this.urls.Len()
}

type ConcurrentURL interface {
	Add(string)
	Get() string
	Over() map[string]bool
	Len() int
}

type commandStruct struct {
	action commandAction
	value  string
	result chan interface{}
}

type commandAction int

const (
	get = iota
	add
	length
	over
)

type concurrentURL chan commandStruct

func (c concurrentURL) Get() string {
	reply := make(chan interface{})
	c <- commandStruct{action: get, result: reply}
	return (<-reply).(string)
}

func (c concurrentURL) Add(url string) {
	c <- commandStruct{action: add, value: url}
}

func (c concurrentURL) Len() int {
	reply := make(chan interface{})
	c <- commandStruct{action: length, result: reply}
	return (<-reply).(int)
}

func (c concurrentURL) Over() map[string]bool {
	reply := make(chan interface{})
	c <- commandStruct{action: over, result: reply}
	return (<-reply).(map[string]bool)
}

func (c concurrentURL) doTask() {
	store := make(map[string]bool)
	// done := make(chan struct{}, workers)
	// if after 5 mins not receive any url
	var timeout time.Duration = 1e9 * 60 * 5
	for cmd := range c {
		switch cmd.action {
		case get:
			timer := time.After(timeout)
			isContinue := true
			for isContinue {
				select {
				case <-timer:
					cmd.result <- ""
					// reset 30 secs
					timeout = 1e9 * 30
					isContinue = false
				default:
					var url string
					var finish bool
					for url, finish = range store {
						if !finish {
							store[url] = true
							cmd.result <- url
							isContinue = false
							break
						}
					}
				}
			}
		case add:
			// the url not exists then add it.
			if _, ok := store[cmd.value]; !ok {
				store[cmd.value] = false
			}
		case length:
			cmd.result <- len(store)
		case over:
			// close concurrentURL for return urls.
			close(c)
			cmd.result <- store
		}
	}
}

func complete(done chan<- struct{}) {
	for i := 0; i < workers; i++ {
		done <- struct{}{}
	}
}

func NewConcurrentURL() ConcurrentURL {
	concurrent := make(concurrentURL, 1)
	go concurrent.doTask()
	return concurrent
}
