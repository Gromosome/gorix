package yaml

import (
	"strings"
)

func splitYAMLKeyValue(line string) (key string, value string, hasValue bool) {
	inSingleQuote := false
	inDoubleQuote := false

	for i, char := range line {
		switch char {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}

		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}

		case ':':
			if !inSingleQuote && !inDoubleQuote {
				key = strings.TrimSpace(line[:i])
				value = strings.TrimSpace(line[i+1:])

				if value == "" {
					return key, "", false
				}

				return key, value, true
			}
		}
	}

	return "", "", false
}
