package hook

import (
	"time"

	"github.com/Gromosome/gorix/gorix/core"
)

type ExecutionContext struct {
	Context *core.Context

	Method core.Method
	Path   string

	Module     string
	Controller string
	Handler    string

	StartTime time.Time
	Response  any
	Error     error
}

type Interceptor interface {
	Before(ctx *ExecutionContext) error
	After(ctx *ExecutionContext) error
}

type InterceptorConfig struct {
	Interceptor Interceptor
	Rule        RouteRule
}

type InterceptorBuilder struct {
	config InterceptorConfig
}

func ApplyInterceptor(interceptor Interceptor) InterceptorBuilder {
	return InterceptorBuilder{
		config: InterceptorConfig{
			Interceptor: interceptor,
			Rule:        RouteRule{},
		},
	}
}

func (b InterceptorBuilder) Only(paths ...string) InterceptorConfig {
	b.config.Rule.OnlyPaths = paths
	return b.config
}

func (b InterceptorBuilder) Except(paths ...string) InterceptorConfig {
	b.config.Rule.ExceptPaths = paths
	return b.config
}

func GlobalInterceptor(interceptor Interceptor) InterceptorConfig {
	return InterceptorConfig{
		Interceptor: interceptor,
		Rule:        RouteRule{},
	}
}
