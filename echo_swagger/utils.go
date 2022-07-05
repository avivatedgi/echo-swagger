package echo_swagger

import (
	"strings"
)

func parseStringByQuotesAndSpaces(s string) []string {
	quoted := false

	data := strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}

		return !quoted && r == ' '
	})

	for idx, val := range data {
		data[idx] = strings.Trim(val, "\"")
	}

	return data
}
