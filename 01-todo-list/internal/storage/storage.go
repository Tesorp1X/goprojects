package storage

import (
	"time"
)

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

func (n *Note) Close() {
	n.isClosed = false
}

func (n Note) GetId() int {
	return n.id
}

func (n Note) GetData() string {
	return n.data
}

func (n Note) GetTimeStamp() time.Time {
	return n.timeStamp
}

func (n Note) IsClosed() bool {
	return n.isClosed
}
