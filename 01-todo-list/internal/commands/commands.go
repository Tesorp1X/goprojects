package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
)

func AddCommand(s storage.Storage, task string) {
	err := s.Save(task)
	if err != nil {
		fmt.Fprintf(s.GetSettings().ErrFile, "error: %s", err.Error())
		s.GetSettings().Logger.Fatalf("error: %s", err.Error())
		return
	}
	fmt.Fprintf(s.GetSettings().OutFile, "Task: '%s' added to your list.", task)
}

func ListCommand(s storage.Storage, allFlag bool) {
	notesList, err := s.GetNotesList()
	if err != nil {
		fmt.Fprintf(s.GetSettings().ErrFile, "error: %s", err.Error())
		s.GetSettings().Logger.Fatalf("error: %s", err.Error())
		return
	}
	minwidth := 5         // minimal cell width including any padding
	tabwidth := 8         // width of tab characters (equivalent number of spaces)
	padding := 1          // padding added to a cell before computing its width
	padchar := byte('\t') // ASCII char used for padding
	w := tabwriter.NewWriter(s.GetSettings().OutFile, minwidth, tabwidth, padding, padchar, tabwriter.AlignRight)
	fmt.Fprint(w, "|ID\tTask\tCreated\tDone\t|\n")
	for _, note := range notesList {
		if note.IsClosed() && !allFlag {
			continue
		}
		fmt.Fprintf(w, "|%d\t%s\t%s\t%t\t|\n",
			note.GetId(),
			note.GetData(),
			note.GetTimeStamp().Format(models.TimeFormat),
			note.IsClosed(),
		)

	}
	w.Flush()
}

func CompleteCommand(s storage.Storage, id int) {
	note, err := s.GetNote(id)
	if err != nil {
		s.GetSettings().Logger.Fatalf("Task with id %d not found.\n", id)
		fmt.Fprintf(s.GetSettings().ErrFile, "Task with id %d not found.\n", id)
		return
	}
	note.Close()
	err = s.AlterNote(*note)
	if err != nil {
		fmt.Fprintln(s.GetSettings().ErrFile, models.SomethingWentWrongError)
		return
	}
	fmt.Fprintln(s.GetSettings().OutFile, "Task closed.")
}

func DeleteCommand(s storage.Storage, id int) {
	err := s.DeleteNote(id)
	if err != nil {
		fmt.Fprintf(s.GetSettings().ErrFile, "Task with id %d not found.", id)
		return
	}
	fmt.Fprintln(s.GetSettings().OutFile, "Task deleted.")
}
