package hook

import (
	"github.com/Gromosome/gorix/gorix/core/context"
)

type ExceptionContext struct {
	Context    *context.Context
	Method     context.Method
	Path       string
	Module     string
	Controller string
	Handler    string
	Error      error
	StatusCode context.StatusCode
}

type Filter interface {
	Catch(ctx *ExceptionContext)
}

type FilterConfig struct {
	Filter Filter
	Rule   RouteRule
}

type FilterBuilder struct {
	config FilterConfig
}

func ApplyFilter(filter Filter) FilterBuilder {
	return FilterBuilder{
		config: FilterConfig{
			Filter: filter,
			Rule:   RouteRule{},
		},
	}
}

func (b FilterBuilder) Only(paths ...string) FilterConfig {
	b.config.Rule.OnlyPaths = paths
	return b.config
}

func (b FilterBuilder) Except(paths ...string) FilterConfig {
	b.config.Rule.ExceptPaths = paths
	return b.config
}

func GlobalFilter(filter Filter) FilterConfig {
	return FilterConfig{
		Filter: filter,
		Rule:   RouteRule{},
	}
}
