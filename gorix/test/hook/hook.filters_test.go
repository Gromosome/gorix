package hook

import (
	"testing"

	hook2 "github.com/Gromosome/gorix/gorix/hook"
)

type testFilter struct{}

func (testFilter) Catch(ctx *hook2.ExceptionContext) {}

func TestApplyFilterBuildsRouteRules(t *testing.T) {
	filter := testFilter{}

	only := hook2.ApplyFilter(filter).Only("/users", "/admin/*")
	if only.Filter != filter {
		t.Fatal("Only did not preserve filter")
	}
	if !only.Rule.Match("/admin/settings") {
		t.Fatal("Only rule did not match expected path")
	}
	if only.Rule.Match("/other") {
		t.Fatal("Only rule matched unexpected path")
	}

	except := hook2.ApplyFilter(filter).Except("/health")
	if !except.Rule.Match("/users") {
		t.Fatal("Except rule rejected allowed path")
	}
	if except.Rule.Match("/health") {
		t.Fatal("Except rule accepted excluded path")
	}
}

func TestGlobalFilterMatchesAllRoutes(t *testing.T) {
	config := hook2.GlobalFilter(testFilter{})
	if !config.Rule.Match("/anything") {
		t.Fatal("global filter should match all routes")
	}
}
