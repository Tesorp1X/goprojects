package util

import (
	"os"
	"strings"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

// Returns true if both notes are identical
func AssertEqualNotes(a, b storage.Note) bool {
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

// Returns true if a[i][j] == b[i][j] for any acceptable i and j.
func AssertEqualRawData(a, b [][]string) bool {
	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] == nil || b[i] == nil {
			return false
		}
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := 0; j < len(a[i]); j++ {
			if strings.Compare(a[i][j], b[i][j]) != 0 {
				return false
			}
		}
	}

	return true
}

// Cleans a file f and closess it. Returns whatever error occured during that process.
func CleanFile(f *os.File) error {
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
