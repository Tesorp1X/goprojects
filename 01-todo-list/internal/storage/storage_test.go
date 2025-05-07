// test of helper-functions from strorage.go
package storage_test

import (
	"slices"
	"testing"
	"time"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
	"github.com/Tesorp1X/goprojects/01-todo-list/tests/util"
)

func TestNewNoteFromRawData(t *testing.T) {
	t.Run("good data", func(t *testing.T) {
		idStr := "1"
		dataStr := "some task"
		newTime := time.Now()
		statusStr := "true"
		goodData := []string{idStr, dataStr, newTime.Format(models.TimeFormat), statusStr}

		gotNote, err := storage.NewNoteFromRawData(goodData)
		if err != nil {
			t.Errorf("didn't expect an error, but got: %v", err)
		}

		var wantedNote storage.Note

		wantedNote.SetId(1)
		wantedNote.SetData(dataStr)
		wantedNote.SetTime(newTime)
		wantedNote.SetStatus(true)

		if !util.AssertEqualNotes(*gotNote, wantedNote) {
			//add note.String() after it's complete
			t.Error("Notes are not the same.")
		}
	})
	t.Run("bad data", func(t *testing.T) {
		idStr := "d1dsf"
		dataStr := "some task"
		newTime := time.Now()
		statusStr := "true"
		badData := []string{idStr, dataStr, newTime.Format(models.TimeFormat), statusStr}

		_, err := storage.NewNoteFromRawData(badData)

		if err == nil {
			t.Error("expected an error, but got non")
		}
		if err.Error() != models.WrongNoteDataError {
			t.Errorf("expected models.WrongNoteDataError, but got: %v", err)
		}
	})
}

func TestGenerateRawDataFromNote(t *testing.T) {
	// Set data
	note := storage.Note{}
	note.SetId(1)
	testTime := time.Now()
	note.SetData("test task")
	note.SetTime(testTime)
	note.SetStatus(true)

	got := storage.GenerateRawDataFromNote(note)
	expected := []string{"1", "test task", testTime.Format(models.TimeFormat), "true"}

	if !slices.Equal(got, expected) {
		t.Errorf("expected: %v, but got: %v", expected, got)
	}
}
