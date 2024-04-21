package auth

import (
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
	router.HandleFunc("/auth/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/auth/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.UserLoginPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid payload", errors.Error())
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid email or password", false)
		return
	}

	if !utils.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid email or password", false)
		return
	}

	secret := []byte(os.Getenv("SECRET_KEY"))
	token, err := middlewares.CreateJWT(secret, u.ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"email":    u.Email,
			"username": u.Username,
		},
	}
	utils.ResponseJSON(w, http.StatusOK, "success", response)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegisterPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ResponseJSON(w, http.StatusBadRequest, errors.Error(), false)
		return
	}

	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "username already exists", user.Username)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, "error", err)
		return
	}

	err = h.store.CreateUser(&models.UserRegisterPayload{
		Email:    user.Email,
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, "error", err)
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, "register successfully", false)
}
