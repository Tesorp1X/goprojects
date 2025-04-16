package main

import (
	"fmt"
	"os"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/commands"
)

func main() {
	args := os.Args
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
