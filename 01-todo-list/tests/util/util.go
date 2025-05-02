package util

import (
	"strings"
	"testing"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

func AssertNotes(t *testing.T, a, b storage.Note) bool {
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
