package test

import (
	"reflect"
	"testing"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/gorix/app"
)

func TestNewAppReturnsAppInstance(t *testing.T) {
	instance := gorix.NewApp()
	if instance == nil {
		t.Fatal("NewApp returned nil")
	}

	if reflect.TypeOf(instance) != reflect.TypeOf(app.NewApp()) {
		t.Fatalf("NewApp returned unexpected type %T", instance)
	}
}
