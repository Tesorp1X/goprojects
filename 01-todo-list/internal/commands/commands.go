package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

func AddCommand(s storage.Storage, task string) {
	err := s.Save(task)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		// TODO logging
		return
	}
	fmt.Fprintf(os.Stdout, "Task: '%s' added to your list.", task)
}

func ListCommand(s storage.Storage, allFlag bool) {
	notesList, err := s.GetNotesList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		// TODO logging
		return
	}
	minwidth := 5         // minimal cell width including any padding
	tabwidth := 8         // width of tab characters (equivalent number of spaces)
	padding := 1          // padding added to a cell before computing its width
	padchar := byte('\t') // ASCII char used for padding
	w := tabwriter.NewWriter(os.Stdout, minwidth, tabwidth, padding, padchar, tabwriter.AlignRight)
	fmt.Fprint(w, "|ID\tTask\tCreated\tDone\t|\n")
	for _, note := range notesList {
		if note.IsClosed() && !allFlag {
			continue
		}
		fmt.Fprintf(w, "|%d\t%s\t%s\t%t\t|\n",
			note.GetId(),
			note.GetData(),
			note.GetTimeStamp().String(),
			note.IsClosed(),
		)

	}
	w.Flush()
}

func CompleteCommand(s storage.Storage, id int) {
	note, err := s.GetNote(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		// TODO logging
		return
	}
	note.Close()
	fmt.Fprintln(os.Stdout, "Task closed.")
}

func DeleteCommand(s storage.Storage, id int) {
	err := s.DeleteNote(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Task with id %d not found.", id)
		return
	}
	fmt.Fprintln(os.Stdout, "Task deleted.")
}
