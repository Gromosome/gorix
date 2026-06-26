package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/internal/access"
)

func TestResponseBodyWritesStatusContentTypeAndCommits(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
	data, _ := ctx.Status(context2.StatusCreated).ResponseEntityJSON(func() (any, error) {
		return map[string]string{"ok": "true"}, nil
	})

	if err := ctx.ResponseBodyInternal(access.Gorix, data); err != nil {
		t.Fatalf("JSON returned error: %v", err)
	}
	if recorder.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", recorder.Header().Get("Content-Type"))
	}
	if !ctx.IsCommitted() {
		t.Fatal("context should be committed")
	}
}
