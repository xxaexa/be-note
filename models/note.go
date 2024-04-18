package models

type NoteStore interface {
	CreateNote(*CreateNotePayload) error
	GetNotes() ([]*Note, error)
	GetNoteByID(id int) (*Note, error)
	UpdateNote(*Note) error
	DeleteNote(id int) error
}

type Note struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateNotePayload struct {
	UserID      int    `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
