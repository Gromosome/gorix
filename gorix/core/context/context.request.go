package context

import (
	"errors"
	"net/http"
	"strings"
)

func (c *Context) Request() *http.Request {
	if c == nil {
		return nil
	}
	return c.R
}

func (c *Context) contentType() string {
	if c == nil || c.R == nil {
		return ""
	}
	value := c.R.Header.Get("Content-Type")
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ";")
	return strings.ToLower(strings.TrimSpace(parts[0]))
}

func (c *Context) IsCancelled() bool {
	return c != nil && c.Err() != nil
}

func (c *Context) SetParams(
	params map[string]string,
) {
	c.params = cloneParams(params)
}

func (c *Context) Param(
	key string,
) string {
	if c == nil || c.params == nil {
		return ""
	}
	return c.params[key]
}

func (c *Context) Query(
	key string,
) string {
	if c == nil || c.R == nil {
		return ""
	}
	return c.R.URL.Query().Get(key)
}

func (c *Context) QueryDefault(
	key string,
	defaultValue string,
) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Context) ParamTo(
	key string,
	target *string,
) *Context {
	if c == nil {
		return c
	}
	if c.bindingErr != nil {
		return c
	}
	if target == nil {
		c.bindingErr = errors.New("param target cannot be nil")
		return c
	}
	if c.params == nil {
		*target = ""
		return c
	}
	*target = c.params[key]
	return c
}

func (c *Context) QueryTo(
	key string,
	target *string,
) *Context {
	if c == nil {
		return c
	}
	if c.bindingErr != nil {
		return c
	}
	if target == nil {
		c.bindingErr = errors.New("query target cannot be nil")
		return c
	}
	if c.R == nil || c.R.URL == nil {
		*target = ""
		return c
	}
	*target = c.R.URL.Query().Get(key)
	return c
}

func (c *Context) QueryDefaultTo(
	key string,
	defaultValue string,
	target *string,
) *Context {
	if c == nil {
		return c
	}
	if c.bindingErr != nil {
		return c
	}
	if target == nil {
		c.bindingErr = errors.New("query default target cannot be nil")
		return c
	}
	if c.R == nil || c.R.URL == nil {
		*target = defaultValue
		return c
	}
	value := c.R.URL.Query().Get(key)
	if value == "" {
		*target = defaultValue
		return c
	}
	*target = value
	return c
}
