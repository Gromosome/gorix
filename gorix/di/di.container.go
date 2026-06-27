package di

import (
	"fmt"
	"reflect"
	"sync"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type Container struct {
	mu        sync.RWMutex
	providers map[Key]*ProviderDefinition
	resolving map[Key]bool
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[Key]*ProviderDefinition),
		resolving: make(map[Key]bool),
	}
}

func (c *Container) RegisterProvider(
	provider any,
	opts ...ProviderOption,
) error {
	return c.registerProvider(
		provider,
		false,
		false, // isOverride
		opts...,
	)
}

func (c *Container) OverrideProvider(
	provider any,
	opts ...ProviderOption,
) error {
	opts = append(opts, Replace())

	return c.registerProvider(
		provider,
		true,
		true, // isOverride
		opts...,
	)
}

func (c *Container) RegisterInstance(
	instance any,
	opts ...ProviderOption,
) error {
	return c.registerInstance(
		instance,
		false,
		false,
		opts...,
	)
}

func (c *Container) OverrideInstance(
	instance any,
	opts ...ProviderOption,
) error {
	opts = append(opts, Replace())

	return c.registerInstance(
		instance,
		true,
		true,
		opts...,
	)
}

func (c *Container) registerProvider(provider any, forceReplace bool, isOverride bool, opts ...ProviderOption) error {
	if provider == nil {
		return fmt.Errorf("gorix di: provider cannot be nil")
	}

	providerValue := reflect.ValueOf(provider)
	providerType := providerValue.Type()

	if providerType.Kind() != reflect.Func {
		return fmt.Errorf(
			"gorix di: provider must be constructor function, got %s",
			providerType.Kind(),
		)
	}

	if providerType.NumOut() != 1 && providerType.NumOut() != 2 {
		return fmt.Errorf(
			"gorix di: provider must return value or value,error",
		)
	}

	if providerType.NumOut() == 2 && !providerType.Out(1).Implements(errorType) {
		return fmt.Errorf(
			"gorix di: second provider return value must be error",
		)
	}

	resultType := providerType.Out(0)

	def := &ProviderDefinition{
		Key: Key{
			Type: resultType,
		},
		ResultType:  resultType,
		Constructor: providerValue,
		Scope:       Singleton,
		IsOverride:  isOverride,
	}

	for _, opt := range opts {
		opt(def)
	}

	for _, err := range def.optionErrors {
		if err != nil {
			return err
		}
	}

	if def.Key.Type == nil {
		return fmt.Errorf("gorix di: provider key type cannot be nil")
	}

	if err := validateAssignable(resultType, def.Key.Type); err != nil {
		return err
	}

	replace := forceReplace || def.Replace

	c.mu.Lock()
	defer c.mu.Unlock()

	existing, exists := c.providers[def.Key]

	if exists {
		if existing.IsOverride && !isOverride {
			return nil
		}

		if isOverride || replace {
			c.providers[def.Key] = def
			return nil
		}

		return fmt.Errorf(
			"gorix di: duplicate provider registered for %s",
			def.Key.Type.String(),
		)
	}

	c.providers[def.Key] = def
	return nil
}

func (c *Container) registerInstance(instance any, forceReplace bool, isOverride bool, opts ...ProviderOption) error {
	if instance == nil {
		return fmt.Errorf("gorix di: instance cannot be nil")
	}

	instanceValue := reflect.ValueOf(instance)
	instanceType := instanceValue.Type()

	def := &ProviderDefinition{
		Key: Key{
			Type: instanceType,
		},
		ResultType:  instanceType,
		Scope:       Singleton,
		Instance:    instanceValue,
		HasInstance: true,
		IsOverride:  isOverride,
	}

	for _, opt := range opts {
		opt(def)
	}

	for _, err := range def.optionErrors {
		if err != nil {
			return err
		}
	}

	if err := validateAssignable(instanceType, def.Key.Type); err != nil {
		return err
	}

	replace := forceReplace || def.Replace

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[def.Key]; exists && !replace {
		return fmt.Errorf(
			"gorix di: duplicate instance registered for %s",
			def.Key.Type.String(),
		)
	}

	c.providers[def.Key] = def
	return nil
}

func (c *Container) Resolve(targetType reflect.Type) (reflect.Value, error) {
	return c.resolveKey(Key{Type: targetType})
}

func (c *Container) ResolveNamed(targetType reflect.Type, name string) (reflect.Value, error) {
	return c.resolveKey(Key{
		Type: targetType,
		Name: name,
	})
}

func (c *Container) ResolveInto(target any) error {
	if target == nil {
		return fmt.Errorf("gorix di: resolve target cannot be nil")
	}

	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer {
		return fmt.Errorf("gorix di: resolve target must be pointer")
	}

	elem := targetValue.Elem()

	if !elem.CanSet() {
		return fmt.Errorf("gorix di: resolve target is not settable")
	}

	value, err := c.Resolve(elem.Type())
	if err != nil {
		return err
	}

	elem.Set(value)
	return nil
}

func (c *Container) resolveKey(key Key) (reflect.Value, error) {
	c.mu.RLock()
	def, ok := c.providers[key]
	if !ok {
		c.mu.RUnlock()
		return reflect.Value{}, fmt.Errorf(
			"gorix di: no provider found for %s",
			key.Type.String(),
		)
	}

	if def.Scope == Singleton && def.HasInstance {
		instance := def.Instance
		c.mu.RUnlock()
		return instance, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	if c.resolving[key] {
		c.mu.Unlock()
		return reflect.Value{}, fmt.Errorf(
			"gorix di: circular dependency detected for %s",
			key.Type.String(),
		)
	}
	c.resolving[key] = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.resolving, key)
		c.mu.Unlock()
	}()

	if !def.Constructor.IsValid() {
		return reflect.Value{}, fmt.Errorf(
			"gorix di: provider for %s has no constructor",
			key.Type.String(),
		)
	}

	instance, err := c.callConstructor(def.Constructor)
	if err != nil {
		return reflect.Value{}, err
	}

	if err := validateAssignable(instance.Type(), key.Type); err != nil {
		return reflect.Value{}, err
	}

	if def.Scope == Singleton {
		c.mu.Lock()
		def.Instance = instance
		def.HasInstance = true
		c.mu.Unlock()
	}

	return instance, nil
}

func (c *Container) Build(constructor any) (reflect.Value, error) {
	if constructor == nil {
		return reflect.Value{}, fmt.Errorf("gorix di: constructor cannot be nil")
	}

	constructorValue := reflect.ValueOf(constructor)
	constructorType := constructorValue.Type()

	if constructorType.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf(
			"gorix di: constructor must be function",
		)
	}

	if constructorType.NumOut() != 1 && constructorType.NumOut() != 2 {
		return reflect.Value{}, fmt.Errorf(
			"gorix di: constructor must return value or value,error",
		)
	}

	if constructorType.NumOut() == 2 && !constructorType.Out(1).Implements(errorType) {
		return reflect.Value{}, fmt.Errorf(
			"gorix di: second constructor return value must be error",
		)
	}

	return c.callConstructor(constructorValue)
}

func (c *Container) callConstructor(constructor reflect.Value) (reflect.Value, error) {
	constructorType := constructor.Type()

	args := make([]reflect.Value, 0, constructorType.NumIn())

	for i := 0; i < constructorType.NumIn(); i++ {
		depType := constructorType.In(i)

		dep, err := c.Resolve(depType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(
				"gorix di: failed to resolve dependency %s: %w",
				depType.String(),
				err,
			)
		}

		args = append(args, dep)
	}

	results := constructor.Call(args)

	if len(results) == 2 {
		errValue := results[1]
		if !errValue.IsNil() {
			return reflect.Value{}, errValue.Interface().(error)
		}
	}

	return results[0], nil
}
