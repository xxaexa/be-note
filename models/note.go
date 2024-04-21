package models

type NoteStore interface {
	CreateNote(*NotePayload) error
	GetNotes() ([]*Note, error)
	GetNoteByID(id int) (*Note, error)
	UpdateNote(id int, note *NotePayload) error
	DeleteNote(id int) error
}

type Note struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      int    `json:"user_id" validate:"required"`
}

type NotePayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      int    `json:"user_id" validate:"required"`
}
