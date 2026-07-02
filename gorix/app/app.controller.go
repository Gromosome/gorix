package app

import "github.com/Gromosome/gorix/gorix/core/context"

type versionedModule interface {
	APIVersion() context.APIVersion
}
type ControllerRegistration struct {
	Constructor any
	BasePath    context.BasePath
	Version     context.APIVersion
}

func Controller(constructor any, args ...any) ControllerRegistration {
	registration := ControllerRegistration{
		Constructor: constructor,
		BasePath:    "",
		Version:     context.VersionNeutral,
	}

	for _, arg := range args {
		switch value := arg.(type) {
		case string:
			registration.BasePath = context.BasePath(value)

		case context.BasePath:
			registration.BasePath = value

		case context.APIVersion:
			registration.Version = value
		}
	}

	return registration
}
