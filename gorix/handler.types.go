package gorix

import (
	"github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/hook"
)

type Context = context.Context
type StatusCode = context.StatusCode

type Method = context.Method
type Path = context.Path
type BasePath = context.BasePath
type RouteInfo = context.RouteInfo

const (
	GET    = context.GET
	POST   = context.POST
	PUT    = context.PUT
	PATCH  = context.PATCH
	DELETE = context.DELETE
)

const (
	// 1xx Informational
	StatusContinue           = context.StatusContinue
	StatusSwitchingProtocols = context.StatusSwitchingProtocols
	StatusProcessing         = context.StatusProcessing
	StatusEarlyHints         = context.StatusEarlyHints

	// 2xx Success
	StatusOK                   = context.StatusOK
	StatusCreated              = context.StatusCreated
	StatusAccepted             = context.StatusAccepted
	StatusNonAuthoritativeInfo = context.StatusNonAuthoritativeInfo
	StatusNoContent            = context.StatusNoContent
	StatusResetContent         = context.StatusResetContent
	StatusPartialContent       = context.StatusPartialContent
	StatusMultiStatus          = context.StatusMultiStatus
	StatusAlreadyReported      = context.StatusAlreadyReported
	StatusIMUsed               = context.StatusIMUsed

	// 3xx Redirection
	StatusMultipleChoices   = context.StatusMultipleChoices
	StatusMovedPermanently  = context.StatusMovedPermanently
	StatusFound             = context.StatusFound
	StatusSeeOther          = context.StatusSeeOther
	StatusNotModified       = context.StatusNotModified
	StatusUseProxy          = context.StatusUseProxy
	StatusTemporaryRedirect = context.StatusTemporaryRedirect
	StatusPermanentRedirect = context.StatusPermanentRedirect

	// 4xx Client Error
	StatusBadRequest                   = context.StatusBadRequest
	StatusUnauthorized                 = context.StatusUnauthorized
	StatusPaymentRequired              = context.StatusPaymentRequired
	StatusForbidden                    = context.StatusForbidden
	StatusNotFound                     = context.StatusNotFound
	StatusMethodNotAllowed             = context.StatusMethodNotAllowed
	StatusNotAcceptable                = context.StatusNotAcceptable
	StatusProxyAuthRequired            = context.StatusProxyAuthRequired
	StatusRequestTimeout               = context.StatusRequestTimeout
	StatusConflict                     = context.StatusConflict
	StatusGone                         = context.StatusGone
	StatusLengthRequired               = context.StatusLengthRequired
	StatusPreconditionFailed           = context.StatusPreconditionFailed
	StatusRequestEntityTooLarge        = context.StatusRequestEntityTooLarge
	StatusRequestURITooLong            = context.StatusRequestURITooLong
	StatusUnsupportedMediaType         = context.StatusUnsupportedMediaType
	StatusRequestedRangeNotSatisfiable = context.StatusRequestedRangeNotSatisfiable
	StatusExpectationFailed            = context.StatusExpectationFailed
	StatusTeapot                       = context.StatusTeapot
	StatusMisdirectedRequest           = context.StatusMisdirectedRequest
	StatusUnprocessableEntity          = context.StatusUnprocessableEntity
	StatusLocked                       = context.StatusLocked
	StatusFailedDependency             = context.StatusFailedDependency
	StatusTooEarly                     = context.StatusTooEarly
	StatusUpgradeRequired              = context.StatusUpgradeRequired
	StatusPreconditionRequired         = context.StatusPreconditionRequired
	StatusTooManyRequests              = context.StatusTooManyRequests
	StatusRequestHeaderFieldsTooLarge  = context.StatusRequestHeaderFieldsTooLarge
	StatusUnavailableForLegalReasons   = context.StatusUnavailableForLegalReasons

	// 5xx Server Error
	StatusInternalServerError           = context.StatusInternalServerError
	StatusNotImplemented                = context.StatusNotImplemented
	StatusBadGateway                    = context.StatusBadGateway
	StatusServiceUnavailable            = context.StatusServiceUnavailable
	StatusGatewayTimeout                = context.StatusGatewayTimeout
	StatusHTTPVersionNotSupported       = context.StatusHTTPVersionNotSupported
	StatusVariantAlsoNegotiates         = context.StatusVariantAlsoNegotiates
	StatusInsufficientStorage           = context.StatusInsufficientStorage
	StatusLoopDetected                  = context.StatusLoopDetected
	StatusNotExtended                   = context.StatusNotExtended
	StatusNetworkAuthenticationRequired = context.StatusNetworkAuthenticationRequired
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
type RouteHandler = context.RouteHandler
type ValidationError = context.ValidationError
type FieldError = context.FieldError

var Apply = hook.Apply
var ApplyInterceptor = hook.ApplyInterceptor
var ApplyFilter = hook.ApplyFilter
var NewValidationError = context.NewValidationError
var NewFieldError = context.NewFieldError
