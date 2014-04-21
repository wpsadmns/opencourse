package spider

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

var goroutine int = 6

var urls []string = []string{
	"http://www.google.com",
	"http://www.baidu.com",
	"http://www.xuchenfeng.con",
	"http://www.google.com",
	"http://www.apple.com",
	"http://www.ibm.com",
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().Unix())
}

func TestDispatch(t *testing.T) {
	d := NewDispatcher(urls, goroutine)
	if d.Len() != len(urls)-1 {
		t.Errorf("expect urls.length = concurrent.length, but not equals")
	}
	var expected, found string
	if expected, found = urls[0], d.Pop(); expected != found {
		t.Errorf("expect %s, but found %s.", expected, found)
	}

	if expected, found = urls[1], d.Pop(); expected != found {
		t.Errorf("expect %s, but found %s.", expected, found)
	}

	if expected, found = urls[2], d.Pop(); expected != found {
		t.Errorf("expect %s, but found %s.", expected, found)
	}

	if expected, found = urls[4], d.Pop(); expected != found {
		t.Errorf("expect %s, but found %s.", expected, found)
	}

	if expected, found = urls[5], d.Pop(); expected != found {
		t.Errorf("expect %s, but found %s.", expected, found)
	}
}

func TestModified(t *testing.T) {
	d := NewDispatcher(nil, goroutine)
	total := 0
	done := make(chan struct{}, goroutine)
	pushSum := goroutine
	// Push
	for i := 0; i < pushSum; i++ {
		go func() {
			sum := rand.Intn(1000)
			total += sum
			for i := 0; i < sum; i++ {
				d.Push(fmt.Sprintf("http://www.%d%d.com", rand.Int(), rand.Int()))
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < pushSum; i++ {
		<-done
	}

	if d.Len() != total {
		t.Errorf("expected %d, but found %d\n", total, d.Len())
	}

	// Pop
	popSum := goroutine / 3
	popTotal := 0
	for i := 0; i < popSum; i++ {
		go func() {
			sum := rand.Intn(1000)
			popTotal += sum
			total -= sum
			for i := 0; i < sum; i++ {
				d.Pop()
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < popSum; i++ {
		<-done
	}

	if d.Len()-popTotal != total {
		t.Errorf("expected %d, but found %d\n", total, d.Len())
	}
}
