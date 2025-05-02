package util

import (
	"strings"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

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

func AssertEqualRawData(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
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
