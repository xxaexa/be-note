package models

type UserStore interface {
	CreateUser(user *UserPayload) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int) (*User, error)
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
