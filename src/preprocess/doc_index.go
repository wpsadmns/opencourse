package preprocess

import (
	"bufio"
	"bytes"
	"constant"
	"crypto/md5"
	"fmt"
	"io"
	"logx"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"util/dbutil"
)

type IndexBuilder struct {
	root     string
	rawFiles []RawFile
	done     chan struct{}
}

func ConstructIndexBuilder(root string) *IndexBuilder {
	docIndex := &DocIndexCreater{root: root}
	docIndex.findSubRawFiles()
	docIndex.done = make(chan struct{}, len(docIndex.rawFiles))
	return docIndex
}

func (this IndexBuilder) ProcessDocument() {
	docChannel := CollectingDocument()
	urlChannel, unsortedChannel := PreInsert(docChannel)
	docIndexChannel := CreateDocIndex(urlChannel)
	go CreateURLIndex(docIndexChannel)
	forwardChannel := CreateForwardIndex(unsortedChannel)
	CreateInvertedIndex(forwardChannel)
}

func PreInsert(docChannel <-chan *Document) (<-chan *URLDocument, <-chan *UnsortedDocument) {
	urlChannel, unsortedChannel := make(chan *URLDocument, 100), make(chan *UnsortedDocument, 100)
	go func() {
		for document := range docChannel {
			urlChannel <- generateURLDocument(document)
			unsorted <- generateUnsortedDocument(document)
		}
		close(urlChannel)
		close(unsortedChannel)
	}()
	return urlChannel, unsortedChannel
}

type URLDocument struct {
}

type UnsortedDocument struct {
}

func generateURLDocument(document *Document) *URLDocument {
	return nil
}

func generateUnsortedChannel(doc *Document) *UnsortedDocument {
	return nil
}

func (this IndexBuilder) CollectingDocument() <-chan Document {
	docChannel := make(chan *Document, 1000)
	for _, rawFile := range this.rawFiles {
		go func(_rawFile RawFile) {
			BuildRawCollector(_rawFile, docChannel, this.done).Collecting()
		}(rawFile)
	}
	go func() {
		for i := 0; i < cap(this.done); i++ {
			<-this.done
		}
		close(docChannel)
	}()
	return docChannel
}

/*
 *
 */
func CreateDocIndex() {
}

/*
 *
 */
func CreateURLIndex() {
}

/*
 *
 */
func CreateForwardIndex() {
}

/*
 *
 */
func CreateInvertedIndex() {
}

func (this DocIndexCreater) Wait2Process(f func()) {
	for i := 0; i < cap(this.done); i++ {
		<-this.done
	}
	f()
}

type RawFile struct {
	path   string
	suffix string
}

func (this *IndexBuilder) findSubRawFiles() {
	findRawFile := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logx.Logger.Printf("[ERROR] from %s find raw file failed, error is %s.", root, err.Error())
			return err
		}
		if !info.IsDir() &&
			strings.HasPrefix(path, constant.RawFilePrefix) {
			this.rawFiles = append(this.rawFiles, RawFile{path, path[strings.LastIndex(".")+1:]})
		}
		return nil
	}
	filepath.Walk(this.root, findRawFile)
}

type Document struct {
	URL      string
	Offset   int64
	LineLen  int
	Content  []byte
	Abstract string
	RawName  string
	RawDocId string
	DocId    string
}

type RawCollector struct {
	rawFile    RawFile
	done       chan<- struct{}
	docChannel chan<- *Document
	rawChannel chan []byte
	offset     int64
}

func BuildRawCollector(rawFile RawFile, docChannel chan<- *Document, done chan<- struct{}) *Raw {
	return &Raw{rawFile: rawFile, 
					docChannel: docChannel, 
					rawChannel: make(chan []byte, 100), 
					done: done}
}

const LF = byte('\n')

func (this *RawCollector) Collecting() {
	defer func() {
		if err := recover(); err != nil && err != io.EOF {
			logx.Logger.Fatalf("[ERROR] analysis raw file %s failed, error is %v.", this.rawpath, err)
		}
	}()

	go func() {
		rawFile, err := os.Open(this.rawpath)
		check(err)
		bufReader := bufio.NewReader(rawFile)
		for {
			line, err := bufReader.ReadBytes(LF)
			if err != nil {
				close(this.rawChannel)
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if len(line) == 1 && line[0] == LF {
				this.rawChannel <- nil
			} else {
				this.rawChannel <- line
			}
		}
	}()

	fileInfo, err := os.Stat(this.rawpath)
	check(err)
	offsetLen := this.getOffsetLen(fileInfo.Size())

	for {
		doc := new(Document)
		doc.LineLen = offsetLen
		doc.RawDocId = this.rawFile.suffix
		over := this.readRawContent(doc, this.readRawHead(doc))
		this.docChannel <- doc
		if over {
			break
		}
	}
	this.done <- struct{}{}
}

func (this *RawCollector) getOffsetLen(n int64) (l int) {
	if n == 0 {
		return 1
	}
	for n > 0 {
		n = n / 10
		l++
	}
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *RawCollector) readRawHead(doc *Document) (buffLen int) {
	doc.Offset = this.offset
	doc.RawName = this.rawFile.path

	for {
		if line, ok := <-this.rawChannel; ok {
			if line == nil {
				this.offset++
				break
			} else {
				// ignored except URL and length
				if theLine := bytes.TrimSpace(line); bytes.HasPrefix(theLine, constant.BheadURL) {
					doc.URL = string(bytes.FieldsFunc(theLine, splitHead)[1])
				} else if bytes.HasPrefix(theLine, constant.BheadLength) {
					buffLen, _ = strconv.Atoi(string(bytes.FieldsFunc(theLine, splitHead)[1]))
				}
				this.offset += int64(len(line))
			}
		} else {
			logx.Logger.Printf("read head info from DocIndex.rawChannel failed, channel closed.")
			break
		}
	}
	return
}

func (this *RawCollector) readRawContent(doc *Document, buffLen int) (over bool) {
	buff := bytes.NewBuffer(make([]byte, 0, buffLen))
	for {
		if line, ok := <-this.rawChannel; ok {
			if line == nil {
				this.offset++
				break
			}
			buff.Write(line)
			this.offset += int64(len(line))
		} else {
			over = true
			logx.Logger.Printf("read content from DocIndex.rawChannel failed, channel closed, offset is %d.", this.offset)
			break
		}
	}
	doc.Abstract = ToMD5(buff.Bytes())
	doc.Content = buff.Bytes()
	return
}

func ToMD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%02x", h.Sum(nil))
}

func splitHead(r rune) bool {
	return r == ':'
}
