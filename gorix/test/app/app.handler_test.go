package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	app2 "github.com/Gromosome/gorix/gorix/app"
	"github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/hook"
)

type routeTestModule struct{}

func (routeTestModule) BasePath() context.BasePath {
	return "/api"
}

func (routeTestModule) Controllers() []any {
	return []any{func() *routeTestController { return &routeTestController{} }}
}

type routeTestController struct{}

func (routeTestController) GetUser() (context.Method, context.Path, context.RouteHandler) {
	return context.GET, "/users/:id", func(c *context.Context) (any, error) {
		return map[string]string{"id": c.Param("id")}, nil
	}
}

func (routeTestController) Summary() (context.Method, context.Path, context.RouteHandler) {
	return context.GET, "/users/summary", func(c *context.Context) (any, error) {
		return map[string]string{"summary": "true"}, nil
	}
}

func TestTryRegisterModulesRegistersAndDispatchesRoutes(t *testing.T) {
	instance := app2.NewApp()
	if err := instance.TryRegisterModules(&routeTestModule{}); err != nil {
		t.Fatalf("TryRegisterModules returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/users/42", nil)
	instance.Dispatch(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", recorder.Code, recorder.Body.String())
	}

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["id"] != "42" {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestDispatchPrefersStaticRouteOverParamRoute(t *testing.T) {
	instance := app2.NewApp()
	if err := instance.TryRegisterModules(&routeTestModule{}); err != nil {
		t.Fatalf("TryRegisterModules returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/users/summary", nil)
	instance.Dispatch(recorder, request)

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["summary"] != "true" {
		t.Fatalf("static route was not preferred: %#v", body)
	}
}

func TestDispatchHandlesNotFoundAndMethodNotAllowed(t *testing.T) {
	instance := app2.NewApp()
	if err := instance.TryRegisterModules(&routeTestModule{}); err != nil {
		t.Fatalf("TryRegisterModules returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	instance.Dispatch(recorder, httptest.NewRequest(http.MethodPost, "/api/users/42", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected method-not-allowed status: %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	instance.Dispatch(recorder, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("unexpected not-found status: %d", recorder.Code)
	}
}

func TestRouteHelpers(t *testing.T) {
	matched, params, score := app2.MatchRoute("/users/:id", "/users/42")
	if !matched {
		t.Fatal("expected route to match")
	}
	if params["id"] != "42" || score != 11 {
		t.Fatalf("unexpected params or score: %#v %d", params, score)
	}

	if got := app2.NormalizeRoute("/api/", "/users/"); got != "/api/users" {
		t.Fatalf("unexpected normalized route: %s", got)
	}
	if got := app2.SplitRoutePath("/api/users/"); !reflect.DeepEqual(got, []string{"api", "users"}) {
		t.Fatalf("unexpected route parts: %#v", got)
	}
}

func TestResolveHooksByPath(t *testing.T) {
	instance := app2.NewApp()
	middleware := func(next hook.Handler) hook.Handler { return next }
	interceptor := testAppInterceptor{}
	filter := testAppFilter{}

	instance.Use(hook.Apply(middleware).Only("/api/*"))
	instance.UseInterceptors(hook.ApplyInterceptor(interceptor).Only("/api/*"))
	instance.UseFilters(hook.ApplyFilter(filter).Only("/api/*"))

	if len(instance.ResolveMiddlewares("/api/users")) != 1 {
		t.Fatal("middleware was not resolved")
	}
	if len(instance.ResolveInterceptors("/api/users")) != 1 {
		t.Fatal("interceptor was not resolved")
	}
	if len(instance.ResolveFilters("/api/users")) != 1 {
		t.Fatal("filter was not resolved")
	}
	if len(instance.ResolveFilters("/other")) != 0 {
		t.Fatal("filter should not match other route")
	}
}

type testAppInterceptor struct{}

func (testAppInterceptor) Before(ctx *hook.ExecutionContext) error { return nil }
func (testAppInterceptor) After(ctx *hook.ExecutionContext) error  { return nil }

type testAppFilter struct{}

func (testAppFilter) Catch(ctx *hook.ExceptionContext) {}
