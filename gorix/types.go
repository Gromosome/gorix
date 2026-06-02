package gorix

import (
	"github.com/Gromosome/gorix/gorix/core"
	"github.com/Gromosome/gorix/gorix/hook"
)

type Context = core.Context
type StatusCode = core.StatusCode

type Method = core.Method
type Path = core.Path
type BasePath = core.BasePath
type RouteInfo = core.RouteInfo

const (
	GET    = core.GET
	POST   = core.POST
	PUT    = core.PUT
	PATCH  = core.PATCH
	DELETE = core.DELETE
)

const (
	StatusOK                  = core.StatusOK
	StatusCreated             = core.StatusCreated
	StatusAccepted            = core.StatusAccepted
	StatusNoContent           = core.StatusNoContent
	StatusBadRequest          = core.StatusBadRequest
	StatusUnauthorized        = core.StatusUnauthorized
	StatusForbidden           = core.StatusForbidden
	StatusNotFound            = core.StatusNotFound
	StatusMethodNotAllowed    = core.StatusMethodNotAllowed
	StatusInternalServerError = core.StatusInternalServerError
)

type Handler = hook.Handler
type Middleware = hook.Middleware
type Interceptor = hook.Interceptor
type Filter = hook.Filter

type ExecutionContext = hook.ExecutionContext
type ExceptionContext = hook.ExceptionContext

type MiddlewareConfig = hook.MiddlewareConfig
type InterceptorConfig = hook.InterceptorConfig
type FilterConfig = hook.FilterConfig

var Apply = hook.Apply
var ApplyInterceptor = hook.ApplyInterceptor
var ApplyFilter = hook.ApplyFilter
