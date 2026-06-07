package app

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/hook"
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

	for _, controllerConstructor := range controllersModule.Controllers() {
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

		if err := a.registerController(moduleName, string(basePath), controllerValue); err != nil {
			return err
		}
	}

	return nil
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

		routeAction, ok := startupOut[2].Interface().(context.RouteHandler)
		if !ok {
			return fmt.Errorf(
				"gorix: controller method %s third return value must be gorix.RouteHandler",
				methodInfo.Name,
			)
		}

		fullPath := normalizeRoute(basePath, string(path))
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

		routeInterceptors := a.resolveInterceptors(fullPath)

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

			response, err := routeAction(c)
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

			return c.Status(context.StatusOK).JSON(execCtx.Response)
		}

		routeMiddlewares := a.resolveMiddlewares(fullPath)
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
func (a *App) dispatch(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path

	for _, route := range a.routeEntries {
		matched, params := matchRoute(route.Path, requestPath)
		if !matched {
			continue
		}

		c := context.NewContext(w, r)
		c.SetParams(params)

		if r.Method != string(route.Method) {
			a.handleException(&hook.ExceptionContext{
				Context:    c,
				Method:     route.Method,
				Path:       route.Path,
				Module:     route.Module,
				Controller: route.Controller,
				Handler:    route.HandlerName,
				Error:      fmt.Errorf("method not allowed"),
				StatusCode: context.StatusMethodNotAllowed,
			})
			return
		}

		if err := route.Handler(c); err != nil {
			a.handleException(&hook.ExceptionContext{
				Context:    c,
				Method:     route.Method,
				Path:       route.Path,
				Module:     route.Module,
				Controller: route.Controller,
				Handler:    route.HandlerName,
				Error:      err,
				StatusCode: context.StatusInternalServerError,
			})
			return
		}

		return
	}

	c := context.NewContext(w, r)

	a.handleException(&hook.ExceptionContext{
		Context:    c,
		Method:     context.Method(r.Method),
		Path:       requestPath,
		Error:      fmt.Errorf("route not found"),
		StatusCode: context.StatusNotFound,
	})
}
func matchRoute(pattern string, actualPath string) (bool, map[string]string) {
	pattern = normalizeRoutePath(pattern)
	actualPath = normalizeRoutePath(actualPath)

	patternParts := splitRoutePath(pattern)
	actualParts := splitRoutePath(actualPath)

	if len(patternParts) != len(actualParts) {
		return false, nil
	}

	params := make(map[string]string)

	for i := range patternParts {
		patternPart := patternParts[i]
		actualPart := actualParts[i]

		if strings.HasPrefix(patternPart, ":") {
			key := strings.TrimPrefix(patternPart, ":")
			params[key] = actualPart
			continue
		}

		if patternPart != actualPart {
			return false, nil
		}
	}

	return true, params
}

func splitRoutePath(path string) []string {
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
	filters := a.resolveFilters(ctx.Path)

	if len(filters) == 0 {
		_ = ctx.Context.Status(ctx.StatusCode).JSON(map[string]any{
			"success": false,
			"error":   ctx.Error.Error(),
			"path":    ctx.Path,
		})
		return
	}

	for _, filter := range filters {
		filter.Catch(ctx)
		return
	}
}

func (a *App) resolveMiddlewares(path string) []hook.Middleware {
	result := make([]hook.Middleware, 0)

	for _, item := range a.middlewares {
		if item.Rule.Match(path) {
			result = append(result, item.Middleware)
		}
	}

	return result
}

func (a *App) resolveInterceptors(path string) []hook.Interceptor {
	result := make([]hook.Interceptor, 0)

	for _, item := range a.interceptors {
		if item.Rule.Match(path) {
			result = append(result, item.Interceptor)
		}
	}

	return result
}

func (a *App) resolveFilters(path string) []hook.Filter {
	result := make([]hook.Filter, 0)

	for _, item := range a.filters {
		if item.Rule.Match(path) {
			result = append(result, item.Filter)
		}
	}

	return result
}

func normalizeRoute(basePath string, path string) string {
	basePath = "/" + strings.Trim(basePath, "/")
	path = "/" + strings.Trim(path, "/")

	if path == "/" {
		return basePath
	}

	return basePath + path
}

func (a *App) PrintRoutes() {
	if len(a.routeInfos) == 0 {
		fmt.Println("No routes registered")
		return
	}

	fmt.Println()
	fmt.Println("Registered Gorix Routes")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%-8s %-30s %-20s %-20s\n", "METHOD", "PATH", "CONTROLLER", "HANDLER")
	fmt.Println("------------------------------------------------------------")

	for _, route := range a.routeInfos {
		fmt.Printf(
			"%-8s %-30s %-20s %-20s\n",
			route.Method,
			route.Path,
			route.Controller,
			route.Handler,
		)
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println()
}
