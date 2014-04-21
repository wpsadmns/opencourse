package main
// package preprocess

import (
	"os"
	"time"
	"regexp"
	"io/ioutil"
	"io"
	"fmt"
	"bufio"
	"bytes"
	"sort"
)

const MaxLen = 2

var words sort.StringSlice

func main() {
	words = gotWords("./data/word.dict")
	words.Sort()
	uInputC := getUserInput()
	var t time.Time
	var c time.Duration
	for input := range uInputC {
		t = time.Now()
		output := SegmentSentenceMM(input)
		c = time.Now().Sub(t)
		fmt.Println(c.Nanoseconds(), c.Seconds(), output)
	}
}

func getUserInput() chan string {
	defer dealWith()
	bfReader := bufio.NewReader(os.Stdin)
	uInputC := make(chan string, 1)
	go func() {
		for {
			line, err := bfReader.ReadBytes('\n')
			check(err)
			if string(line) == "quit" {
				close(uInputC)
				break
			}
			uInputC <- string(line)
		}
	}()
	return uInputC
}

func confirmMax(max int, dest int) int {
	if max > dest {
		max = dest
	} 
	return max
}

func CutWord(s1 string) string {
	var input []rune = []rune(s1)
	var output []rune
	max := confirmMax(5, len(input))
	for len(input) >= max {
		w := input[0:max]
		for len(w) > 1 && FindWords(string(w)) == -1{
			w = w[:len(w)-1]
		}
		for _, v := range w {
			output = append(output, v)
		}
		output = append(output, '/')
		input = input[len(w):]
	}
	if 1 == 1 {
		return string(output)
	}

	var s2 string = ""
	fmt.Println(s2)
	for s1 != "" && len(s1) > MaxLen{
		w := s1[0:MaxLen]
		fmt.Println(w)
		for len(w) > 1 {
			if FindWords(w) == -1 {
				w = w[:len(w)-1]
			}
		}
		s2 = w + "/"
		s1 = s1[len(w):]
	}
	return string(output)
}

func FindWords(word string) int {
	for i, v := range words {
		if word == v {
			return i
		}
	}
	return -1
}

func FoundWords(word string) int {
	var (
		min int = 0
		max int = len(words) - 1
		mid int
	)
	for max > mid {
		mid = (max + min) >> 1
		if word == words[mid] {
			return mid
		} else if word > words[mid] {
			min = mid + 1
		} else {
			max = mid - 1
		}
	}
	return -1
}

var zhRx *regexp.Regexp = regexp.MustCompile("[\u4e00-\u9fa5]+")

func gotWords(path string) sort.StringSlice {
	defer dealWith()
	file, err := os.Open(path)
	check(err)
	data, err := ioutil.ReadAll(file)
	check(err)
	words := sort.StringSlice{}
	for _, submatch := range zhRx.FindAllSubmatch(data, -1) {
		words = append(words, string(submatch[0]))
	}
	return words
}

func dealWith() {
	if err := recover(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getWords(path string) []string {
	defer dealWith()
	file, err := os.Open(path)
	check(err)
	bfReader := bufio.NewReader(file)
	words := []string{}
	for {
		line, err := bfReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			check(err)
		}
		fmt.Println(words)
		words = append(words, string(bytes.TrimSpace(line)))
	}
	return words
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func SegmentSentenceMM(s string) string {
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
				s2 = sAppend(s2)
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
						s2 = sAppend(s2)
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

		s2 = append(s2, SegmentHzStrMM(s1)...)

		if i <= l {
			s1 = s1[i:]
		} else {
			break
		}
	}

	return string(s2)
}

func SegmentHzStrMM(s1 []rune) []rune {
	var s2 []rune

	for len(s1) != 0 {
		l := len(s1)
		if l > 4 {
			l = 4
		}
		w := s1[0:l]
		isw := FindWords(string(w))
		for l > 2 && isw == -1 {
			l -= 2
			w = w[0:l]
			isw = FindWords(string(w))
		}

		s2 = append(s2, w...)
		s2 = sAppend(s2)
		s1 = s1[len(w):]
	}

	return s2
}

func sAppend(r []rune) []rune {
	r = append(r, '/')
	r = append(r, ' ')
	r = append(r, ' ')
	return r
}
