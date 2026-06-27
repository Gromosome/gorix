package gorix

import "github.com/Gromosome/gorix/gorix/di"

type Scope = di.Scope

const (
	Singleton = di.Singleton
	Transient = di.Transient
	Request   = di.Request
)

type ProviderOption = di.ProviderOption
type ProviderRegistration = di.Registration

var As = di.As
var Named = di.Named
var WithScope = di.WithScope
var Replace = di.Replace

var Provider = di.Provider
var Instance = di.Instance
