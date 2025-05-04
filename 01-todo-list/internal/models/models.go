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
	IdNotFoundError         = "error: given id not found"
	WrongNoteDataError      = "error: couldn't create new note, because wrong data was given"
	NoIdWasGivenError       = "error: no id was given"
	InvalidIdError          = "error: invalid id was given"
	SomethingWentWrongError = "error: something went wrong"
)

// Helps
const (
	CompleteCommandHelp = "To mark task as complete use: task complete <task_id>"
)
