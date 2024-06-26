package note

import (
	"database/sql"
	"go-note/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateNote(note *models.NotePayload) error {
	sqlQuery := `INSERT INTO notes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id`
	_, err := s.db.Exec(sqlQuery, note.Title, note.Description, note.UserID)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) GetNotes() ([]*models.Note, error) {
	sqlQuery := `SELECT * FROM notes`
	rows, err := s.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}

	notes := make([]*models.Note, 0)
	for rows.Next() {
		note, err := scanRowsIntoNotes(rows)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil

}

func (s *Store) GetNoteByID(id int) (*models.Note, error) {

	exists, err := checkID(id, s.db)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	sqlQuery := `SELECT * FROM notes WHERE id = $1`
	rows, err := s.db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}

	note := new(models.Note)
	for rows.Next() {
		note, err = scanRowsIntoNotes(rows)
		if err != nil {
			return nil, err
		}
	}

	return note, nil
}

func (s *Store) UpdateNote(id int, note *models.NotePayload) error {

	exists, err := checkID(id, s.db)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	sqlQuery := `UPDATE notes SET title = $1, description = $2, user_id = $3 WHERE id = $4`
	_, err = s.db.Exec(sqlQuery, note.Title, note.Description, note.UserID, id)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) DeleteNote(id int) error {

	exists, err := checkID(id, s.db)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	sqlQuery := `DELETE FROM notes WHERE id = $1`
	_, err = s.db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	return err
}

func checkID(id int, db *sql.DB) (bool, error) {
	exists := false

	sqlQuery := `SELECT EXISTS (SELECT * FROM notes WHERE id = $1)`
	err := db.QueryRow(sqlQuery, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func scanRowsIntoNotes(rows *sql.Rows) (*models.Note, error) {
	note := new(models.Note)

	err := rows.Scan(
		&note.ID,
		&note.Title,
		&note.Description,
		&note.UserID,
	)
	if err != nil {
		return nil, err
	}

	return note, nil

}
