package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
)

type Storage interface {
	Save(string) error
	GetNote(int) (*Note, error)
	GetNotesList() ([]Note, error)
	DeleteNote(int) error
	AlterNote(int, Note)
	GetLastId() (int, error)
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

// CsvStorage is a tool to manage '.csv' storage.
// Create only via NewCsvStorage!
type CsvStorage struct {
	storageFile *os.File
	rawData     [][]string // matrix with all csv data in it
	appSettings *models.Settings
}

func NewCsvStorage(file *os.File, settings *models.Settings) (*CsvStorage, error) {
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintln(settings.ErrFile, err.Error())
		return nil, err
	}

	return &CsvStorage{storageFile: file, rawData: data, appSettings: settings}, nil
}
