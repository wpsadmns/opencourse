package preprocess

import (
	"testing"
)

func TestSegmentSentence(t *testing.T) {
	ws := map[string]string{
		"Java C":    "Java C/  ",
		"Love清华大学":  "Love/  清华大学/  ",
		"南京":        "南京/  ",
		"南京java":    "南京/  java/  ",
		"C++南京java": "C++/  南京/  java/  ",
	}

	for k, v := range ws {
		if theV := segment.SegmentSentence(k); v != theV {
			t.Errorf("translate |%s| is |%s|, not |%s|.", k, theV, v)
		}
	}
}
