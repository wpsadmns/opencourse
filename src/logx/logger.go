package logx

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	logfile, err := os.OpenFile("./opencourse.log", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("Open spider.log s failed, the error is %s\n", err.Error())
	}
	if 1 == 3 {
		Logger = log.New(logfile, "", log.Lshortfile|log.LstdFlags)
	} else {
		Logger = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	}
}
