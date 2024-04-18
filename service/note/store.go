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

func (s *Store) CreateNote(note *models.CreateNotePayload) error {
	sqlQuery := `INSERT INTO "mst_note"(user_id, title, description) VALUES ($1, $2, $3) RETURNING id`
	_, err := s.db.Exec(sqlQuery, note.UserID, note.Title, note.Description)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) GetNotes() ([]*models.Note, error) {
	sqlQuery := "SELECT * FROM mst_note"
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

	sqlQuery := "SELECT * FROM mst_note WHERE id = $1"
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

func (s *Store) UpdateNote(note *models.Note) error {
	sqlQuery := `UPDATE "mst_note" SET user_id = $1, title = $2, description = $3 WHERE id = $4`
	_, err := s.db.Exec(sqlQuery, note.UserID, note.Title, note.Description, note.ID)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) DeleteNote(id int) error {
	sqlQuery := `DELETE FROM mst_note WHERE id = $1`
	_, err := s.db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}

	return err
}

func scanRowsIntoNotes(rows *sql.Rows) (*models.Note, error) {
	note := new(models.Note)

	err := rows.Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Description,
	)
	if err != nil {
		return nil, err
	}

	return note, nil

}
