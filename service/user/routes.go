package auth

import (
	"go-note/models"
	"go-note/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(store models.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user", h.handleGetUser).Methods("GET")
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.ResponseJSON(w, http.StatusBadRequest, "missing user ID", false)
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid user ID", false)
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, "error", err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "success", user)
}
