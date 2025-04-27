package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
		os.O_CREATE|os.O_RDWR, 0644)

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
		//tasks.exe add <taskStr>
		taskStr := strings.Join(args[2:], " ")
		commands.AddCommand(csvStorage, taskStr)
	case "list":
		//tasks.exe list
		//tasks.exe list -a
		//tasks.exe list --all
		var allFlag bool
		if len(args) > 2 {
			flag := args[2]
			switch flag {
			case "-a", "--all":
				allFlag = true
			default:
				fmt.Fprintf(appSettings.ErrFile, "invalid flag for 'list' %s", flag)
			}
		}
		commands.ListCommand(csvStorage, allFlag)
	case "complete":
		//TODO parse taskId
		if len(args) < 3 {
			fmt.Fprintln(appSettings.ErrFile, models.NoIdWasGivenError)
			fmt.Fprintln(appSettings.OutFile, models.CompleteCommandHelp)
			return
		}
		taskId, err := strconv.ParseInt(args[2], 10, 0)
		if err != nil {
			fmt.Println(appSettings.ErrFile, models.InvalidIdError)
			return
		}
		commands.CompleteCommand(csvStorage, int(taskId))
	case "delete", "remove", "-r":
		// task delete <taskId>
		// task remove <taskId>
		// task -r <taskId>
		taskId, err := strconv.ParseInt(args[2], 10, 0)
		if err != nil {
			fmt.Println(appSettings.ErrFile, models.InvalidIdError)
			return
		}
		commands.DeleteCommand(csvStorage, int(taskId))
	default:
		fmt.Fprintf(appSettings.ErrFile, "Error: unknown command %s", command)
	}

}
