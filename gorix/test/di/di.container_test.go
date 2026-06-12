package di

import (
	"reflect"
	"testing"

	di2 "github.com/Gromosome/gorix/gorix/di"
)

type testDependency struct {
	Value string
}

type testService struct {
	Dependency *testDependency
}

type circularFirst struct {
	Second *circularSecond
}

type circularSecond struct {
	First *circularFirst
}

func TestContainerRegistersAndResolvesProviders(t *testing.T) {
	container := di2.NewContainer()

	if err := container.RegisterProvider(func() *testDependency {
		return &testDependency{Value: "ok"}
	}); err != nil {
		t.Fatalf("RegisterProvider dependency returned error: %v", err)
	}
	if err := container.RegisterProvider(func(dep *testDependency) *testService {
		return &testService{Dependency: dep}
	}); err != nil {
		t.Fatalf("RegisterProvider service returned error: %v", err)
	}

	value, err := container.Resolve(reflect.TypeOf(&testService{}))
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	service := value.Interface().(*testService)
	if service.Dependency.Value != "ok" {
		t.Fatalf("unexpected dependency value: %s", service.Dependency.Value)
	}
}

func TestContainerBuildInjectsDependencies(t *testing.T) {
	container := di2.NewContainer()
	dependency := &testDependency{Value: "instance"}
	if err := container.RegisterInstance(dependency); err != nil {
		t.Fatalf("RegisterInstance returned error: %v", err)
	}

	value, err := container.Build(func(dep *testDependency) *testService {
		return &testService{Dependency: dep}
	})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if value.Interface().(*testService).Dependency != dependency {
		t.Fatal("Build did not inject registered instance")
	}
}

func TestContainerRejectsInvalidRegistrations(t *testing.T) {
	container := di2.NewContainer()

	if err := container.RegisterProvider(nil); err == nil {
		t.Fatal("nil provider should be rejected")
	}
	if err := container.RegisterProvider("not func"); err == nil {
		t.Fatal("non-function provider should be rejected")
	}
	if err := container.RegisterProvider(func() {}); err == nil {
		t.Fatal("provider without return should be rejected")
	}
	if err := container.RegisterInstance(nil); err == nil {
		t.Fatal("nil instance should be rejected")
	}
}

func TestContainerDetectsCircularDependencies(t *testing.T) {
	container := di2.NewContainer()
	_ = container.RegisterProvider(func(s *circularSecond) *circularFirst {
		return &circularFirst{Second: s}
	})
	_ = container.RegisterProvider(func(f *circularFirst) *circularSecond {
		return &circularSecond{First: f}
	})

	if _, err := container.Resolve(reflect.TypeOf(&circularFirst{})); err == nil {
		t.Fatal("expected circular dependency error")
	}
}
