package yaml

import (
	"os"
	"strconv"
	"strings"
)

func parseScalarOrInline(value string) YAMLValue {
	value = strings.TrimSpace(value)
	value = expandEnvVars(value)

	if value == "" {
		return ""
	}

	if value == "null" || value == "~" {
		return nil
	}

	if value == "true" {
		return true
	}

	if value == "false" {
		return false
	}

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return parseInlineList(value)
	}

	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return parseInlineMap(value)
	}

	if isQuoted(value) {
		return unquote(value)
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}

	return value
}

func isQuoted(value string) bool {
	if len(value) < 2 {
		return false
	}

	return (value[0] == '"' && value[len(value)-1] == '"') ||
		(value[0] == '\'' && value[len(value)-1] == '\'')
}

func unquote(value string) string {
	if !isQuoted(value) {
		return value
	}

	value = value[1 : len(value)-1]
	value = strings.ReplaceAll(value, `\"`, `"`)
	value = strings.ReplaceAll(value, `\'`, `'`)

	return value
}

func parseInlineList(value string) []YAMLValue {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")

	if strings.TrimSpace(value) == "" {
		return []YAMLValue{}
	}

	parts := splitCommaAware(value)

	result := make([]YAMLValue, 0, len(parts))

	for _, part := range parts {
		result = append(result, parseScalarOrInline(strings.TrimSpace(part)))
	}

	return result
}
func parseInlineMap(value string) map[string]YAMLValue {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "{")
	value = strings.TrimSuffix(value, "}")

	result := make(map[string]YAMLValue)

	if strings.TrimSpace(value) == "" {
		return result
	}

	parts := splitCommaAware(value)

	for _, part := range parts {
		key, val, ok := splitYAMLKeyValue(strings.TrimSpace(part))
		if !ok {
			continue
		}

		result[key] = parseScalarOrInline(val)
	}

	return result
}
func splitCommaAware(value string) []string {
	parts := make([]string, 0)

	inSingleQuote := false
	inDoubleQuote := false
	depthSquare := 0
	depthCurly := 0

	start := 0

	for i, char := range value {
		switch char {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}

		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}

		case '[':
			if !inSingleQuote && !inDoubleQuote {
				depthSquare++
			}

		case ']':
			if !inSingleQuote && !inDoubleQuote && depthSquare > 0 {
				depthSquare--
			}

		case '{':
			if !inSingleQuote && !inDoubleQuote {
				depthCurly++
			}

		case '}':
			if !inSingleQuote && !inDoubleQuote && depthCurly > 0 {
				depthCurly--
			}

		case ',':
			if !inSingleQuote && !inDoubleQuote && depthSquare == 0 && depthCurly == 0 {
				parts = append(parts, strings.TrimSpace(value[start:i]))
				start = i + 1
			}
		}
	}

	parts = append(parts, strings.TrimSpace(value[start:]))

	return parts
}
func expandEnvVars(value string) string {
	return os.ExpandEnv(value)
}
