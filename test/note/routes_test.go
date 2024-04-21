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

	t.Run("should fail creating a note if the payload is missing", func(t *testing.T) {
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

	t.Run("should handle creating a note", func(t *testing.T) {
		payload := models.NotePayload{
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

	t.Run("should fail updating a note if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/notes/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes/{id}", handler.HandleUpdateNote).Methods("PUT")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle updating a note", func(t *testing.T) {
		payload := models.NotePayload{
			Title:       "test",
			Description: "test description",
			UserID:      1,
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, "/notes/1", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes/{id}", handler.HandleUpdateNote).Methods(http.MethodPut)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle deleting a note", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodDelete, "/notes/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/notes/{id}", handler.HandleDeleteNote).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockNoteStore struct{}

func (m *mockNoteStore) CreateNote(note *models.NotePayload) error {
	return nil
}

func (m *mockNoteStore) GetNotes() ([]*models.Note, error) {
	return []*models.Note{}, nil
}

func (m *mockNoteStore) GetNoteByID(id int) (*models.Note, error) {
	return &models.Note{}, nil
}

func (m *mockNoteStore) UpdateNote(id int, note *models.NotePayload) error {
	return nil
}

func (m *mockNoteStore) DeleteNote(id int) error {
	return nil
}
