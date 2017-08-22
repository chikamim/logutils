package logltsv

import (
	"bytes"
	"strings"
)

func ToJSON(b []byte) []byte {
	j := bytes.Buffer{}
	j.WriteString("{")
	s := string(b[0 : len(b)-1])
	pairs := strings.Split(string(s), "\t")
	for i, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) < 2 {
			//TODO: parse error, set default key ("message")?
			return b
		}
		j.WriteString("\"" + kv[0] + "\":\"" + kv[1] + "\"")
		if i < len(pairs)-1 {
			j.WriteString(", ")
		}
	}
	j.WriteString("}\n")
	return j.Bytes()
}
