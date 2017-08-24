package logltsv

import (
	"bytes"
	"strings"
)

func ToJSON(b []byte) []byte {
	j := bytes.Buffer{}
	j.WriteString("{")
	s := strings.TrimRight(string(b), "\n")
	pairs := strings.Split(s, "\t")
	for i, pair := range pairs {
		kv := strings.SplitN(escape(pair), ":", 2)
		if len(kv) < 2 {
			//TODO: parse error. better set default key "payload"?
			return b
		}
		j.WriteString("\"" + kv[0] + "\":\"" + kv[1] + "\"")
		if i < len(pairs)-1 {
			j.WriteString(", ")
		}
	}
	j.WriteString("}\n") // TODO: check log.Print()
	return j.Bytes()
}

func escape(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	return s
}
