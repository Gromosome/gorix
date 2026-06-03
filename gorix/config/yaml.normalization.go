package config

import "strings"

func normalizeYAMLLines(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	rawLines := strings.Split(content, "\n")

	lines := make([]string, 0, len(rawLines))

	for _, line := range rawLines {
		line = removeYAMLComment(line)

		if strings.TrimSpace(line) == "" {
			lines = append(lines, "")
			continue
		}

		lines = append(lines, line)
	}

	return lines
}

func removeYAMLComment(line string) string {
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

		case '#':
			if !inSingleQuote && !inDoubleQuote {
				return strings.TrimRight(line[:i], " ")
			}
		}
	}

	return line
}

func countLeadingSpaces(line string) int {
	count := 0

	for _, char := range line {
		if char == ' ' {
			count++
			continue
		}

		break
	}

	return count
}
