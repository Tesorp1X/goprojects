// Tests for exportable functions and methods
package tests

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"testing"

	//"github.com/spf13/afero"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

func TestCsvSave(t *testing.T) {
	t.Run("ordinary note", func(t *testing.T) {
		outputBuff := &bytes.Buffer{}
		errBuff := &bytes.Buffer{}
		logBuff := &bytes.Buffer{}
		logger := log.New(logBuff, "TEST: ", log.Ldate|log.Ltime)
		settings := models.InitSettings(outputBuff, errBuff, logger)

		testCsvFile, err := os.CreateTemp(".", "test.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer testCsvFile.Close()
		defer os.Remove(testCsvFile.Name())

		testCsvData := `ID,Description,CreatedAt,IsComplete
1,My new task,2024-07-27T16:45:19-05:00,true
2,Finish this video,2024-07-27T16:45:26-05:00,true
3,Find a video editor,2024-07-27T16:45:31-05:00,false`

		if _, err := testCsvFile.Write([]byte(testCsvData)); err != nil {
			t.Fatal(err)
		}
		csvStorage, _ := storage.NewCsvStorage(testCsvFile, settings)
		csvStorage.Save("test this thing")
		gotBytes := make([]byte, 512)
		_, errB := testCsvFile.Read(gotBytes)
		if errB != nil && errB != io.EOF {
			t.Fatal(errB)
		}

		expected := testCsvData + "\n4,test this thing," + time.Now().Format(models.TimeFormat) + ",false"
		if c := strings.Compare(string(gotBytes), expected); c != 0 {
			t.Errorf("expected %s got %s", expected, string(gotBytes))
		}

	})
}

func TestGetNote(t *testing.T) {
	compareNotes := func(a, b *storage.Note) bool {
		idComp := a.GetId() != b.GetId()
		dataComp := strings.Compare(a.GetData(), b.GetData())
		timeComp := strings.Compare(a.GetTimeStamp().Format(models.TimeFormat), b.GetTimeStamp().Format(models.TimeFormat))
		statusComp := a.IsClosed() == b.IsClosed()
		if idComp || dataComp != 0 ||
			timeComp != 0 ||
			!statusComp {
			return false
		}
		return true
	}
	outputBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	logBuff := &bytes.Buffer{}
	logger := log.New(logBuff, "TEST: ", log.Ldate|log.Ltime)
	settings := models.InitSettings(outputBuff, errBuff, logger)
	testCsvFile, err := os.Open("test_3lines.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer testCsvFile.Close()
	csvStorage, _ := storage.NewCsvStorage(testCsvFile, settings)

	t.Run("get note from file", func(t *testing.T) {
		note, err := csvStorage.GetNote(7)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		rawNoteData := strings.Split("7,Find a video editor,2025-04-26 03:05:39 +0700 +07,false", ",")
		wantedNote, err := storage.NewNoteFromRawData(rawNoteData)
		if err != nil {
			t.Fatal(err)
		}
		if !compareNotes(note, wantedNote) {
			t.Error("notes doesnt match")
		}
	})
	t.Run("wrong id error", func(t *testing.T) {
		_, err := csvStorage.GetNote(10000)
		wantedErr := models.IdNotFoundError

		if err == nil {
			t.Errorf("expected an error (%s), but got nothing", wantedErr)
		}

		if strings.Compare(err.Error(), wantedErr) != 0 {
			t.Errorf("expected: '%s' got: '%s'", wantedErr, err)
		}
	})
}
