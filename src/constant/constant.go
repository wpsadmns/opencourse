package constant

import (
	"os"
	"io/ioutil"
	"regexp"
	"bytes"
)

const (
	LF          = byte('\n')
	HeadVersion = "version:"
	HeadURL     = "url:"
	HeadDate    = "date:"
	HeadIP      = "ip:"
	HeadLength  = "length:"

	RawFilePrefix = "OpenCourse_Raw"
)

var BheadURL []byte = []byte("url:")
var BheadLength []byte = []byte("length:")


// Raw Config File
type Config struct {
	prop map[string]string
	filename string
}

var configRx *regexp.Regexp = regexp.MustCompile(`([^=]*?)=(.*)`)
var ConfProp *Config

func init() {
	ConfProp = NewConfig("raw_config.properties")
}

func NewConfig(filename string) *Config {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	check(err)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	check(err)
	config := new(Config)
	config.filename = filename
	config.prop = map[string]string{}
	for _, submatch := range configRx.FindAllSubmatch(data, -1) {
		config.prop[string(bytes.TrimSpace(submatch[1]))] = string(bytes.TrimSpace(submatch[2]))
	}
	return config
}

func (this *Config) GetProperty(key string) string {
	return this.prop[key]
}

func (this *Config) SetProperty(key, value string) {
	this.prop[key] = value
}

func (this *Config) Len() int {
	return len(this.prop)
}

func (this *Config) getKey(key string) string {
	return this.filename + "." + key
}

func (this *Config) KeySet() []string {
	keys := []string{}
	for k, _ := range this.prop {
		keys = append(keys, k)
	}
	return keys
}

func (this *Config) Truncate() {
	for k, _ := range this.prop {
		delete(this.prop, k)
	}
	this.Store()
}

func (this *Config) Store() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	file, err := os.OpenFile(this.filename, os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	err = file.Truncate(0)
	check(err)
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	for k, v := range this.prop {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte(LF)
	}
	_, err = file.Write(buf.Bytes())
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
