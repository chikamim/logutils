package logutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	Level   LogLevel
	Time    time.Time
	Message string
	Caller  string
	Data    string
	Empty   bool
}

func NewEntry(line []byte) *Entry {
	entry := &Entry{}
	entry.Level, line = logLevel(line)
	entry.Time = time.Now()
	if len(line) == 0 {
		entry.Empty = true
	}
	_, file, lno, _ := runtime.Caller(4)
	entry.Caller = fmt.Sprintf("%v:%v", file, lno)
	m := bytes.SplitN(line, []byte("\t"), 2)
	if len(m) < 2 {
		entry.Message = string(line)
		return entry
	}
	entry.Message = string(m[0])
	entry.Data = string(m[1])
	return entry
}

func (e *Entry) JSON() []byte {
	if e.Empty {
		return []byte{}
	}
	buf := bytes.Buffer{}
	buf.WriteString("level:")
	buf.WriteString(strings.ToLower(string(e.Level)))
	buf.WriteString("\tts:")
	buf.WriteString(epochTime(e.Time))
	buf.WriteString("\tcaller:")
	buf.WriteString(trimPath(e.Caller))
	buf.WriteString("\tmsg:")
	buf.WriteString(string(e.Message))
	if len(e.Data) > 0 {
		buf.WriteString("\t")
		buf.WriteString(e.Data)
	}
	buf.WriteString("\n")
	return []byte(ltsvToJSON(buf.String()))
}

func (e *Entry) Text() []byte {
	if e.Empty {
		return []byte{}
	}
	buf := bytes.Buffer{}
	buf.WriteString(e.Time.Format("2006-01-02T15:04:05.000Z0700 "))
	buf.WriteString("\t")
	buf.WriteString(string(e.Level))
	buf.WriteString("\t")
	buf.WriteString(trimPath(e.Caller))
	buf.WriteString("\t")
	buf.WriteString(e.Message)
	if len(e.Data) > 0 {
		buf.WriteString("\t")
		buf.WriteString(ltsvToJSON(e.Data))
	}
	return buf.Bytes()
}

func ltsvToJSON(s string) string {
	j := bytes.Buffer{}
	j.WriteString("{")
	ln := strings.HasSuffix(s, "\n")
	if ln {
		s = strings.TrimRight(s, "\n")
	}
	pairs := strings.Split(s, "\t")
	for i, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) < 2 {
			j.WriteString(escape(kv[0]) + ":\"\"")
			continue
		}
		if isNumber(kv[1]) {
			j.WriteString(escape(kv[0]) + ":" + kv[1])
		} else {
			j.WriteString(escape(kv[0]) + ":" + escape(kv[1]))
		}
		if i < len(pairs)-1 {
			j.WriteString(",")
		}
	}
	j.WriteString("}")
	if ln {
		j.WriteString("\n")
	}
	return j.String()
}

func escape(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return s
	}
	return string(b)
}

func isNumber(s string) bool {
	s = strings.Trim(s, " ")
	if len(s) == 0 || s == "." {
		return false
	}
	dot := false
	for i, c := range s {
		var a, b byte
		if i > 0 {
			b = s[i-1]
		}
		if i < len(s)-1 {
			a = s[i+1]
		}
		if c == '.' {
			if dot || ((b < '0' || b > '9') && (a < '0' || a > '9')) {
				return false
			}
			dot = true
		} else if c >= '0' && c <= '9' {
			continue
		} else {
			return false
		}
	}
	return true
}

func trimPath(s string) string {
	idx := strings.LastIndexByte(s, '/')
	if idx == -1 {
		return s
	}
	idx = strings.LastIndexByte(s[:idx], '/')
	if idx == -1 {
		return s
	}
	return s[idx+1:]
}

func epochTime(t time.Time) string {
	sec := float64(t.UnixNano()) / float64(time.Second)
	return fmt.Sprintf("%.7f", sec)
}
