package context

import (
	native "context"
	"net/http"
)

type ResponseType string

const (
	ResponseTypeAuto      ResponseType = "auto"
	ResponseTypeJSON      ResponseType = "json"
	ResponseTypeXML       ResponseType = "xml"
	ResponseTypeText      ResponseType = "text"
	ResponseTypeHTML      ResponseType = "html"
	ResponseTypeFile      ResponseType = "file"
	ResponseTypeDownload  ResponseType = "download"
	ResponseTypeStream    ResponseType = "stream"
	ResponseTypeNoContent ResponseType = "no_content"
	ResponseTypeRedirect  ResponseType = "redirect"
)

type Context struct {
	native native.Context

	W            http.ResponseWriter
	R            *http.Request
	responseType ResponseType
	status       StatusCode
	committed    bool
	params       map[string]string
	err          error
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
		err:    nil,
	}
}

func (c *Context) setError(err error) {
	c.err = err
}

func (c *Context) ResponseType() ResponseType {
	if c.responseType == "" {
		return ResponseTypeAuto
	}
	return c.responseType
}

func (c *Context) setResponseType(responseType ResponseType) {
	c.responseType = responseType
}

func Background() *Context {
	return &Context{
		native: native.Background(),
		params: make(map[string]string),
	}
}
func (c *Context) IsCommitted() bool {
	return c.committed
}
func (c *Context) Params() map[string]string {
	return c.params
}

func TODO() *Context {
	return &Context{
		native: native.TODO(),
		params: make(map[string]string),
	}
}
