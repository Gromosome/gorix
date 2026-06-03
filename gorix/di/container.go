package di

import (
	"fmt"
	"reflect"
)

type Container struct {
	providers map[reflect.Type]reflect.Value
	instances map[reflect.Type]reflect.Value
	resolving map[reflect.Type]bool
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]reflect.Value),
		instances: make(map[reflect.Type]reflect.Value),
		resolving: make(map[reflect.Type]bool),
	}
}

func (c *Container) RegisterProvider(provider any) error {
	if provider == nil {
		return fmt.Errorf("gorix di: provider cannot be nil")
	}

	providerValue := reflect.ValueOf(provider)
	providerType := providerValue.Type()

	if providerType.Kind() != reflect.Func {
		return fmt.Errorf("gorix di: provider must be constructor function, got %s", providerType.Kind())
	}

	if providerType.NumOut() != 1 {
		return fmt.Errorf("gorix di: provider constructor must return exactly one value")
	}

	resultType := providerType.Out(0)

	if resultType.Kind() != reflect.Pointer && resultType.Kind() != reflect.Struct && resultType.Kind() != reflect.Interface {
		return fmt.Errorf("gorix di: provider return type must be pointer, struct, or interface, got %s", resultType.Kind())
	}

	if _, exists := c.providers[resultType]; exists {
		return fmt.Errorf("gorix di: duplicate provider registered for %s", resultType.String())
	}

	c.providers[resultType] = providerValue

	return nil
}

func (c *Container) Resolve(targetType reflect.Type) (reflect.Value, error) {
	if instance, ok := c.instances[targetType]; ok {
		return instance, nil
	}

	provider, ok := c.providers[targetType]
	if !ok {
		return reflect.Value{}, fmt.Errorf("gorix di: no provider found for %s", targetType.String())
	}

	if c.resolving[targetType] {
		return reflect.Value{}, fmt.Errorf("gorix di: circular dependency detected for %s", targetType.String())
	}

	c.resolving[targetType] = true
	defer delete(c.resolving, targetType)

	providerType := provider.Type()

	args := make([]reflect.Value, 0, providerType.NumIn())

	for i := 0; i < providerType.NumIn(); i++ {
		depType := providerType.In(i)

		dep, err := c.Resolve(depType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(
				"gorix di: failed to resolve dependency %s for %s: %w",
				depType.String(),
				targetType.String(),
				err,
			)
		}

		args = append(args, dep)
	}

	results := provider.Call(args)

	instance := results[0]
	c.instances[targetType] = instance

	return instance, nil
}

func (c *Container) Build(constructor any) (reflect.Value, error) {
	constructorValue := reflect.ValueOf(constructor)
	constructorType := constructorValue.Type()

	if constructorType.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("gorix di: controller must be constructor function")
	}

	if constructorType.NumOut() != 1 {
		return reflect.Value{}, fmt.Errorf("gorix di: constructor must return exactly one value")
	}

	args := make([]reflect.Value, 0, constructorType.NumIn())

	for i := 0; i < constructorType.NumIn(); i++ {
		depType := constructorType.In(i)

		dep, err := c.Resolve(depType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(
				"gorix di: failed to resolve constructor dependency %s: %w",
				depType.String(),
				err,
			)
		}

		args = append(args, dep)
	}

	result := constructorValue.Call(args)[0]

	return result, nil
}
