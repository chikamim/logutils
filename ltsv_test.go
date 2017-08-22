package logltsv

import "testing"

func TestParseField(t *testing.T) {
	text := "2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message:ok"
	got := ParseField([]byte(text))
	if string(got) != "time:2009/01/23 01:23:23.123123\tfile:/a/b/c/d.go\tline:23\tmessage:ok" {
		t.Errorf("Fail: %s", got)
	}
}

func TestParseTime_long(t *testing.T) {
	text := "2009/01/23 01:23:23 message:ok"
	parsed, rest := ParseTime([]byte(text))
	if string(parsed) != "time:2009/01/23 01:23:23\t" {
		t.Errorf("Fail: %s", parsed)
	}
	if string(rest) != "message:ok" {
		t.Errorf("Fail: %s", rest)
	}
}

func TestParseTime_long_micro(t *testing.T) {
	text := "2009/01/23 01:23:23.123123 message:ok"
	parsed, rest := ParseTime([]byte(text))
	if string(parsed) != "time:2009/01/23 01:23:23.123123\t" {
		t.Errorf("Fail: %s", parsed)
	}
	if string(rest) != "message:ok" {
		t.Errorf("Fail: %s", rest)
	}
}

func TestParseFilename_long(t *testing.T) {
	text := "/a/b/c/d.go:23: message:ok"

	parsed, rest := ParseFilename([]byte(text))
	if string(parsed) != "file:/a/b/c/d.go\tline:23\t" {
		t.Errorf("Fail: %s", parsed)
	}
	if string(rest) != "message:ok" {
		t.Errorf("Fail: %s", rest)
	}
}

func TestParseFilename_short(t *testing.T) {
	text := "d.go:23: message:ok"

	parsed, rest := ParseFilename([]byte(text))
	if string(parsed) != "file:d.go\tline:23\t" {
		t.Errorf("Fail: %s", parsed)
	}
	if string(rest) != "message:ok" {
		t.Errorf("Fail: %s", rest)
	}
}
