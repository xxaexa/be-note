package note

import (
	"go-note/models"
	"go-note/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store models.NoteStore
}

func NewHandler(store models.NoteStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/note", h.handleCreateNote).Methods("POST")
	router.HandleFunc("/notes", h.handleGetNotes).Methods("GET")
	router.HandleFunc("/note/{id}", h.handleGetNoteByID).Methods("GET")
	router.HandleFunc("/note", h.handleUpdateNote).Methods("PUT")
	router.HandleFunc("/note/{id}", h.handleDeleteNote).Methods("DELETE")
}

func (h *Handler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.CreateNotePayload
	if err := utils.ParseJSON(r, &note); err != nil {
		response := map[string]string{"message": err.Error()}
		utils.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	err := h.store.CreateNote(&note)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"message": "success"}
	utils.ResponseJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetNotes()

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, notes)
}

func (h *Handler) handleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	notes, err := h.store.GetNoteByID(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, notes)
}

func (h *Handler) handleUpdateNote(w http.ResponseWriter, r *http.Request) {

	var note models.Note
	if err := utils.ParseJSON(r, &note); err != nil {
		response := map[string]string{"message": err.Error()}
		utils.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	err := h.store.UpdateNote(&note)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"message": "success"}
	utils.ResponseJSON(w, http.StatusOK, response)
}

func (h *Handler) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.DeleteNote(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"message": "success"}
	utils.ResponseJSON(w, http.StatusOK, response)
}
