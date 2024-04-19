package auth

import (
	"fmt"
	"go-note/middlewares"
	"go-note/models"
	"go-note/utils"
	"net/http"
	"os"

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
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/auth/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.UserLoginPayload
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

	u, err := h.store.GetUserByEmail(user.Email)
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

	secret := []byte(os.Getenv("SECRET_KEY"))
	token, err := middlewares.CreateJWT(secret, u.ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}
	response := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"email":    u.Email,
			"username": u.Username,
		},
	}
	utils.ResponseJSON(w, http.StatusOK, response)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegisterPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.ResponseJSON(w, http.StatusBadRequest, fmt.Errorf("user with Username %s already exists", user.Username))
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(&models.UserRegisterPayload{
		Email:    user.Email,
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]string{"message": "register successfully"}
	utils.ResponseJSON(w, http.StatusCreated, response)
}
