package context

import (
	native "context"
	"net/http"

	"github.com/Gromosome/gorix/gorix/logger"
)

type ResponseType string

const (
	ResponseTypeSOAP12 ResponseType = "soap12"
	ResponseTypeSOAP11 ResponseType = "soap11"
	ResponseTypeJSON   ResponseType = "json"
	ResponseTypeXML    ResponseType = "xml"
)

type Context struct {
	native native.Context

	W            http.ResponseWriter
	R            *http.Request
	responseType ResponseType
	status       int
	committed    bool
	params       map[string]string
	bindingErr   error
	logger       *logger.Logger
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
		native:     base,
		W:          w,
		R:          r,
		params:     make(map[string]string),
		bindingErr: nil,
		logger:     logger.NewLogger(logger.Wrap1),
	}
}
func (c *Context) setStatus(status Code) {
	c.status = status.Int()
}
func (c *Context) GetStatusOrDefault(status Code) int {
	if c == nil || c.status == 0 {
		return status.Int()
	}
	return c.status
}
func (c *Context) Status(status StatusCode) *Context {
	if c == nil {
		return c
	}
	c.setStatus(status)
	return c
}
func (c *Context) SOAPStatus(status SOAPStatusCode) *Context {
	if c == nil {
		return c
	}
	c.setStatus(status)
	return c
}
func (c *Context) SetHeader(key, value string) *Context {
	if c != nil && c.W != nil {
		c.W.Header().Set(key, value)
	}
	return c
}
func (c *Context) GetHeader(key string) string {
	return c.R.Header.Get(key)
}
func (c *Context) GetHeaderOf(key string, target *string) *Context {
	*target = c.Request().Header.Get(key)
	return c
}

func (c *Context) setBindingError(err error) {
	c.bindingErr = err
}

func (c *Context) GetBindingErr() error {
	return c.bindingErr
}

func (c *Context) ResponseType() ResponseType {
	if c.responseType == "" {
		return ResponseTypeJSON
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
