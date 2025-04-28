package storage

import (
	"strings"
	"testing"
	"time"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
)

func assertNotes(t *testing.T, a, b Note) bool {
	t.Helper()
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

func TestNewNoteFromRawData(t *testing.T) {
	t.Run("good data", func(t *testing.T) {
		idStr := "1"
		dataStr := "some task"
		newTime := time.Now()
		statusStr := "true"
		goodData := []string{idStr, dataStr, newTime.Format(models.TimeFormat), statusStr}

		gotNote, err := NewNoteFromRawData(goodData)
		if err != nil {
			t.Errorf("didn't expect an error, but got: %v", err)
		}

		var wantedNote Note

		wantedNote.SetId(1)
		wantedNote.SetData(dataStr)
		wantedNote.SetTime(newTime)
		wantedNote.SetStatus(true)

		if !assertNotes(t, *gotNote, wantedNote) {
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

		_, err := NewNoteFromRawData(badData)

		if err == nil {
			t.Error("expected an error, but got non")
		}
		if strings.Compare(err.Error(), models.WrongNoteDataError) != 0 {
			t.Errorf("expected models.WrongNoteDataError, but got: %v", err)
		}
	})
}
