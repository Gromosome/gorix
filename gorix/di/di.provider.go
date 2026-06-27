package di

import (
	"fmt"
	"reflect"
)

type Scope string

const (
	Singleton Scope = "singleton"
	Transient Scope = "transient"
	Request   Scope = "request" // future: request-scoped cache
)

type Key struct {
	Type reflect.Type
	Name string
}

type ProviderDefinition struct {
	Key         Key
	ResultType  reflect.Type
	Constructor reflect.Value

	Scope       Scope
	Replace     bool
	IsOverride  bool
	Instance    reflect.Value
	HasInstance bool

	optionErrors []error
}

type ProviderOption func(*ProviderDefinition)

func As(typeHint any) ProviderOption {
	return func(def *ProviderDefinition) {
		t, err := typeFromHint(typeHint)
		if err != nil {
			def.optionErrors = append(def.optionErrors, err)
			return
		}
		def.Key.Type = t
	}
}

func Named(name string) ProviderOption {
	return func(def *ProviderDefinition) {
		def.Key.Name = name
	}
}

func WithScope(scope Scope) ProviderOption {
	return func(def *ProviderDefinition) {
		def.Scope = scope
	}
}

func Replace() ProviderOption {
	return func(def *ProviderDefinition) {
		def.Replace = true
	}
}

func typeFromHint(typeHint any) (reflect.Type, error) {
	if typeHint == nil {
		return nil, fmt.Errorf("gorix di: type hint cannot be nil")
	}

	t := reflect.TypeOf(typeHint)

	// Common Go interface-binding pattern:
	// As((*UserRepository)(nil))
	if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Interface {
		return t.Elem(), nil
	}

	return t, nil
}

func validateAssignable(actual reflect.Type, target reflect.Type) error {
	if actual.AssignableTo(target) {
		return nil
	}

	if target.Kind() == reflect.Interface && actual.Implements(target) {
		return nil
	}

	return fmt.Errorf(
		"gorix di: %s is not assignable to %s",
		actual.String(),
		target.String(),
	)
}
