package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
)

type Storage interface {
	Save(string) error
	GetNote(int) (*Note, error)
	GetNotesList() ([]Note, error)
	DeleteNote(int) error
	AlterNote(Note) error
	GetLastId() (int, error)
}

type Note struct {
	id        int
	data      string
	timeStamp time.Time
	isClosed  bool
}

// NewNoteWithId is a Note constructor. Creates a Note object with given id.
func NewNoteWithId(id int, taskStr string) *Note {
	return &Note{id: id, data: taskStr, timeStamp: time.Now(), isClosed: false}
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

func (n *Note) SetId(id int) {
	n.id = id
}

func (n *Note) SetData(taskStr string) {
	n.data = taskStr
}

func (n *Note) SetTime(newTime time.Time) {
	n.timeStamp = newTime
}

func (n *Note) SetStatus(status bool) {
	n.isClosed = status
}

var CSV_HEADERS []string = []string{"ID", "Task", "Created", "Done"}

// CsvStorage is a tool to manage '.csv' storage.
// Create only via NewCsvStorage!
// Implements Storage interface.
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
	if len(data) == 0 {
		data = append(data, CSV_HEADERS)
	}

	return &CsvStorage{storageFile: file, rawData: data, appSettings: settings}, nil
}

// Save method saves a new task Note and writes it to the CSV file.
func (s *CsvStorage) Save(taskStr string) error {
	writer := csv.NewWriter(s.storageFile)
	defer writer.Flush()

	newId, _ := s.GetLastId()
	newId++
	newIdStr := strconv.FormatInt(int64(newId), 10)
	timeStr := time.Now().String()
	statusStr := strconv.FormatBool(false)
	newLine := []string{newIdStr, taskStr, timeStr, statusStr}

	s.rawData = append(s.rawData, newLine)
	writer.Write(newLine)
	return nil
}

func (s *CsvStorage) GetNote(noteId int) (*Note, error) {
	return nil, nil
}

func (s *CsvStorage) GetNotesList() ([]Note, error) {
	return nil, nil
}

func (s *CsvStorage) DeleteNote(noteId int) error {
	return nil
}

func (s *CsvStorage) AlterNote(newNote Note) error {
	return nil
}

func (s *CsvStorage) GetLastId() (int, error) {
	return len(s.rawData) - 1, nil
}
