package config

import (
	"strings"
)

type YAMLValue = any

func ParseYAML(content string) map[string]YAMLValue {
	root := make(map[string]YAMLValue)

	lines := normalizeYAMLLines(content)

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

		key, value, hasValue := splitYAMLKeyValue(trimmed)
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

		parent[key] = parseScalarOrInline(value)
	}

	return root
}
