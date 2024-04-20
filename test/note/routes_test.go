package note

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-note/models"
	"go-note/service/note"

	"github.com/gorilla/mux"
)

func TestNoteServiceHandlers(t *testing.T) {
	noteStore := &mockNoteStore{}
	handler := note.NewHandler(noteStore)

	t.Run("should handle get notes", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/notes", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes", handler.HandleGetNotes).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail if the note ID is not a number", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/notes/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes/{id}", handler.HandleGetNoteByID).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle get note by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/notes/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes/{id}", handler.HandleGetNoteByID).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail creating a product if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/notes", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes", handler.HandleCreateNote).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle creating a product", func(t *testing.T) {
		payload := models.CreateNotePayload{
			Title:       "test",
			Description: "test description",
			UserID:      1,
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes", handler.HandleCreateNote).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockNoteStore struct{}

func (m *mockNoteStore) CreateNote(note *models.CreateNotePayload) error {
	return nil
}

func (m *mockNoteStore) GetNotes() ([]*models.Note, error) {
	return []*models.Note{}, nil
}

func (m *mockNoteStore) GetNoteByID(id int) (*models.Note, error) {
	return &models.Note{}, nil
}

func (m *mockNoteStore) UpdateNote(note *models.Note) error {
	return nil
}

func (m *mockNoteStore) DeleteNote(id int) error {
	return nil
}
