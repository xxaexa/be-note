package auth

import (
	"fmt"
	"go-note/middlewares"
	"go-note/models"
	"go-note/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(store models.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user", h.handleGetUser).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.UserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		response := map[string]string{"message": errors.Error()}
		utils.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	u, err := h.store.GetUserByUsername(user.Username)
	if err != nil {
		response := map[string]string{"message": "invalid username or password"}
		utils.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if !utils.ComparePasswords(u.Password, []byte(user.Password)) {
		response := map[string]string{"message": "invalid username or password"}
		utils.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	secret := []byte(os.Getenv("SCRT_KEY"))
	token, err := middlewares.CreateJWT(secret, u.ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.UserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByUsername(user.Username)
	if err == nil {
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("user with Username %s already exists", user.Username))
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(&models.UserPayload{
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, user)
}
