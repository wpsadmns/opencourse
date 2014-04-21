package preprocess

import (
	"io/ioutil"
	"logx"
	"os"
	"regexp"
	"strconv"
)

var segment *DocSegment

func init() {
	segment = new(DocSegment)
	segment.dictMap = LoadDict("../data/word.dict")
}

func LoadDict(dictPath string) map[string]uint {
	file, err := os.Open(dictPath)
	if err != nil {
		logx.Logger.Fatalf("[ERROR] open dict file %s failed, error is %s.", dictPath, err.Error())
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logx.Logger.Fatalf("[ERROR] read dict file %s failed, error is %s.", dictPath, err.Error())
	}
	dictMap := map[string]uint{}
	pattern := "([\u4e00-\u9fa5]+)\\s*(\\d+)"
	for _, submatch := range regexp.MustCompile(pattern).FindAllSubmatch(data, -1) {
		n, err := strconv.Atoi(string(submatch[2]))
		if err != nil {
			logx.Logger.Printf("[ERROR] %s has error number %s", string(submatch[1]), string(submatch[2]))
			n = 0
		}
		dictMap[string(submatch[1])] = uint(n)
	}
	return dictMap
}

type DocSegment struct {
	dictMap map[string]uint
}

func (this *DocSegment) HasWords(s string) (ok bool) {
	_, ok = this.dictMap[s]
	return
}

func (this *DocSegment) SegmentSentence(s string) string {
	s1 := []rune(s)
	var s2 []rune
	var i, l int

	for len(s1) != 0 {
		ch := s1[0]
		if ch < 128 {
			i = 1
			l = len(s1)

			for i < l && s1[i] < 128 && s1[i] != 10 && s1[i] != 13 {
				i++
			}

			if ch != 32 && ch != 10 && ch != 13 {
				s2 = append(s2, s1[:i]...)
				s2 = this.sAppend(s2)
			} else {
				if ch == 10 || ch == 13 {
					s2 = append(s2, s1[:i]...)
				}
			}

			if i < len(s1) {
				s1 = s1[i:]
			} else {
				break
			}
			continue
		} else {
			if ch < 176 {
				i = 0
				l = len(s1)

				for i < l && s1[i] < 176 && s1[i] >= 161 && !(s1[i] == 161 && s1[i+1] >= 162 && s1[i+1] <= 168) && !(s1[i] == 161 && s1[i+1] >= 171 && s1[i+1] <= 191) && (!(s1[i] == 163 && (s1[i+1] == 172 || s1[i+1] == 161) || s1[i+1] == 168 || s1[i+1] <= 169 || s1[i+1] == 186 || s1[i+1] == 187 || s1[i+1] == 191)) {
					i += 2
				}

				if i == 0 {
					i += 2
				}

				if !(ch == 161 && s1[1] == 161) {
					if i < len(s1) {
						s2 = append(s2, s1[0:i]...)
						s2 = this.sAppend(s2)
					} else {
						break
					}
				}

				if i <= len(s1) {
					s1 = s1[i:]
				} else {
					break
				}
				continue
			}
		}

		i, l = 2, len(s1)
		for i < l && s1[i] >= 176 {
			i += 2
		}

		s2 = append(s2, this.SegmentZh(s1)...)

		if i <= l {
			s1 = s1[i:]
		} else {
			break
		}
	}

	return string(s2)
}

func (this *DocSegment) SegmentZh(s1 []rune) []rune {
	var s2 []rune

	for len(s1) != 0 {
		l := len(s1)
		if l > 4 {
			l = 4
		}
		w := s1[0:l]
		isw := this.HasWords(string(w))
		for l > 2 && !isw {
			l -= 2
			w = w[0:l]
			isw = this.HasWords(string(w))
		}

		s2 = append(s2, w...)
		s2 = this.sAppend(s2)
		s1 = s1[len(w):]
	}

	return s2
}

func (*DocSegment) sAppend(r []rune) []rune {
	r = append(r, '/')
	r = append(r, ' ')
	r = append(r, ' ')
	return r
}
