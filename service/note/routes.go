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
	router.HandleFunc("/notes", h.HandleCreateNote).Methods("POST")
	router.HandleFunc("/notes", h.HandleGetNotes).Methods("GET")
	router.HandleFunc("/notes/{id}", h.HandleGetNoteByID).Methods("GET")
	router.HandleFunc("/notes", h.HandleUpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", h.HandleDeleteNote).Methods("DELETE")
}

func (h *Handler) HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.CreateNotePayload
	if err := utils.ParseJSON(r, &note); err != nil {
		response := map[string]string{"message": err.Error()}
		utils.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	err := h.store.CreateNote(&note)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"message": "success"}
	utils.ResponseJSON(w, http.StatusCreated, response)
}

func (h *Handler) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetNotes()

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, notes)
}

func (h *Handler) HandleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	notes, err := h.store.GetNoteByID(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, notes)
}

func (h *Handler) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {

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

func (h *Handler) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
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
