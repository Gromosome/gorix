package user

import (
	"testing"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user"
)

func TestUserServiceWithMockRepository(t *testing.T) {
	app := gorix.NewApp()

	err := app.OverrideProvider(
		NewMockUserRepository,
		gorix.As((*user.UserRepositoryPort)(nil)),
	)
	if err != nil {
		t.Fatal(err)
	}

	app.RegisterModules(
		promotion.NewPromotionModule(),
		user.NewUserModule(),
	)

	var service user.UserServicePort
	if err := app.Resolve(&service); err != nil {
		t.Fatal(err)
	}

	result, err := service.GetByID(nil, 10)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "Mock User" {
		t.Fatalf("expected mock user, got %s", result.Name)
	}
}
