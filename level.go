package logutils

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel string

type Encoder int

const (
	// Plain output standard log format
	Plain Encoder = iota
	// ZapDevelop output zap develop text format
	ZapDevelop
	// ZapProduction output zap production json format
	ZapProduction
)

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

	// Encoder is the log output format
	Encoder Encoder

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

// NewFilter builds a development Logger that writes DebugLevel and above logs to standard error in a standart log format.
func NewFilter() *LevelFilter {
	return &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: LogLevel("DEBUG"),
		Writer:   os.Stderr,
		Encoder:  Plain,
	}
}

// NewDevelopment builds a development Logger that writes DebugLevel and above logs to standard error in a human-friendly format.
func NewDevelopment() *LevelFilter {
	log.SetFlags(0)
	return &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: LogLevel("DEBUG"),
		Writer:   os.Stderr,
		Encoder:  ZapDevelop,
	}
}

// NewProduction builds a sensible production Logger that writes InfoLevel and above logs to standard error as JSON.
func NewProduction() *LevelFilter {
	log.SetFlags(0)
	return &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: LogLevel("INFO"),
		Writer:   os.Stderr,
		Encoder:  ZapProduction,
	}
}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	f.once.Do(f.init)
	level, _ := logLevel(line)
	_, ok := f.badLevels[level]
	return !ok
}

func logLevel(line []byte) (level LogLevel, rest []byte) {
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = LogLevel(line[x+1 : x+y])
			rest = line[x+y+2:]
		}
	}
	return
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

	e := NewEntry(p)
	if f.Encoder == ZapProduction {
		p = e.JSON()
	} else if f.Encoder == ZapDevelop {
		p = e.Text()
	}

	return f.Writer.Write(p)
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
