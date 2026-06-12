package hook

import (
	"testing"

	hook2 "github.com/Gromosome/gorix/gorix/hook"
)

func TestRouteRuleMatch(t *testing.T) {
	tests := []struct {
		name string
		rule hook2.RouteRule
		path string
		want bool
	}{
		{name: "empty rule matches all", rule: hook2.RouteRule{}, path: "/users", want: true},
		{name: "only exact match", rule: hook2.RouteRule{OnlyPaths: []string{"/users"}}, path: "users/", want: true},
		{name: "only rejects non-match", rule: hook2.RouteRule{OnlyPaths: []string{"/users"}}, path: "/admin", want: false},
		{name: "wildcard matches all", rule: hook2.RouteRule{OnlyPaths: []string{"*"}}, path: "/admin", want: true},
		{name: "prefix wildcard matches child", rule: hook2.RouteRule{OnlyPaths: []string{"/api/*"}}, path: "/api/users", want: true},
		{name: "prefix wildcard matches prefix", rule: hook2.RouteRule{OnlyPaths: []string{"/api/*"}}, path: "/api", want: true},
		{name: "except rejects match", rule: hook2.RouteRule{ExceptPaths: []string{"/health"}}, path: "/health", want: false},
		{name: "except allows non-match", rule: hook2.RouteRule{ExceptPaths: []string{"/health"}}, path: "/users", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Match(tt.path); got != tt.want {
				t.Fatalf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := map[string]string{
		"":          "/",
		"users":     "/users",
		"/users/":   "/users",
		"/users/42": "/users/42",
	}

	for input, want := range tests {
		if got := hook2.NormalizePath(input); got != want {
			t.Fatalf("NormalizePath(%q) = %q, want %q", input, got, want)
		}
	}
}
