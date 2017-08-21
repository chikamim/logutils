// Package logutils augments the standard log package with levels.
package logutils

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger that
// will filter out log messages that aren't at least a certain level.
//
// Once the filter is in use somewhere, it is not safe to modify
// the structure.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"DEBUG", "WARN", "ERROR"}.
	Levels []LogLevel

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	f.once.Do(f.init)

	// Check for a log level
	prefix := "level:"
	var level LogLevel
	x := bytes.Index(line, []byte(prefix))
	if x >= 0 {
		y := bytes.IndexByte(line[x:], '\t')
		if y >= 0 {
			level = LogLevel(line[x+len(prefix) : x+y])
		}
	}

	_, ok := f.badLevels[level]
	return !ok
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line. We use that as a slight optimization within
	// this method, assuming we're dealing with a single, complete line
	// of log data.

	if !f.Check(p) {
		return len(p), nil
	}
	p = ToTLSVDate(p)

	return f.Writer.Write(ToJSON(p))
}

func ToTLSVDate(b []byte) []byte {
	if b[4] == byte('/') && b[16] == byte(':') {
		t := b[0:19]
		p := b[20:]
		p = append([]byte("\t"), p...)
		p = append(t, p...)
		p = append([]byte("time:"), p...)
		return p
	}
	return b
}

func ToInfluxLine(b []byte) []byte {
	// ltsv
	return b
}

func ToJSON(b []byte) []byte {
	j := bytes.Buffer{}
	j.WriteString("{")
	s := string(b[0 : len(b)-1])
	pairs := strings.Split(s, "\t")
	for i, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		j.WriteString("\"" + kv[0] + "\":\"" + kv[1] + "\"")
		if i < len(pairs)-1 {
			j.WriteString(", ")
		}
	}
	j.WriteString("}\n")
	return j.Bytes()
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	badLevels := make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if level == f.MinLevel {
			break
		}
		badLevels[level] = struct{}{}
	}
	f.badLevels = badLevels
}
