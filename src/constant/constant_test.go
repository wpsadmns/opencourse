package constant

import (
	"testing"
)

func TestConfig(t *testing.T) {
	filename := "./OpenCourse_Raw.123892"
	config, err := NewConfig(filename)
	if err != nil {
		t.Errorf("open config file %s failed, error is %v.\n", filename, err)
	}
	config.SetProperty("aaa", "bbb")
	if value := config.GetProperty("bbb"); value != "aaa" {
		t.Errorf("got: %s, want: %s\n", value, "aaa")
	}
	config.SetProperty("bbb", "aaa")
	if value := config.GetProperty("aaa"); value != "bbb" {
		t.Errorf("got: %s, want: %s\n", value, "bbb")
	}
	err = config.Store()
	if err != nil {
		t.Errorf("Write config file %s failed, error is %v.\n", filename, err)
	}
}
