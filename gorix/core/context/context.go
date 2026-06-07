package context

import (
	native "context"
	"net/http"
)

type Context struct {
	native native.Context

	W http.ResponseWriter
	R *http.Request

	status StatusCode
	params map[string]string
}

func NewContext(
	w http.ResponseWriter,
	r *http.Request,
) *Context {
	base := native.Background()

	if r != nil {
		base = r.Context()
	}

	return &Context{
		native: base,
		W:      w,
		R:      r,
		params: make(map[string]string),
	}
}

func Background() *Context {
	return &Context{
		native: native.Background(),
		params: make(map[string]string),
	}
}

func TODO() *Context {
	return &Context{
		native: native.TODO(),
		params: make(map[string]string),
	}
}
