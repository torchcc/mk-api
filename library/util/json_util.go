package util

import (
	"encoding/json"
	"strings"
)

func JsonDump(v interface{}, indent int) string {
	b, err := json.MarshalIndent(v, " ", strings.Repeat(" ", indent))
	if err != nil {
		return ""
	}
	return string(b)
}
