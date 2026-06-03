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
	// 1xx Informational
	StatusContinue           = core.StatusContinue
	StatusSwitchingProtocols = core.StatusSwitchingProtocols
	StatusProcessing         = core.StatusProcessing
	StatusEarlyHints         = core.StatusEarlyHints

	// 2xx Success
	StatusOK                   = core.StatusOK
	StatusCreated              = core.StatusCreated
	StatusAccepted             = core.StatusAccepted
	StatusNonAuthoritativeInfo = core.StatusNonAuthoritativeInfo
	StatusNoContent            = core.StatusNoContent
	StatusResetContent         = core.StatusResetContent
	StatusPartialContent       = core.StatusPartialContent
	StatusMultiStatus          = core.StatusMultiStatus
	StatusAlreadyReported      = core.StatusAlreadyReported
	StatusIMUsed               = core.StatusIMUsed

	// 3xx Redirection
	StatusMultipleChoices   = core.StatusMultipleChoices
	StatusMovedPermanently  = core.StatusMovedPermanently
	StatusFound             = core.StatusFound
	StatusSeeOther          = core.StatusSeeOther
	StatusNotModified       = core.StatusNotModified
	StatusUseProxy          = core.StatusUseProxy
	StatusTemporaryRedirect = core.StatusTemporaryRedirect
	StatusPermanentRedirect = core.StatusPermanentRedirect

	// 4xx Client Error
	StatusBadRequest                   = core.StatusBadRequest
	StatusUnauthorized                 = core.StatusUnauthorized
	StatusPaymentRequired              = core.StatusPaymentRequired
	StatusForbidden                    = core.StatusForbidden
	StatusNotFound                     = core.StatusNotFound
	StatusMethodNotAllowed             = core.StatusMethodNotAllowed
	StatusNotAcceptable                = core.StatusNotAcceptable
	StatusProxyAuthRequired            = core.StatusProxyAuthRequired
	StatusRequestTimeout               = core.StatusRequestTimeout
	StatusConflict                     = core.StatusConflict
	StatusGone                         = core.StatusGone
	StatusLengthRequired               = core.StatusLengthRequired
	StatusPreconditionFailed           = core.StatusPreconditionFailed
	StatusRequestEntityTooLarge        = core.StatusRequestEntityTooLarge
	StatusRequestURITooLong            = core.StatusRequestURITooLong
	StatusUnsupportedMediaType         = core.StatusUnsupportedMediaType
	StatusRequestedRangeNotSatisfiable = core.StatusRequestedRangeNotSatisfiable
	StatusExpectationFailed            = core.StatusExpectationFailed
	StatusTeapot                       = core.StatusTeapot
	StatusMisdirectedRequest           = core.StatusMisdirectedRequest
	StatusUnprocessableEntity          = core.StatusUnprocessableEntity
	StatusLocked                       = core.StatusLocked
	StatusFailedDependency             = core.StatusFailedDependency
	StatusTooEarly                     = core.StatusTooEarly
	StatusUpgradeRequired              = core.StatusUpgradeRequired
	StatusPreconditionRequired         = core.StatusPreconditionRequired
	StatusTooManyRequests              = core.StatusTooManyRequests
	StatusRequestHeaderFieldsTooLarge  = core.StatusRequestHeaderFieldsTooLarge
	StatusUnavailableForLegalReasons   = core.StatusUnavailableForLegalReasons

	// 5xx Server Error
	StatusInternalServerError           = core.StatusInternalServerError
	StatusNotImplemented                = core.StatusNotImplemented
	StatusBadGateway                    = core.StatusBadGateway
	StatusServiceUnavailable            = core.StatusServiceUnavailable
	StatusGatewayTimeout                = core.StatusGatewayTimeout
	StatusHTTPVersionNotSupported       = core.StatusHTTPVersionNotSupported
	StatusVariantAlsoNegotiates         = core.StatusVariantAlsoNegotiates
	StatusInsufficientStorage           = core.StatusInsufficientStorage
	StatusLoopDetected                  = core.StatusLoopDetected
	StatusNotExtended                   = core.StatusNotExtended
	StatusNetworkAuthenticationRequired = core.StatusNetworkAuthenticationRequired
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
type RouteHandler = core.RouteHandler

var Apply = hook.Apply
var ApplyInterceptor = hook.ApplyInterceptor
var ApplyFilter = hook.ApplyFilter
