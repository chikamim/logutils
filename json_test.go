package logltsv

import "testing"

func TestToJSON(t *testing.T) {
	const want = `{"time":"2009/01/23 01:23:23.123123", "file":"/a/b/c/d.go", "line":"23", "message":"ok"}` + "\n"
	got := ToJSON([]byte("time:2009/01/23 01:23:23.123123\tfile:/a/b/c/d.go\tline:23\tmessage:ok\n"))
	if string(got) != want {
		t.Errorf("Fail: %s", got)
	}
}

func TestToJSONWithEscape(t *testing.T) {
	const want = `{"time":"2009/01/23 01:23:23.123123", "file":"/a/b/c/d.go", "line":"23", "message":"\"ok\\\""}` + "\n"
	got := ToJSON([]byte("time:2009/01/23 01:23:23.123123\tfile:/a/b/c/d.go\tline:23\tmessage:\"ok\\\"\n"))
	if string(got) != want {
		t.Errorf("Fail: %s", got)
	}
}

func TestEscapeDoubleQuotes(t *testing.T) {
	const want = `\"\"`
	got := escape(`""`)
	if got != want {
		t.Errorf("got %v; want %s", got, want)
	}
}

func TestEscapeBackSlashAndDoubleQuotes(t *testing.T) {
	const want = `\\\"\"`
	got := escape(`\""`)
	if got != want {
		t.Errorf("got %v; want %s", got, want)
	}
}
