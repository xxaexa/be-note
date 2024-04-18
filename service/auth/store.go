package auth

import (
	"database/sql"
	"fmt"
	"go-note/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user *models.UserPayload) error {
	sqlQuery := `INSERT INTO "mst_user" (username, password) VALUES ($1, $2)`
	_, err := s.db.Exec(sqlQuery, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	rows, err := s.db.Query(`SELECT * FROM mst_user WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}

	u := new(models.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(models.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
