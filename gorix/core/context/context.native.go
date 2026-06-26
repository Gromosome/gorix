package context

import (
	native "context"
	"time"
)

type CancelFunc func()

func (c CancelFunc) Cancel() {
	if c != nil {
		c()
	}
}

func (c *Context) Deadline() (
	time.Time,
	bool,
) {
	return c.ensureNative().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ensureNative().Done()
}

func (c *Context) Err() error {
	return c.ensureNative().Err()
}

func (c *Context) Value(key any) any {
	return c.ensureNative().Value(key)
}

func (c *Context) ensureNative() native.Context {
	if c == nil || c.native == nil {
		return native.Background()
	}

	return c.native
}

func (c *Context) Native() native.Context {
	return c.ensureNative()
}

func (c *Context) SetNative(
	value native.Context,
) *Context {
	if value == nil {
		value = native.Background()
	}

	c.native = value

	if c.R != nil {
		c.R = c.R.WithContext(value)
	}

	return c
}

func WithCancel(
	parent *Context,
) (
	*Context,
	CancelFunc,
) {
	if parent == nil {
		parent = Background()
	}
	child, cancel := native.WithCancel(
		parent.ensureNative(),
	)
	return parent.cloneWithNative(child), CancelFunc(cancel)
}

func WithTimeout(
	parent *Context,
	timeout time.Duration,
) (
	*Context,
	CancelFunc,
) {
	if parent == nil {
		parent = Background()
	}
	child, cancel := native.WithTimeout(
		parent.ensureNative(),
		timeout,
	)
	return parent.cloneWithNative(child), CancelFunc(cancel)
}

func WithDeadline(
	parent *Context,
	deadline time.Time,
) (
	*Context,
	CancelFunc,
) {
	if parent == nil {
		parent = Background()
	}
	child, cancel := native.WithDeadline(
		parent.ensureNative(),
		deadline,
	)
	return parent.cloneWithNative(child), CancelFunc(cancel)
}

func WithValue(
	parent *Context,
	key any,
	value any,
) *Context {
	if parent == nil {
		parent = Background()
	}
	child := native.WithValue(
		parent.ensureNative(),
		key,
		value,
	)
	return parent.cloneWithNative(child)
}

func (c *Context) cloneWithNative(
	value native.Context,
) *Context {
	copyContext := &Context{
		native: value,
		W:      c.W,
		R:      c.R,
		status: c.status,
		params: cloneParams(c.params),
	}
	if copyContext.R != nil {
		copyContext.R = copyContext.R.WithContext(value)
	}
	return copyContext
}

func cloneParams(
	source map[string]string,
) map[string]string {
	target := make(map[string]string, len(source))
	for key, value := range source {
		target[key] = value
	}
	return target
}
