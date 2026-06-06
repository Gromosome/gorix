package context

func (c *Context) SetParams(params map[string]string) {
	c.params = params
}

func (c *Context) Param(key string) string {
	if c.params == nil {
		return ""
	}

	return c.params[key]
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) QueryDefault(key string, defaultValue string) string {
	value := c.Query(key)

	if value == "" {
		return defaultValue
	}

	return value
}
