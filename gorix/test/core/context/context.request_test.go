package context

import (
	"net/http/httptest"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestRequestHelpers(t *testing.T) {
	request := httptest.NewRequest("GET", "/users?name=bob", nil)
	ctx := context2.NewContext(nil, request)
	ctx.SetParams(map[string]string{"id": "42"})

	if ctx.Param("id") != "42" {
		t.Fatalf("unexpected param: %s", ctx.Param("id"))
	}
	if ctx.Query("name") != "bob" {
		t.Fatalf("unexpected query: %s", ctx.Query("name"))
	}
	if ctx.QueryDefault("missing", "fallback") != "fallback" {
		t.Fatal("QueryDefault did not return fallback")
	}
}

func TestNilContextRequestHelpersAreSafe(t *testing.T) {
	var ctx *context2.Context
	if ctx.Request() != nil {
		t.Fatal("nil context Request should return nil")
	}
	if ctx.Param("id") != "" || ctx.Query("q") != "" {
		t.Fatal("nil context accessors should return empty strings")
	}
}
