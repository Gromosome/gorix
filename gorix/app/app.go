package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Gromosome/gorix/gorix/app/linter"
	"github.com/Gromosome/gorix/gorix/config"
	"github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/core/database"
	"github.com/Gromosome/gorix/gorix/di"
	"github.com/Gromosome/gorix/gorix/hook"
)

type routeEntry struct {
	Method      context.Method
	Path        string
	HandlerName string
	Module      string
	Controller  string
	Handler     hook.Handler
}

type basePathModule interface {
	BasePath() context.BasePath
}

type providerModule interface {
	Providers() []any
}

type controllerModule interface {
	Controllers() []any
}
type App struct {
	routes       map[string]bool
	routeInfos   []context.RouteInfo
	routeEntries []routeEntry

	projectRoot string
	config      config.Config
	container   *di.Container

	middlewares  []hook.MiddlewareConfig
	interceptors []hook.InterceptorConfig
	filters      []hook.FilterConfig
	databases    *database.Manager
}

func NewApp() *App {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	cfg := config.LoadConfig(wd)

	container := di.NewContainer()
	databaseManager := database.NewManager()

	if err := container.RegisterInstance(databaseManager); err != nil {
		panic(fmt.Errorf(
			"gorix: failed to register database manager: %w",
			err,
		))
	}

	return &App{
		routes:       make(map[string]bool),
		routeInfos:   make([]context.RouteInfo, 0),
		routeEntries: make([]routeEntry, 0),

		projectRoot: wd,
		config:      cfg,
		container:   container,
		databases:   databaseManager,

		middlewares:  make([]hook.MiddlewareConfig, 0),
		interceptors: make([]hook.InterceptorConfig, 0),
		filters:      make([]hook.FilterConfig, 0),
	}
}
func (a *App) RouteEntries() []routeEntry {
	return a.routeEntries
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

		moduleValue := reflect.ValueOf(module)
		if moduleValue.Kind() != reflect.Pointer {
			return fmt.Errorf("gorix: module must be pointer, got %s", moduleValue.Kind())
		}

		moduleType := moduleValue.Type()
		moduleName := moduleType.Elem().Name()

		if providersModule, ok := module.(providerModule); ok {
			for _, provider := range providersModule.Providers() {
				if err := a.container.RegisterProvider(provider); err != nil {
					return fmt.Errorf("gorix: module %s provider registration failed: %w", moduleName, err)
				}
			}
		}
	}

	for _, module := range modules {
		if err := a.registerModuleControllers(module); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) TryListen(addr string) error {
	printGorixLogo()
	if !a.config.IsProd() {
		if err := linter.ValidateProject(a.projectRoot); err != nil {
			return err
		}

		fmt.Println("Gorix validation passed")
	}

	startupContext, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	if err := a.connectDatabases(startupContext); err != nil {
		return err
	}

	defer func() {
		if err := a.databases.Close(); err != nil {
			log.Printf(
				"gorix: database shutdown error: %v",
				err,
			)
		}
	}()

	a.PrintRoutes()

	fmt.Println("Gorix server running on", addr)

	return http.ListenAndServe(
		addr,
		http.HandlerFunc(a.Dispatch),
	)
}
func printGorixLogo() {
	// ANSI 24-bit TrueColor sequences for a perfectly smooth 6-step gradient
	const t1 = "\033[38;2;255;170;0m" // Bright golden orange (Top)
	const t2 = "\033[38;2;255;145;0m" // Light orange
	const t3 = "\033[38;2;245;115;0m" // Medium orange
	const t4 = "\033[38;2;230;85;0m"  // Rich orange
	const t5 = "\033[38;2;210;55;0m"  // Deep orange
	const t6 = "\033[38;2;180;35;0m"  // Dark burnished shadow (Bottom)
	const reset = "\033[0m"

	// Bold, clean, symmetrical geometric typography
	fmt.Print(t1)
	fmt.Println(`   ▄██████▄       ▄██████▄     ████████▄     ███    ███      ███ `)
	fmt.Print(t2)
	fmt.Println(`  ███    ███     ███    ███    ███    ███    ███      ███  ███   `)
	fmt.Print(t3)
	fmt.Println(`  ███            ███    ███    ████████▀     ███        ████     `)
	fmt.Print(t4)
	fmt.Println(`  ███  █████     ███    ███    ███   ███▄    ███        ████     `)
	fmt.Print(t5)
	fmt.Println(`  ███    ███     ███    ███    ███    ███    ███      ███  ███   `)
	fmt.Print(t6)
	fmt.Println(`   ▀██████▀       ▀██████▀     ███    ███    ███    ███      ███ `)

	fmt.Print(reset)
}
