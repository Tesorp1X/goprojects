package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/commands"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

const (
	STD_STORAGE_FILE_DIR  = "./storage/"
	STD_STORAGE_FILE_NAME = "tasks"
	CSV_FILE_SUFFIX       = ".csv"
	STD_LOG_FILE_NAME     = "./log/tasks.log"
	LOG_MSG_PREFIX        = "INFO: "
)

func main() {
	args := os.Args

	//TODO flags
	//Init logger
	logFile, errLog := os.OpenFile(STD_LOG_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errLog != nil {
		panic("Log file is missing.")
	}
	defer logFile.Close()
	logger := log.New(logFile, LOG_MSG_PREFIX, log.Ldate|log.Ltime)

	//Init  Settings
	appSettings := models.InitSettings(os.Stdout, os.Stderr, logger)

	//Init storage
	storageFile, err := os.OpenFile(STD_STORAGE_FILE_DIR+STD_STORAGE_FILE_NAME+CSV_FILE_SUFFIX,
		os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		fmt.Fprintln(appSettings.ErrFile, "Storage file is missing")
	}
	defer storageFile.Close()
	csvStorage, errStorage := storage.NewCsvStorage(storageFile, appSettings)
	if errStorage != nil {
		panic(errStorage)
	}
	//Parse command
	command := args[1]
	switch command {
	case "add":
		//TODO parse taskStr
		taskStr := "task"
		commands.AddCommand(csvStorage, taskStr)
	case "list":
		//TODO parse -a
		commands.ListCommand(csvStorage)
	case "complete":
		//TODO parse taskId
		taskId := 1
		commands.CompleteCommand(csvStorage, taskId)
	default:
		fmt.Fprintf(appSettings.ErrFile, "Error: unknown command %s", command)
	}

}
