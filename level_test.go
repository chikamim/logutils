package logltsv

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestLevelFilter_impl(t *testing.T) {
	var _ io.Writer = new(Output)
}

func TestLevelFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &Output{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   buf,
	}

	logger := log.New(filter, "", 0)
	logger.Print("level:WARN\tfoo")
	logger.Println("level:ERROR\tbar")
	logger.Println("level:DEBUG\tbaz")
	logger.Println("level:WARN\tbuzz")

	result := buf.String()
	expected := "level:WARN\tfoo\nlevel:ERROR\tbar\nlevel:WARN\tbuzz\n"

	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
}

func TestLevelFilterCheck(t *testing.T) {
	filter := &Output{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   nil,
	}

	testCases := []struct {
		line  string
		check bool
	}{
		{"level:WARN\tfoo\n", true},
		{"level:ERROR\tbar\n", true},
		{"level:DEBUG\tbaz\n", false},
		{"level:WARN\tbuzz\n", true},
	}

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.check {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}

func TestLevelFilter_SetMinLevel(t *testing.T) {
	filter := &Output{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "ERROR",
		Writer:   nil,
	}

	testCases := []struct {
		line        string
		checkBefore bool
		checkAfter  bool
	}{
		{"level:WARN\tfoo\n", false, true},
		{"level:ERROR\tbar\n", true, true},
		{"level:DEBUG\t baz\n", false, false},
		{"level:WARN\tbuzz\n", false, true},
	}

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.checkBefore {
			t.Errorf("Fail: %s", testCase.line)
		}
	}

	// Update the minimum level to WARN
	filter.SetMinLevel("WARN")

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.checkAfter {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}
