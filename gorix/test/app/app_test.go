package app

import (
	"strings"
	"testing"

	app2 "github.com/Gromosome/gorix/gorix/app"
	"github.com/Gromosome/gorix/gorix/core/context"
)

type providerDependency struct {
	value string
}

type providerModuleFixture struct{}

func (providerModuleFixture) Providers() []any {
	return []any{func() *providerDependency { return &providerDependency{value: "ok"} }}
}

func (providerModuleFixture) BasePath() context.BasePath {
	return "/provider"
}

func (providerModuleFixture) Controllers() []any {
	return []any{func(dep *providerDependency) *providerController { return &providerController{dep: dep} }}
}

type providerController struct {
	dep *providerDependency
}

func (c providerController) Index() (context.Method, context.Path, context.RouteHandler) {
	return context.GET, "/", func(ctx *context.Context) (any, error) {
		return c.dep.value, nil
	}
}

func TestTryRegisterModulesValidatesModules(t *testing.T) {
	instance := app2.NewApp()

	if err := instance.TryRegisterModules(nil); err == nil {
		t.Fatal("nil module should be rejected")
	}
	if err := instance.TryRegisterModules(providerModuleFixture{}); err == nil {
		t.Fatal("non-pointer module should be rejected")
	}
}

func TestTryRegisterModulesRegistersProvidersBeforeControllers(t *testing.T) {
	instance := app2.NewApp()

	if err := instance.TryRegisterModules(&providerModuleFixture{}); err != nil {
		t.Fatalf("TryRegisterModules returned error: %v", err)
	}
	if len(instance.RouteEntries()) != 1 {
		t.Fatalf("unexpected route count: %d", len(instance.RouteEntries()))
	}
}

func TestTryRegisterModulesRejectsDuplicateRoutes(t *testing.T) {
	instance := app2.NewApp()

	err := instance.TryRegisterModules(&routeTestModule{}, &routeTestModule{})
	if err == nil {
		t.Fatal("expected duplicate route error")
	}
	if !strings.Contains(err.Error(), "duplicate route") {
		t.Fatalf("unexpected error: %v", err)
	}
}
