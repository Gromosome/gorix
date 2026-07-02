package app

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/hook"
	"github.com/Gromosome/gorix/gorix/internal/access"
	"github.com/Gromosome/gorix/gorix/internal/global"
)

func (a *App) registerModuleControllers(module any) error {
	moduleValue := reflect.ValueOf(module)

	if moduleValue.Kind() != reflect.Pointer {
		return fmt.Errorf("gorix: module must be pointer, got %s", moduleValue.Kind())
	}

	moduleType := moduleValue.Type()
	moduleName := moduleType.Elem().Name()

	basePathProvider, ok := module.(basePathModule)
	if !ok {
		return fmt.Errorf("gorix: module %s must implement BasePath() gorix.BasePath", moduleName)
	}

	basePath := basePathProvider.BasePath()

	controllersModule, ok := module.(controllerModule)
	if !ok {
		return fmt.Errorf("gorix: module %s must implement Controllers() []any", moduleName)
	}

	for _, controllerItem := range controllersModule.Controllers() {
		controllerConstructor, controllerBasePath, controllerVersion, err := resolveControllerRegistration(controllerItem)
		if err != nil {
			return fmt.Errorf(
				"gorix: module %s invalid controller registration: %w",
				moduleName,
				err,
			)
		}

		controllerValue, err := a.container.Build(controllerConstructor)
		if err != nil {
			return fmt.Errorf("gorix: module %s controller build failed: %w", moduleName, err)
		}

		if controllerValue.Kind() == reflect.Pointer {
			controllerValue = controllerValue.Elem()
		}

		if controllerValue.Kind() != reflect.Struct {
			return fmt.Errorf(
				"gorix: module %s controller constructor must return controller struct or pointer to struct",
				moduleName,
			)
		}
		moduleVersion := context.VersionNeutral

		if versioned, ok := module.(versionedModule); ok {
			moduleVersion = versioned.APIVersion()
		}

		fullBasePath := buildVersionedBasePath(
			a.apiPrefix,
			moduleVersion,
			controllerVersion,
			basePath,
			controllerBasePath,
		)

		if err := a.registerController(moduleName, fullBasePath, controllerValue); err != nil {
			return err
		}
	}

	return nil
}
func resolveControllerRegistration(item any) (
	constructor any,
	basePath context.BasePath,
	controllerVersion context.APIVersion,
	err error,
) {
	switch controller := item.(type) {
	case ControllerRegistration:
		if controller.Constructor == nil {
			return nil, "", "", fmt.Errorf("controller constructor cannot be nil")
		}

		return controller.Constructor, controller.BasePath, controller.Version, nil

	case *ControllerRegistration:
		if controller == nil {
			return nil, "", "", fmt.Errorf("controller registration cannot be nil")
		}

		if controller.Constructor == nil {
			return nil, "", "", fmt.Errorf("controller constructor cannot be nil")
		}

		return controller.Constructor, controller.BasePath, controller.Version, nil

	default:
		// Backward compatibility:
		// Controllers() []any { return []any{NewUserController} }
		return item, "", "", nil
	}
}
func buildVersionedBasePath(
	apiPrefix context.BasePath,
	moduleVersion context.APIVersion,
	controllerVersion context.APIVersion,
	moduleBasePath context.BasePath,
	controllerBasePath context.BasePath,
) string {
	version := moduleVersion

	if controllerVersion != context.VersionNeutral {
		version = controllerVersion
	}

	if version == context.VersionNeutral {
		return NormalizeRoute(
			string(moduleBasePath),
			string(controllerBasePath),
		)
	}

	return NormalizeRoute(
		string(apiPrefix),
		version.String(),
		string(moduleBasePath),
		string(controllerBasePath),
	)
}
func (a *App) registerController(moduleName string, basePath string, controllerValue reflect.Value) error {
	controllerName := controllerValue.Type().Name()

	controllerPointer := reflect.New(controllerValue.Type())
	controllerPointer.Elem().Set(controllerValue)

	controllerType := controllerPointer.Type()

	for i := 0; i < controllerPointer.NumMethod(); i++ {
		method := controllerPointer.Method(i)
		methodInfo := controllerType.Method(i)
		methodType := method.Type()

		if methodType.NumIn() != 0 {
			continue
		}

		if methodType.NumOut() != 3 {
			return fmt.Errorf(
				"gorix: controller method %s must return 3 values: gorix.Method, gorix.Path, gorix.RouteHandler",
				methodInfo.Name,
			)
		}

		startupOut := method.Call(nil)

		httpMethod, ok := startupOut[0].Interface().(context.Method)
		if !ok {
			return fmt.Errorf(
				"gorix: controller method %s first return value must be gorix.Method",
				methodInfo.Name,
			)
		}

		path, ok := startupOut[1].Interface().(context.Path)
		if !ok {
			return fmt.Errorf(
				"gorix: controller method %s second return value must be gorix.Path",
				methodInfo.Name,
			)
		}

		routeAction := startupOut[2].Interface()

		routeInvoker, err := newRouteActionInvoker(routeAction)
		if err != nil {
			return fmt.Errorf(
				"gorix: controller method %s invalid route handler: %w",
				methodInfo.Name,
				err,
			)
		}

		fullPath := NormalizeRoute(basePath, string(path))
		routeKey := string(httpMethod) + " " + fullPath

		if a.routes[routeKey] {
			return fmt.Errorf("gorix: duplicate route found: %s", routeKey)
		}

		a.routes[routeKey] = true

		a.routeInfos = append(a.routeInfos, context.RouteInfo{
			Method:     httpMethod,
			Path:       fullPath,
			Handler:    methodInfo.Name,
			Module:     moduleName,
			Controller: controllerName,
		})

		routeInterceptors := a.ResolveInterceptors(fullPath)

		routeHandler := func(c *context.Context) error {
			if c.R.Method != string(httpMethod) {
				a.handleException(&hook.ExceptionContext{
					Context:    c,
					Method:     httpMethod,
					Path:       fullPath,
					Module:     moduleName,
					Controller: controllerName,
					Handler:    methodInfo.Name,
					Error:      fmt.Errorf("method not allowed"),
					StatusCode: context.StatusMethodNotAllowed,
				})
				return nil
			}

			execCtx := &hook.ExecutionContext{
				Context:    c,
				Method:     httpMethod,
				Path:       fullPath,
				Module:     moduleName,
				Controller: controllerName,
				Handler:    methodInfo.Name,
				StartTime:  time.Now(),
			}

			defer func() {
				if rec := recover(); rec != nil {
					a.handleException(&hook.ExceptionContext{
						Context:    c,
						Method:     httpMethod,
						Path:       fullPath,
						Module:     moduleName,
						Controller: controllerName,
						Handler:    methodInfo.Name,
						Error:      fmt.Errorf("%v", rec),
						StatusCode: context.StatusInternalServerError,
					})
				}
			}()

			for _, interceptor := range routeInterceptors {
				if err := interceptor.Before(execCtx); err != nil {
					execCtx.Error = err

					a.handleException(&hook.ExceptionContext{
						Context:    c,
						Method:     httpMethod,
						Path:       fullPath,
						Module:     moduleName,
						Controller: controllerName,
						Handler:    methodInfo.Name,
						Error:      err,
						StatusCode: context.StatusInternalServerError,
					})
					return nil
				}
			}

			response, err := routeInvoker.Invoke(c)
			if err != nil {
				execCtx.Error = err

				a.handleException(&hook.ExceptionContext{
					Context:    c,
					Method:     httpMethod,
					Path:       fullPath,
					Module:     moduleName,
					Controller: controllerName,
					Handler:    methodInfo.Name,
					Error:      err,
					StatusCode: context.StatusInternalServerError,
				})
				return nil
			}

			execCtx.Response = response

			for j := len(routeInterceptors) - 1; j >= 0; j-- {
				if err := routeInterceptors[j].After(execCtx); err != nil {
					execCtx.Error = err

					a.handleException(&hook.ExceptionContext{
						Context:    c,
						Method:     httpMethod,
						Path:       fullPath,
						Module:     moduleName,
						Controller: controllerName,
						Handler:    methodInfo.Name,
						Error:      err,
						StatusCode: context.StatusInternalServerError,
					})
					return nil
				}
			}
			return c.Status(context.StatusCode(c.GetStatusOrDefault(context.StatusOK))).ResponseBodyInternal(access.Gorix, execCtx.Response)
		}

		routeMiddlewares := a.ResolveMiddlewares(fullPath)
		finalHandler := hook.ChainMiddlewares(routeHandler, routeMiddlewares...)

		a.routeEntries = append(a.routeEntries, routeEntry{
			Method:      httpMethod,
			Path:        fullPath,
			HandlerName: methodInfo.Name,
			Module:      moduleName,
			Controller:  controllerName,
			Handler:     finalHandler,
		})
	}

	return nil
}
func (a *App) Dispatch(w http.ResponseWriter, r *http.Request) {
	requestPath := strings.TrimSpace(r.URL.Path) // To avoid request failure due to whitespace
	var methodNotAllowedRoute *routeEntry
	var methodNotAllowedContext *context.Context
	methodNotAllowedScore := -1
	var bestRoute *routeEntry
	var bestContext *context.Context
	bestScore := -1
	//Register Routes based on Route Scoring
	for _, route := range a.routeEntries {
		matched, params, score := MatchRoute(route.Path, requestPath)
		if !matched {
			continue
		}

		c := context.NewContext(w, r)
		c.SetParams(params)

		if r.Method != string(route.Method) {
			if methodNotAllowedRoute == nil ||
				score > methodNotAllowedScore {
				methodNotAllowedRoute = &route
				methodNotAllowedContext = c
				methodNotAllowedScore = score
			}
			continue
		}

		if score > bestScore {
			bestRoute = &route
			bestContext = c
			bestScore = score
		}
	}

	if bestRoute != nil {
		if err := bestRoute.Handler(bestContext); err != nil {
			a.handleException(&hook.ExceptionContext{
				Context:    bestContext,
				Method:     bestRoute.Method,
				Path:       bestRoute.Path,
				Module:     bestRoute.Module,
				Controller: bestRoute.Controller,
				Handler:    bestRoute.HandlerName,
				Error:      err,
				StatusCode: context.StatusInternalServerError,
			})
		}
		return
	}

	c := context.NewContext(w, r)

	if methodNotAllowedRoute != nil {
		a.handleException(&hook.ExceptionContext{
			Context:    methodNotAllowedContext,
			Method:     methodNotAllowedRoute.Method,
			Path:       methodNotAllowedRoute.Path,
			Module:     methodNotAllowedRoute.Module,
			Controller: methodNotAllowedRoute.Controller,
			Handler:    methodNotAllowedRoute.HandlerName,
			Error:      fmt.Errorf("method not allowed"),
			StatusCode: context.StatusMethodNotAllowed,
		})
		return
	}

	a.handleException(&hook.ExceptionContext{
		Context:    c,
		Method:     context.Method(r.Method),
		Path:       requestPath,
		Error:      fmt.Errorf("route not found"),
		StatusCode: context.StatusNotFound,
	})
}

/*
Route priority is determined by scoring each path segment.

	 Static segments receive a higher score than dynamic parameters:
	  		/user/summary -> 10 + 10 = 20
		   	/user/:id     -> 10 + 1  = 11
	    Routes with higher scores are matched first, preventing a dynamic route
		such as /user/:id from incorrectly capturing /user/summary, regardless
*/
func MatchRoute(pattern string, actualPath string) (bool, map[string]string, int) {
	pattern = normalizeRoutePath(pattern)
	actualPath = normalizeRoutePath(actualPath)

	patternParts := SplitRoutePath(pattern)
	actualParts := SplitRoutePath(actualPath)

	if len(patternParts) != len(actualParts) {
		return false, nil, 0
	}

	params := make(map[string]string)
	score := 0

	for i := range patternParts {
		patternPart := patternParts[i]
		actualPart := actualParts[i]

		if strings.HasPrefix(patternPart, ":") {
			key := strings.TrimPrefix(patternPart, ":")
			params[key] = actualPart
			score += 1
			continue
		}

		if patternPart != actualPart {
			return false, nil, 0
		}

		score += 10
	}

	return true, params, score
}

func SplitRoutePath(path string) []string {
	path = normalizeRoutePath(path)

	if path == "/" {
		return []string{}
	}

	path = strings.Trim(path, "/")

	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func normalizeRoutePath(path string) string {
	if path == "" {
		return "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}

	return path
}
func (a *App) handleException(ctx *hook.ExceptionContext) {
	filters := a.ResolveFilters(ctx.Path)

	if len(filters) == 0 {
		_ = ctx.Context.Status(ctx.StatusCode).ResponseBodyInternal(access.Gorix,
			global.ErrorDTO{
				Success: false,
				Message: ctx.Error.Error(),
				Error:   ctx.Error.Error() + "on" + ctx.Path,
			})
		return
	}

	for _, filter := range filters {
		filter.Catch(ctx)
		return
	}
}

func (a *App) ResolveMiddlewares(path string) []hook.Middleware {
	result := make([]hook.Middleware, 0)

	for _, item := range a.middlewares {
		if item.Rule.Match(path) {
			result = append(result, item.Middleware)
		}
	}

	return result
}

func (a *App) ResolveInterceptors(path string) []hook.Interceptor {
	result := make([]hook.Interceptor, 0)

	for _, item := range a.interceptors {
		if item.Rule.Match(path) {
			result = append(result, item.Interceptor)
		}
	}

	return result
}

func (a *App) ResolveFilters(path string) []hook.Filter {
	result := make([]hook.Filter, 0)

	for _, item := range a.filters {
		if item.Rule.Match(path) {
			result = append(result, item.Filter)
		}
	}

	return result
}

func NormalizeRoute(parts ...string) string {
	cleanParts := make([]string, 0)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "/")

		if part == "" {
			continue
		}

		cleanParts = append(cleanParts, part)
	}

	if len(cleanParts) == 0 {
		return "/"
	}

	return "/" + strings.Join(cleanParts, "/")
}

func (a *App) PrintRoutes() {
	const (
		Orange = "\033[38;5;208m"
		Reset  = "\033[0m"
	)

	if len(a.routeInfos) == 0 {
		fmt.Println(Orange + "No routes registered" + Reset)
		return
	}

	fmt.Println()
	fmt.Println(Orange + "Registered Gorix Routes" + Reset)
	fmt.Println(Orange + "------------------------------------------------------------" + Reset)
	fmt.Printf(
		Orange+"%-8s %-30s %-20s %-20s\n"+Reset,
		"METHOD",
		"PATH",
		"CONTROLLER",
		"HANDLER",
	)
	fmt.Println(Orange + "------------------------------------------------------------" + Reset)

	for _, route := range a.routeInfos {
		fmt.Printf(
			Orange+"%-8s %-30s %-20s %-20s\n"+Reset,
			route.Method,
			route.Path,
			route.Controller,
			route.Handler,
		)
	}

	fmt.Println(Orange + "------------------------------------------------------------" + Reset)
	fmt.Println()
}

type routeActionInvoker struct {
	fn   reflect.Value
	args []routeActionArg
}

type routeActionArg struct {
	typ       reflect.Type
	source    context.BindSource
	isContext bool
}

var (
	contextPtrType = reflect.TypeOf((*context.Context)(nil))
	bindingArgType = reflect.TypeOf((*context.BindingArg)(nil)).Elem()
	errorType      = reflect.TypeOf((*error)(nil)).Elem()
)

func newRouteActionInvoker(handler any) (*routeActionInvoker, error) {
	fn := reflect.ValueOf(handler)

	if !fn.IsValid() || fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("route handler must be a function")
	}

	fnType := fn.Type()

	if fnType.NumOut() != 2 {
		return nil, fmt.Errorf("route handler must return 2 values: any, error")
	}

	if !fnType.Out(1).Implements(errorType) {
		return nil, fmt.Errorf("route handler second return value must be error")
	}

	invoker := &routeActionInvoker{
		fn:   fn,
		args: make([]routeActionArg, fnType.NumIn()),
	}

	bodyCount := 0

	for i := 0; i < fnType.NumIn(); i++ {
		argType := fnType.In(i)

		if argType == contextPtrType {
			invoker.args[i] = routeActionArg{
				typ:       argType,
				isContext: true,
			}
			continue
		}

		if argType.Implements(bindingArgType) {
			arg := reflect.New(argType).Elem()
			bindingArg := arg.Interface().(context.BindingArg)

			source := bindingArg.BindSource()

			if source == context.BindSourceBody {
				bodyCount++
				if bodyCount > 1 {
					return nil, fmt.Errorf("route handler cannot have more than one body argument")
				}
			}

			invoker.args[i] = routeActionArg{
				typ:    argType,
				source: source,
			}
			continue
		}

		return nil, fmt.Errorf(
			"unsupported route handler argument at index %d: %s",
			i,
			argType.String(),
		)
	}

	return invoker, nil
}

func (i *routeActionInvoker) Invoke(c *context.Context) (any, error) {
	args := make([]reflect.Value, len(i.args))

	for index, arg := range i.args {
		if arg.isContext {
			args[index] = reflect.ValueOf(c)
			continue
		}

		value, err := bindRouteActionArg(c, arg.typ, arg.source)
		if err != nil {
			return nil, err
		}

		args[index] = value
	}

	out := i.fn.Call(args)

	if !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}

	return out[0].Interface(), nil
}

func bindRouteActionArg(
	c *context.Context,
	argType reflect.Type,
	source context.BindSource,
) (reflect.Value, error) {
	arg := reflect.New(argType).Elem()

	valueField := arg.FieldByName("Value")
	if !valueField.IsValid() {
		return reflect.Value{}, fmt.Errorf("%s must have Value field", argType.String())
	}

	if !valueField.CanAddr() {
		return reflect.Value{}, fmt.Errorf("%s.Value cannot be addressed", argType.String())
	}

	target := valueField.Addr().Interface()

	switch source {
	case context.BindSourceParams:
		return arg, c.BindParams(target)

	case context.BindSourceQuery:
		return arg, c.BindQuery(target)

	case context.BindSourceBody:
		return arg, c.BindBody(target)

	case context.BindSourceHeaders:
		return arg, c.BindHeaders(target)

	default:
		return reflect.Value{}, fmt.Errorf("unsupported bind source: %s", source)
	}
}
