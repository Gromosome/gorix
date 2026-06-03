package config

import "strings"

func detectNextBlockType(lines []string, start int, parentIndent int) string {
	for j := start; j < len(lines); j++ {
		line := lines[j]

		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := countLeadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		if indent <= parentIndent {
			return "map"
		}

		if strings.HasPrefix(trimmed, "- ") {
			return "list"
		}

		return "map"
	}

	return "map"
}

func parseBlockList(lines []string, index *int, parentIndent int) []YAMLValue {
	result := make([]YAMLValue, 0)

	i := *index + 1

	for i < len(lines) {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		indent := countLeadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		if indent <= parentIndent {
			break
		}

		if !strings.HasPrefix(trimmed, "- ") {
			i++
			continue
		}

		itemRaw := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))

		if itemRaw == "" {
			childMap, nextIndex := parseNestedMapForListItem(lines, i+1, indent)
			result = append(result, childMap)
			i = nextIndex
			continue
		}

		key, value, hasValue := splitYAMLKeyValue(itemRaw)
		if hasValue && key != "" {
			itemMap := make(map[string]YAMLValue)
			itemMap[key] = parseScalarOrInline(value)

			childMap, nextIndex := parseNestedMapForListItem(lines, i+1, indent)
			for k, v := range childMap {
				itemMap[k] = v
			}

			result = append(result, itemMap)
			i = nextIndex
			continue
		}

		result = append(result, parseScalarOrInline(itemRaw))
		i++
	}

	*index = i - 1

	return result
}

func parseNestedMapForListItem(lines []string, start int, listItemIndent int) (map[string]YAMLValue, int) {
	result := make(map[string]YAMLValue)

	i := start

	for i < len(lines) {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		indent := countLeadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		if indent <= listItemIndent {
			break
		}

		if strings.HasPrefix(trimmed, "- ") {
			break
		}

		key, value, hasValue := splitYAMLKeyValue(trimmed)
		if key == "" {
			i++
			continue
		}

		if hasValue {
			result[key] = parseScalarOrInline(value)
			i++
			continue
		}

		nextType := detectNextBlockType(lines, i+1, indent)
		if nextType == "list" {
			list := parseBlockList(lines, &i, indent)
			result[key] = list
			i++
			continue
		}

		nestedMap, nextIndex := parseNestedMapForListItem(lines, i+1, indent)
		result[key] = nestedMap
		i = nextIndex
	}

	return result, i
}
