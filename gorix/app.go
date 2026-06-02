package gorix

import "github.com/Gromosome/gorix/gorix/engine"

type App = engine.App

func NewApp() *App {
	return engine.NewApp()
}
