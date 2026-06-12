package hook

import (
	"reflect"
	"testing"

	"github.com/Gromosome/gorix/gorix/core/context"
	hook2 "github.com/Gromosome/gorix/gorix/hook"
)

func TestApplyMiddlewareBuildsRouteRules(t *testing.T) {
	middleware := func(next hook2.Handler) hook2.Handler {
		return next
	}

	only := hook2.Apply(middleware).Only("/api/*")
	if reflect.ValueOf(only.Middleware).Pointer() != reflect.ValueOf(middleware).Pointer() {
		t.Fatal("Only did not preserve middleware")
	}
	if !only.Rule.Match("/api/users") || only.Rule.Match("/web/users") {
		t.Fatalf("unexpected Only rule matching")
	}

	except := hook2.Apply(middleware).Except("/health")
	if !except.Rule.Match("/users") || except.Rule.Match("/health") {
		t.Fatalf("unexpected Except rule matching")
	}
}

func TestChainMiddlewaresWrapsInDeclarationOrder(t *testing.T) {
	calls := make([]string, 0)

	first := func(next hook2.Handler) hook2.Handler {
		return func(c *context.Context) error {
			calls = append(calls, "first-before")
			err := next(nil)
			calls = append(calls, "first-after")
			return err
		}
	}
	second := func(next hook2.Handler) hook2.Handler {
		return func(c *context.Context) error {
			calls = append(calls, "second-before")
			err := next(nil)
			calls = append(calls, "second-after")
			return err
		}
	}
	base := func(c *context.Context) error {
		calls = append(calls, "base")
		return nil
	}

	handler := hook2.ChainMiddlewares(hook2.Handler(base), hook2.Middleware(first), hook2.Middleware(second))
	if err := handler(nil); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	want := []string{"first-before", "second-before", "base", "second-after", "first-after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected middleware order: %#v", calls)
	}
}
