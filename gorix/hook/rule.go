package hook

import "strings"

type RouteRule struct {
	OnlyPaths   []string
	ExceptPaths []string
}

func (r RouteRule) Match(path string) bool {
	if len(r.OnlyPaths) > 0 {
		return matchAnyPattern(path, r.OnlyPaths)
	}

	if len(r.ExceptPaths) > 0 {
		return !matchAnyPattern(path, r.ExceptPaths)
	}

	return true
}

func matchAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(path, pattern) {
			return true
		}
	}

	return false
}

func matchPattern(path string, pattern string) bool {
	path = normalizePath(path)
	pattern = normalizePath(pattern)

	if pattern == "*" || pattern == "/*" {
		return true
	}

	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}

	return path == pattern
}

func normalizePath(path string) string {
	if path == "" {
		return "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}

	return path
}
