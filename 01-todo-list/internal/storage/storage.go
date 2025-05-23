package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	GetSettings() *models.Settings
}

type Note struct {
	id        int
	data      string
	timeStamp time.Time
	isClosed  bool
}

// CreateNewNoteWithId is a Note constructor. Creates a Note object with given id.
func CreateNewNoteWithId(id int, taskStr string, timeStamp time.Time, status bool) *Note {
	return &Note{id: id, data: taskStr, timeStamp: timeStamp, isClosed: status}
}

// Creates a [Note] from raw data.
// [rawData] must be in format: ["number", "string", "models.TimeFormat", "bool"],
// otherwise returns a [WrongNoteDataError].
func NewNoteFromRawData(rawData []string) (*Note, error) {
	id, errId := strconv.ParseInt(rawData[0], 10, 0)
	data := rawData[1]
	timeStamp, errTime := time.Parse(models.TimeFormat, rawData[2])
	status, errStatus := strconv.ParseBool(rawData[3])
	if errId != nil || errTime != nil || errStatus != nil {
		return nil, errors.New(models.WrongNoteDataError)
	}
	return CreateNewNoteWithId(int(id), data, timeStamp, status), nil
}

// Generates raw data -- slice of string in order: [id, data, timeStamp, isClosed]
func GenerateRawDataFromNote(n Note) []string {
	var rawData []string
	id := strconv.Itoa(n.id)
	timeStamp := n.timeStamp.Format(models.TimeFormat)
	status := strconv.FormatBool(n.isClosed)
	rawData = append(rawData, id, n.data, timeStamp, status)
	return rawData
}

// Converts note to string like that: {ID: 1; Data: test task; Created at: 2006-01-02 15:04:05 -0700 MST; Done: true}
func (n Note) String() string {
	var res string
	timeStamp := n.timeStamp.Format(models.TimeFormat)
	res = fmt.Sprintf("{ID: %d; Data: %s; Created at: %s; Done: %t}", n.id, n.data, timeStamp, n.isClosed)
	return res
}

func (n *Note) Close() {
	n.isClosed = true
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
	stagedData  []*Note // data to save. Must call .flush() to save it
}

// Writes [stagedData] to [storageFile].
func (s *CsvStorage) flush() error {
	if len(s.stagedData) == 0 {
		return nil
	}
	w := csv.NewWriter(s.storageFile)
	defer w.Flush()
	for _, note := range s.stagedData {
		id := strconv.Itoa(note.id)
		timeStamp := note.timeStamp.Format(models.TimeFormat)
		status := strconv.FormatBool(note.isClosed)

		err := w.Write([]string{id, note.data, timeStamp, status})
		if err != nil {
			s.appSettings.Logger.Fatalf("Failed to flush csv file: %v", err)
			return err
		}
		s.appSettings.Logger.Printf("Note saved: {ID:%s,Data:%s,TimeStamp:%s,Status:%s}",
			id, note.data, timeStamp, status)
	}
	s.appSettings.Logger.Println("Flush successful")
	return nil
}

func (s *CsvStorage) GetSettings() *models.Settings {
	return s.appSettings
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
		w := csv.NewWriter(file)
		if err := w.Write(CSV_HEADERS); err != nil {
			return nil, err
		}
		w.Flush()

	}

	return &CsvStorage{storageFile: file, rawData: data, appSettings: settings}, nil
}

// Save method saves a new task Note and writes it to the CSV file.
// Returns whatever error occured during [WriteAll] method call.
func (s *CsvStorage) Save(taskStr string) error {
	newId, _ := s.GetLastId()
	newId++

	newNote := CreateNewNoteWithId(newId, taskStr, time.Now(), false)

	s.stagedData = append(s.stagedData, newNote)

	return s.flush()
}

func lookForId(data [][]string, id int) (found bool, noteRaw []string) {
	for _, row := range data {
		if rowId, _ := strconv.ParseInt(row[0], 10, 0); rowId == int64(id) {
			found = true
			noteRaw = row
			return
		}
	}
	return
}

func (s *CsvStorage) GetNote(noteId int) (*Note, error) {
	found, noteRaw := lookForId(s.rawData, noteId)
	if !found {
		return nil, errors.New(models.IdNotFoundError)
	}

	return NewNoteFromRawData(noteRaw)
}

// Returns a slice of Note and any error, that NewNoteFromRawData produced.
func (s *CsvStorage) GetNotesList() ([]Note, error) {
	var notes []Note

	for _, line := range s.rawData[1:] {
		note, err := NewNoteFromRawData(line)
		if err != nil {
			return notes, err
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

// Helper muthod, that clears all data from [storageFile] and pushes
// [CSV_HEADERS] back in the file. Also clears [rawData].
func (s *CsvStorage) clearAll() {
	if err := s.storageFile.Truncate(0); err != nil {
		s.appSettings.Logger.Fatalf("Failed to truncate: %v", err)
	}
	s.storageFile.Seek(0, 0)
	s.rawData = [][]string{CSV_HEADERS}
	if _, err := s.storageFile.WriteString(strings.Join(CSV_HEADERS, ",") + "\n"); err != nil {
		s.appSettings.Logger.Fatalf("Failed to write to csv file: %v", err)
	}

}

// Deletes a note with given [noteId]. Returns a [IdNotFoundError] if
// note was note found in the file. Method will delete all data from the file,
//
//	and rewrite altered data (without note with given id).
func (s *CsvStorage) DeleteNote(noteId int) error {
	if found, _ := lookForId(s.rawData, noteId); !found {
		s.appSettings.Logger.Fatalf("Failed to find ID: %d", noteId)
		return errors.New(models.IdNotFoundError)
	}
	for _, line := range s.rawData[1:] {
		note, _ := NewNoteFromRawData(line)
		if note.id != noteId {
			s.stagedData = append(s.stagedData, note)
		}
	}
	s.clearAll()
	if err := s.flush(); err != nil {
		return err
	}
	s.appSettings.Logger.Printf("Deleted note with id: %d", noteId)
	return nil
}

// AlterNote recieves a note and updates csv-file with new data.
// Returns whatever error, that happens during clearAll or flush method calls.
func (s *CsvStorage) AlterNote(newNote Note) error {
	for _, row := range s.rawData[1:] {
		currentId, _ := strconv.ParseInt(row[0], 10, 0)
		if currentId == int64(newNote.id) {
			s.stagedData = append(s.stagedData, &newNote)
		} else {
			n, _ := NewNoteFromRawData(row)
			s.stagedData = append(s.stagedData, n)
		}
	}
	s.clearAll()
	if err := s.flush(); err != nil {

	}
	return nil
}

// Returns id of last Note, if Csv-file is empty, then method returns 0.
// Returns [ErrSyntax] error, if last id is not int.
func (s *CsvStorage) GetLastId() (int, error) {
	if len(s.rawData) < 2 {
		return 0, nil
	}

	lastId, err := strconv.ParseInt(s.rawData[len(s.rawData)-1][0], 10, 0)

	return int(lastId), err
}
