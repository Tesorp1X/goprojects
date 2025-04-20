package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/commands"
)

type Settings struct {
	OutFile *os.File
	ErrFile *os.File
	Logger  *log.Logger
}

func InitSettings(outF, errF *os.File, log *log.Logger) *Settings {
	return &Settings{OutFile: outF, ErrFile: errF, Logger: log}
}

var APP_SETTINGS *Settings

const (
	STD_STORAGE_FILE_DIR  = "/storage/"
	STD_STORAGE_FILE_NAME = "tasks"
	CSV_FILE_SUFFIX       = ".csv"
	STD_LOG_FILE_NAME     = "/log/tasks.log"
	LOG_MSG_PREFIX        = "INFO: "
)

func main() {
	args := os.Args

	//TODO flags
	//Init logger
	logFile, errLog := os.Open(STD_LOG_FILE_NAME)
	if errLog != nil {
		panic("Log file is missing.")
	}
	defer logFile.Close()
	logger := log.New(logFile, LOG_MSG_PREFIX, log.Ldate|log.Ltime)

	//Init  Settings
	APP_SETTINGS = InitSettings(os.Stdout, os.Stderr, logger)

	//Init storage
	storageFile, err := os.Open(STD_STORAGE_FILE_DIR + STD_STORAGE_FILE_NAME + CSV_FILE_SUFFIX)
	if err != nil {
		fmt.Fprintln(APP_SETTINGS.ErrFile, "Storage file is missing")
	}
	defer storageFile.Close()

	//Parse command
	command := args[1]
	switch command {
	case "add":
		taskStr := "task"
		commands.AddCommand(taskStr)
	case "list":
		commands.ListCommand()
	case "complete":
		taskId := 1
		commands.CompleteCommand(taskId)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %s", command)
	}

}
