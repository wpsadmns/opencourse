package dbutil

import (
	"labix.org/v2/mgo"
	"logx"
)

const (
	User string = "xuchenfeng"
)

var (
	session *mgo.Session
	dbname = "xcf_opencourse"
	PageIndex = "page_index"
)

func init() {
	var err error
	session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		logx.Logger.Fatalf("[ERROR] connect mongodb failed, error is %v", err)
	}
}

func openSession() *mgo.Session {
	return session.Clone()
}

func exec(collection string, f func(*mgo.Collection)) {
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			logx.Logger.Println("[ERROR] M", err)
		}
	}()

	c := openSession().DB(dbname).C(collection)
	f(c)
}

func PreInsert() (chan interface{}, chan struct{}) {
	dataChannel := make(chan interface{}, 50)
	dbDone := make(chan struct{}, 1)
	go Insert(dataChannel, dbDone)
	return dataChannel, dbDone
}

func Insert(datas <-chan interface{}, dbDone chan<- struct{}) {
	defer func() {
		if err := recover(); err != nil {
			logx.Logger.Printf("[ERROR] insert page index failed, error is %v", err)
		}
	}()
	builder := func(c *mgo.Collection) {
		logx.Logger.Println("[INFO] Prepare receive from dataChannel...")
		var total int64
		for data := range datas {
			err := c.Insert(data)
			if err != nil {
				panic(err)
			}
			total++
		}
		logx.Logger.Printf("[INFO] Insert %d record.", total)
	}
	exec(PageIndex, builder)
	dbDone <- struct{}{}
}
