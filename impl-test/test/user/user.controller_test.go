package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/impl-test/user/controller"
	"github.com/Gromosome/gorix/impl-test/user/repository"
)

func TestUserControllerFindByID(t *testing.T) {
	mockService := NewMockUserService()

	// promotionService is nil because this test does not call /promotions.
	controller := controller.NewUserController(
		mockService,
		nil,
	)

	method, path, handler := controller.FindByID()

	if method != "GET" {
		t.Fatalf("expected GET, got %s", method)
	}

	if path != "/:id" {
		t.Fatalf("expected /:id, got %s", path)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/user/10",
		nil,
	)
	res := httptest.NewRecorder()

	ctx := gorixcontext.NewContext(res, req)
	ctx.SetParams(map[string]string{
		"id": "10",
	})

	response, err := handler(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result, ok := response.(*repository.User)
	if !ok {
		t.Fatalf(
			"expected *user.User, got %T",
			response,
		)
	}

	if result.ID != 10 {
		t.Fatalf("expected ID 10, got %d", result.ID)
	}

	if result.Name != "Mock Controller User" {
		t.Fatalf(
			"expected mock user, got %s",
			result.Name,
		)
	}

	if !mockService.GetByIDCalled {
		t.Fatal("expected GetByID to be called")
	}

	if mockService.GetByIDValue != 10 {
		t.Fatalf(
			"expected GetByID id 10, got %d",
			mockService.GetByIDValue,
		)
	}
}
