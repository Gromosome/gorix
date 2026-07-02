package context

import (
	"net/http/httptest"
	"strings"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

type bindDTO struct {
	Name   string   `json:"name" query:"name" param:"name" validate:"required"`
	Age    int      `json:"age" query:"age" param:"age"`
	Active bool     `json:"active" query:"active"`
	Tags   []string `json:"tags" query:"tags"`
}

func TestBindBodyDecodesJSONAndValidates(t *testing.T) {
	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"bob","age":30,"active":true,"tags":["a","b"]}`))
	request.Header.Set("Content-Type", "application/json")
	ctx := context2.NewContext(nil, request)

	var target bindDTO
	err := ctx.BindJSONBody(&target)
	if err != nil {
		t.Fatalf("BindBody returned error: %v", err)
	}
	if target.Name != "bob" || target.Age != 30 || !target.Active || len(target.Tags) != 2 {
		t.Fatalf("unexpected bound body: %#v", target)
	}
}

func TestBindBodyRejectsUnknownFields(t *testing.T) {
	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"bob","unknown":true}`))
	ctx := context2.NewContext(nil, request)

	var target bindDTO
	err := ctx.BindJSONBody(&target)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestBindQueryBindsScalarAndSliceValues(t *testing.T) {
	request := httptest.NewRequest("GET", "/?name=bob&age=30&active=true&tags=a,b&tags=c", nil)
	ctx := context2.NewContext(nil, request)

	var target bindDTO
	err := ctx.BindQuery(&target)
	if err != nil {
		t.Fatalf("BindQuery returned error: %v", err)
	}
	if target.Name != "bob" || target.Age != 30 || !target.Active {
		t.Fatalf("unexpected bound query: %#v", target)
	}
	if len(target.Tags) != 3 || target.Tags[2] != "c" {
		t.Fatalf("unexpected tags: %#v", target.Tags)
	}
}

func TestBindParamsBindsPathParameters(t *testing.T) {
	ctx := context2.Background()
	ctx.SetParams(map[string]string{"name": "bob", "age": "30"})

	var target bindDTO
	err := ctx.BindParams(&target)
	if err != nil {
		t.Fatalf("BindParams returned error: %v", err)
	}
	if target.Name != "bob" || target.Age != 30 {
		t.Fatalf("unexpected bound params: %#v", target)
	}
}

func TestBindValuesRejectsUnsupportedType(t *testing.T) {
	type invalidDTO struct {
		Value map[string]string `query:"value"`
	}
	request := httptest.NewRequest("GET", "/?value=x", nil)
	ctx := context2.NewContext(nil, request)

	var target invalidDTO
	err := ctx.BindQuery(&target)
	if err != nil {
		t.Logf("expected unsupported type bind error %v", err)
	} else {
		t.Fatal("not throwing expected unsupported type bind error")
	}
}
