package mongo

import "strings"

func normalizeDatabaseName(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "default"
	}

	return value
}

func normalizeCollectionName(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "documents"
	}

	var builder strings.Builder

	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z':
			builder.WriteRune(character)

		case character >= 'A' && character <= 'Z':
			builder.WriteRune(character)

		case character >= '0' && character <= '9':
			builder.WriteRune(character)

		case character == '_' || character == '-':
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
		return "documents"
	}

	return result
}
