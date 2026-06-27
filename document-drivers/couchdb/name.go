package couchdb

import "strings"

func normalizeDatabaseName(value string) string {
	value = strings.ToLower(
		strings.TrimSpace(value),
	)

	if value == "" {
		return "default"
	}

	var builder strings.Builder

	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z':
			builder.WriteRune(character)

		case character >= '0' && character <= '9':
			builder.WriteRune(character)

		case character == '_':
			builder.WriteRune(character)

		default:
			builder.WriteByte('_')
		}
	}

	result := strings.Trim(
		builder.String(),
		"_",
	)

	if result == "" {
		return "default"
	}

	return result
}
