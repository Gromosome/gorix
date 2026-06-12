package yaml

import (
	"strconv"
	"strings"
)

type YAMLValue = any

func ParseYAML(content string) map[string]YAMLValue {
	root := make(map[string]YAMLValue)

	lines := NormalizeYAMLLines(content)

	type stackItem struct {
		indent int
		key    string
		obj    map[string]YAMLValue
	}

	stack := []stackItem{
		{
			indent: -1,
			key:    "",
			obj:    root,
		},
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := countLeadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		for len(stack) > 1 && indent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}

		parent := stack[len(stack)-1].obj

		if strings.HasPrefix(trimmed, "- ") {
			// List item under current parent key is handled when parent created list.
			continue
		}

		key, value, hasValue := SplitYAMLKeyValue(trimmed)
		if key == "" {
			continue
		}

		if !hasValue {
			nextValue := detectNextBlockType(lines, i+1, indent)

			if nextValue == "list" {
				list := parseBlockList(lines, &i, indent)
				parent[key] = list
				continue
			}

			child := make(map[string]YAMLValue)
			parent[key] = child

			stack = append(stack, stackItem{
				indent: indent,
				key:    key,
				obj:    child,
			})

			continue
		}

		parent[key] = ParseScalarOrInline(value)
	}

	return root
}

func GetMap(data map[string]YAMLValue, path string) (map[string]YAMLValue, bool) {
	value, ok := Get(data, path)
	if !ok {
		return nil, false
	}

	m, ok := value.(map[string]YAMLValue)
	return m, ok
}

func GetString(data map[string]YAMLValue, path string, fallback string) string {
	value, ok := Get(data, path)
	if !ok || value == nil {
		return fallback
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fallback
	}
}

func GetInt(data map[string]YAMLValue, path string, fallback int) int {
	value, ok := Get(data, path)
	if !ok || value == nil {
		return fallback
	}

	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed
		}
	}

	return fallback
}

func GetBoolPtr(data map[string]YAMLValue, path string) *bool {
	value, ok := Get(data, path)
	if !ok || value == nil {
		return nil
	}

	switch v := value.(type) {
	case bool:
		return &v
	case string:
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

func GetStringSlice(data map[string]YAMLValue, path string) []string {
	value, ok := Get(data, path)
	if !ok || value == nil {
		return nil
	}

	list, ok := value.([]YAMLValue)
	if !ok {
		return nil
	}

	result := make([]string, 0, len(list))

	for _, item := range list {
		switch v := item.(type) {
		case string:
			result = append(result, v)
		case int:
			result = append(result, strconv.Itoa(v))
		case float64:
			result = append(result, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			result = append(result, strconv.FormatBool(v))
		}
	}

	return result
}

func Get(data map[string]YAMLValue, path string) (YAMLValue, bool) {
	parts := strings.Split(path, ".")

	var current YAMLValue = data

	for _, part := range parts {
		currentMap, ok := current.(map[string]YAMLValue)
		if !ok {
			return nil, false
		}

		value, ok := currentMap[part]
		if !ok {
			return nil, false
		}

		current = value
	}

	return current, true
}
