package commands

import (
	"fmt"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
	"github.com/aquasecurity/table"
	"github.com/mergestat/timediff"
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
	//Table setup
	t := table.New(s.GetSettings().OutFile)
	t.SetRowLines(true)
	t.SetHeaders("ID", "Task", "Created at", "Done")
	t.SetAlignment(table.AlignRight, table.AlignRight, table.AlignRight)
	t.SetDividers(table.UnicodeRoundedDividers)

	notes, err := s.GetNotesList()
	if err != nil {
		s.GetSettings().Logger.Fatalf("Failed to retrieve a list from storage")
		fmt.Fprintln(s.GetSettings().ErrFile, models.SomethingWentWrongError)
		return
	}
	for _, note := range notes {
		//	If command was called without flag -a and
		// task is complete it will not show.
		if note.IsClosed() && !allFlag {
			continue
		}

		completed := "❌"
		if note.IsClosed() {
			completed = "✔️"
		}
		rawNote := storage.GenerateRawDataFromNote(note)
		rawNote[2] = timediff.TimeDiff(note.GetTimeStamp()) // setting [Created at] to a human readable, relative time differences
		rawNote[3] = completed                              // setting  [Done] col value to emoji ✔️ or ❌
		t.AddRow(rawNote...)
	}
	t.Render()
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
