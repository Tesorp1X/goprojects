package storage

import "time"

type Storage interface {
	Save(string) error
	GetNote(int) (*Note, error)
	GetNotesList() ([]Note, error)
	DeleteNote(int) error
	AlterNote(int, Note)
}

type Note struct {
	id        int
	data      string
	timeStamp time.Time
	isClosed  bool
}
