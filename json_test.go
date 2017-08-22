package logltsv

import "testing"

func TestToJSON(t *testing.T) {
	text := "time:2009/01/23 01:23:23.123123\tfile:/a/b/c/d.go\tline:23\tmessage:ok\n"
	got := ToJSON([]byte(text))
	if string(got) != `{"time":"2009/01/23 01:23:23.123123", "file":"/a/b/c/d.go", "line":"23", "message":"ok"}
` {
		t.Errorf("Fail: %s", got)
	}
}
