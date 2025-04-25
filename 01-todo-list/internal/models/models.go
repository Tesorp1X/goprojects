package models

import (
	"io"
	"log"
)

type Settings struct {
	OutFile io.Writer
	ErrFile io.Writer
	Logger  *log.Logger
}

func InitSettings(outF, errF io.Writer, log *log.Logger) *Settings {
	return &Settings{OutFile: outF, ErrFile: errF, Logger: log}
}

const (
	TimeFormat = "2006-01-02 15:04:05 -0700 MST"
)

// Errors
const (
	IdOutOfRangeError  = "error: given id is out range"
	WrongNoteDataError = "error: couldn't create new note, because wrong data was given"
)
