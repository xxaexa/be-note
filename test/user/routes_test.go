package note

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-note/models"
	"go-note/service/auth"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := auth.NewHandler(userStore)

	t.Run("should fail creating a user if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/auth/register", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/auth/register", handler.HandleRegister).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle creating a user", func(t *testing.T) {
		payload := models.UserRegisterPayload{
			Email:    "test@mail.com",
			Username: "test",
			Password: "123456",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/auth/register", handler.HandleRegister).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

type mockUserStore struct{}

func (m *mockUserStore) CreateUser(user *models.UserRegisterPayload) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{}, nil
}

func (m *mockUserStore) GetUserByID(id int) (*models.User, error) {
	return &models.User{}, nil
}
