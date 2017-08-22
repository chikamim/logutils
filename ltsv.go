package logltsv

import "bytes"

func parseField(b []byte) []byte {
	buf := bytes.Buffer{}
	p, r := parseTime(b)
	buf.Write(p)
	p, r = parseFilename(r)
	buf.Write(p)
	buf.WriteString("\t")
	buf.Write(r)
	return buf.Bytes()
}

func parseTime(b []byte) (parsed, rest []byte) {
	if b[4] == byte('/') && b[16] == byte(':') && b[19] == byte('.') {
		time := b[0:26]
		rest := b[27:]
		buf := bytes.Buffer{}
		buf.WriteString("time:")
		buf.Write(time)
		buf.WriteString("\t")
		return buf.Bytes(), rest
	} else if b[4] == byte('/') && b[16] == byte(':') {
		time := b[0:19]
		rest := b[20:]
		buf := bytes.Buffer{}
		buf.WriteString("time:")
		buf.Write(time)
		buf.WriteString("\t")
		return buf.Bytes(), rest
	}
	return parsed, b
}

func parseFilename(b []byte) (parsed, rest []byte) {
	e := bytes.Index(b, []byte(".go:"))
	l := bytes.Index(b, []byte(": "))
	if e > 0 && l > 0 {
		file := b[0 : e+3]
		line := b[e+4 : l]
		rest := b[l+2:]
		buf := bytes.Buffer{}
		buf.WriteString("file:")
		buf.Write(file)
		buf.WriteString("\tline:")
		buf.Write(line)
		return buf.Bytes(), rest
	}
	return parsed, b
}
