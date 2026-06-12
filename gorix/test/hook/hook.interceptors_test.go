package hook

import (
	"testing"

	hook2 "github.com/Gromosome/gorix/gorix/hook"
)

type testInterceptor struct{}

func (testInterceptor) Before(ctx *hook2.ExecutionContext) error { return nil }
func (testInterceptor) After(ctx *hook2.ExecutionContext) error  { return nil }

func TestApplyInterceptorBuildsRouteRules(t *testing.T) {
	interceptor := testInterceptor{}

	only := hook2.ApplyInterceptor(interceptor).Only("/api/*")
	if only.Interceptor != interceptor {
		t.Fatal("Only did not preserve interceptor")
	}
	if !only.Rule.Match("/api/users") || only.Rule.Match("/web/users") {
		t.Fatalf("unexpected Only rule matching")
	}

	except := hook2.ApplyInterceptor(interceptor).Except("/private/*")
	if !except.Rule.Match("/public") || except.Rule.Match("/private/users") {
		t.Fatalf("unexpected Except rule matching")
	}
}

func TestGlobalInterceptorMatchesAllRoutes(t *testing.T) {
	config := hook2.GlobalInterceptor(testInterceptor{})
	if !config.Rule.Match("/anything") {
		t.Fatal("global interceptor should match all routes")
	}
}
