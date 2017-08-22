package logltsv

import "bytes"

func ParseField(b []byte) []byte {
	buf := bytes.Buffer{}
	p, r := ParseTime(b)
	buf.Write(p)
	p, r = ParseFilename(r)
	buf.Write(p)
	buf.Write(r)
	return buf.Bytes()
}

func ParseTime(b []byte) (parsed, rest []byte) {
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

func ParseFilename(b []byte) (parsed, rest []byte) {
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
		buf.WriteString("\t")
		return buf.Bytes(), rest
	}
	return parsed, b
}
