package engine

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/Gromosome/gorix/gorix/core"
	"github.com/Gromosome/gorix/gorix/hook"
)

type App struct {
	mux         *http.ServeMux
	routes      map[string]bool
	routeInfos  []core.RouteInfo
	projectRoot string
	config      core.Config

	middlewares  []hook.MiddlewareConfig
	interceptors []hook.InterceptorConfig
	filters      []hook.FilterConfig
}

func NewApp() *App {
	wd, err := os.Getwd()
	fmt.Println("Project root:", wd)
	if err != nil {
		wd = "."
	}
	config := core.LoadConfig(wd)
	return &App{
		mux:          http.NewServeMux(),
		routes:       make(map[string]bool),
		routeInfos:   make([]core.RouteInfo, 0),
		projectRoot:  wd,
		config:       config,
		middlewares:  make([]hook.MiddlewareConfig, 0),
		interceptors: make([]hook.InterceptorConfig, 0),
		filters:      make([]hook.FilterConfig, 0),
	}
}
func (a *App) Use(items ...any) {
	for _, item := range items {
		switch v := item.(type) {
		case hook.Middleware:
			a.middlewares = append(a.middlewares, hook.GlobalMiddleware(v))

		case hook.MiddlewareConfig:
			a.middlewares = append(a.middlewares, v)

		default:
			log.Fatalf("gorix: invalid middleware type %T", item)
		}
	}
}

func (a *App) UseInterceptors(items ...any) {
	for _, item := range items {
		switch v := item.(type) {
		case hook.Interceptor:
			a.interceptors = append(a.interceptors, hook.GlobalInterceptor(v))

		case hook.InterceptorConfig:
			a.interceptors = append(a.interceptors, v)

		default:
			log.Fatalf("gorix: invalid interceptor type %T", item)
		}
	}
}

func (a *App) UseFilters(items ...any) {
	for _, item := range items {
		switch v := item.(type) {
		case hook.Filter:
			a.filters = append(a.filters, hook.GlobalFilter(v))

		case hook.FilterConfig:
			a.filters = append(a.filters, v)

		default:
			log.Fatalf("gorix: invalid filter type %T", item)
		}
	}
}

func (a *App) Listen() {
	if err := a.TryListen(a.config.Address()); err != nil {
		log.Fatal(err)
	}
}
func (a *App) RegisterModules(modules ...any) {
	if err := a.TryRegisterModules(modules...); err != nil {
		log.Fatal(err)
	}
}
func (a *App) TryRegisterModules(modules ...any) error {
	for _, module := range modules {
		if module == nil {
			return fmt.Errorf("gorix: nil module cannot be registered")
		}

		if err := a.registerModule(module); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) TryListen(addr string) error {
	if !a.config.IsProd() {
		if err := ValidateProject(a.projectRoot); err != nil {
			return err
		}
		fmt.Println("Gorix validation passed")
	} else {
		fmt.Println("Gorix validation skipped: gorix.app.prod=true")
	}
	a.PrintRoutes()
	fmt.Println("Gorix server running on", addr)
	return http.ListenAndServe(addr, a.mux)
}

func (a *App) registerModule(module any) error {
	moduleValue := reflect.ValueOf(module)

	if moduleValue.Kind() != reflect.Pointer {
		return fmt.Errorf("gorix: module must be pointer, got %s", moduleValue.Kind())
	}

	moduleType := moduleValue.Type()

	for i := 0; i < moduleValue.NumMethod(); i++ {
		method := moduleValue.Method(i)
		methodInfo := moduleType.Method(i)
		methodType := method.Type()

		if methodType.NumIn() != 0 {
			continue
		}

		if methodType.NumOut() != 2 {
			continue
		}

		out := method.Call(nil)

		basePathValue := out[0]
		controllerValue := out[1]

		basePath, ok := basePathValue.Interface().(core.BasePath)
		if !ok {
			return fmt.Errorf(
				"gorix: module method %s must return gorix.BasePath as first return value",
				methodInfo.Name,
			)
		}

		if controllerValue.Kind() != reflect.Struct {
			return fmt.Errorf(
				"gorix: module method %s second return value must be controller struct",
				methodInfo.Name,
			)
		}

		moduleName := moduleType.Elem().Name()

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
				"gorix: controller method %s must return 3 values: gorix.Method, gorix.Path, response",
				methodInfo.Name,
			)
		}

		startupOut := method.Call(nil)

		httpMethod, ok := startupOut[0].Interface().(core.Method)
		if !ok {
			return fmt.Errorf(
				"gorix: controller method %s first return value must be gorix.Method",
				methodInfo.Name,
			)
		}

		path, ok := startupOut[1].Interface().(core.Path)
		if !ok {
			return fmt.Errorf(
				"gorix: controller method %s second return value must be gorix.Path",
				methodInfo.Name,
			)
		}

		fullPath := normalizeRoute(basePath, string(path))
		routeKey := string(httpMethod) + " " + fullPath

		if a.routes[routeKey] {
			return fmt.Errorf("gorix: duplicate route found: %s", routeKey)
		}

		a.routes[routeKey] = true

		a.routeInfos = append(a.routeInfos, core.RouteInfo{
			Method:     httpMethod,
			Path:       fullPath,
			Handler:    methodInfo.Name,
			Module:     moduleName,
			Controller: controllerName,
		})

		routeInterceptors := a.resolveInterceptors(fullPath)

		routeHandler := func(c *core.Context) error {
			if c.R.Method != string(httpMethod) {
				a.handleException(&hook.ExceptionContext{
					Context:    c,
					Method:     httpMethod,
					Path:       fullPath,
					Module:     moduleName,
					Controller: controllerName,
					Handler:    methodInfo.Name,
					Error:      fmt.Errorf("method not allowed"),
					StatusCode: core.StatusMethodNotAllowed,
				})
				return nil
			}

			ctx := &hook.ExecutionContext{
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
						StatusCode: core.StatusInternalServerError,
					})
				}
			}()

			for _, interceptor := range routeInterceptors {
				if err := interceptor.Before(ctx); err != nil {
					ctx.Error = err

					a.handleException(&hook.ExceptionContext{
						Context:    c,
						Method:     httpMethod,
						Path:       fullPath,
						Module:     moduleName,
						Controller: controllerName,
						Handler:    methodInfo.Name,
						Error:      err,
						StatusCode: core.StatusInternalServerError,
					})
					return nil
				}
			}

			result := method.Call(nil)
			response := result[2].Interface()
			ctx.Response = response

			for j := len(routeInterceptors) - 1; j >= 0; j-- {
				if err := routeInterceptors[j].After(ctx); err != nil {
					ctx.Error = err

					a.handleException(&hook.ExceptionContext{
						Context:    c,
						Method:     httpMethod,
						Path:       fullPath,
						Module:     moduleName,
						Controller: controllerName,
						Handler:    methodInfo.Name,
						Error:      err,
						StatusCode: core.StatusInternalServerError,
					})
					return nil
				}
			}

			return c.Status(core.StatusOK).JSON(ctx.Response)
		}
		routeMiddlewares := a.resolveMiddlewares(fullPath)
		finalHandler := hook.ChainMiddlewares(routeHandler, routeMiddlewares...)

		a.mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
			c := core.NewContext(w, r)

			if err := finalHandler(c); err != nil {
				a.handleException(&hook.ExceptionContext{
					Context:    c,
					Method:     httpMethod,
					Path:       fullPath,
					Module:     moduleName,
					Controller: controllerName,
					Handler:    methodInfo.Name,
					Error:      err,
					StatusCode: core.StatusInternalServerError,
				})
			}
		})
	}
	return nil
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
