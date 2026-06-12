package context

import (
	native "context"
	"net/http/httptest"
	"testing"
	"time"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestWithCancelCancelsChildContextOnly(t *testing.T) {
	parent := context2.Background()
	parent.SetParams(map[string]string{"id": "1"})

	child, cancel := context2.WithCancel(parent)
	cancel.Cancel()

	if child.Err() == nil {
		t.Fatal("child context was not cancelled")
	}
	if parent.Err() != nil {
		t.Fatal("parent context was unexpectedly cancelled")
	}
	if child.Param("id") != "1" {
		t.Fatal("child context did not clone params")
	}
}

func TestWithTimeoutSetsDeadline(t *testing.T) {
	child, cancel := context2.WithTimeout(nil, time.Second)
	defer cancel.Cancel()

	if _, ok := child.Deadline(); !ok {
		t.Fatal("WithTimeout child has no deadline")
	}
}

func TestWithValueAndSetNative(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	ctx := context2.NewContext(nil, request)

	ctx = context2.WithValue(ctx, "key", "value")
	if ctx.Value("key") != "value" {
		t.Fatal("WithValue did not store value")
	}

	replacement := native.WithValue(native.Background(), "other", "value")
	ctx.SetNative(replacement)
	if ctx.R.Context().Value("other") != "value" {
		t.Fatal("SetNative did not update request context")
	}
}
