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
