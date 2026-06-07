package context

import "net/http"

func (c *Context) Request() *http.Request {
	if c == nil {
		return nil
	}

	return c.R
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
