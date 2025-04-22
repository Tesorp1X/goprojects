package models

import (
	"log"
	"os"
)

type Settings struct {
	OutFile *os.File
	ErrFile *os.File
	Logger  *log.Logger
}

func InitSettings(outF, errF *os.File, log *log.Logger) *Settings {
	return &Settings{OutFile: outF, ErrFile: errF, Logger: log}
}

const (
	TimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

// Errors
const (
	IdOutOfRangeError  = "error: given id is out range"
	WrongNoteDataError = "error: couldn't create new note, because wrong data was given"
)
