package context

import (
	"net/http/httptest"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestNewContextUsesRequestContextAndInitializesParams(t *testing.T) {
	request := httptest.NewRequest("GET", "/users", nil)
	ctx := context2.NewContext(httptest.NewRecorder(), request)

	if ctx.Request() != request {
		t.Fatal("context did not retain request")
	}
	if ctx.Params() == nil {
		t.Fatal("context params map was not initialized")
	}
	if ctx.Native() != request.Context() {
		t.Fatal("context did not use request native context")
	}
}

func TestBackgroundAndTODOInitializeNativeContexts(t *testing.T) {
	if context2.Background().Native() == nil {
		t.Fatal("Background native context is nil")
	}
	if context2.TODO().Native() == nil {
		t.Fatal("TODO native context is nil")
	}
}
