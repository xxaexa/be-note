package note

import (
	"database/sql"
	"go-note/middlewares"
	"go-note/models"
	"go-note/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store models.NoteStore
}

func NewHandler(store models.NoteStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	noteRouter := router.PathPrefix("/notes").Subrouter()
	noteRouter.Use(middlewares.JWTMiddleware)

	noteRouter.HandleFunc("/", h.HandleGetNotes).Methods("GET")
	noteRouter.HandleFunc("/{id}", h.HandleGetNoteByID).Methods("GET")
	noteRouter.HandleFunc("/{id}", h.HandleUpdateNote).Methods("PUT")
	noteRouter.HandleFunc("/{id}", h.HandleDeleteNote).Methods("DELETE")

}

func (h *Handler) HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.NotePayload
	if err := utils.ParseJSON(r, &note); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	if err := utils.Validate.Struct(note); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, errors.Error(), false)
		return
	}

	err := h.store.CreateNote(&note)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, "create success", false)
}

func (h *Handler) HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetNotes()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "success", notes)
}

func (h *Handler) HandleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	notes, err := h.store.GetNoteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, http.StatusNotFound, "note not found", false)
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "success", notes)
}

func (h *Handler) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	var note models.NotePayload
	if err := utils.ParseJSON(r, &note); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	if err := utils.Validate.Struct(note); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, errors.Error(), false)
		return
	}

	err = h.store.UpdateNote(id, &note)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, http.StatusNotFound, "note not found", false)
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "update success", id)
}

func (h *Handler) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryID(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	err = h.store.DeleteNote(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, http.StatusNotFound, "note not found", false)
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "delete success", id)
}
