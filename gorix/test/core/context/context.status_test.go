package context

import (
	"net/http"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestStatusCode(t *testing.T) {
	if context2.StatusCreated.Int() != http.StatusCreated {
		t.Fatalf("unexpected created status: %d", context2.StatusCreated.Int())
	}
	if context2.StatusNotFound.Text() != http.StatusText(http.StatusNotFound) {
		t.Fatalf("unexpected status text: %s", context2.StatusNotFound.Text())
	}
}
