// Tests for exportable functions and methods
package tests

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"testing"

	"github.com/Tesorp1X/goprojects/01-todo-list/internal/models"
	"github.com/Tesorp1X/goprojects/01-todo-list/internal/storage"
	"github.com/Tesorp1X/goprojects/01-todo-list/tests/util"
)

func TestCsvSave(t *testing.T) {
	rawDataToStr := func(data [][]string) string {
		t.Helper()
		var str string
		for _, row := range data {
			str += fmt.Sprintf("%s\n", strings.Join(row, ","))
		}
		return str
	}

	var testCsvData [][]string
	testCsvData = append(testCsvData, []string{"ID", "Task", "Created", "Done"})
	testCsvData = append(testCsvData, []string{"1", "My new task", "2025-04-26 03:05:39 +0700 +07", "true"})
	testCsvData = append(testCsvData, []string{"2", "Finish this video", "2025-04-26 03:05:39 +0700 +07", "true"})
	testCsvData = append(testCsvData, []string{"3", "Find a video editor", "2025-04-26 03:05:39 +0700 +07", "false"})

	prepareFile := func(fileName string) {
		t.Helper()
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		if err := writer.WriteAll(testCsvData); err != nil {
			t.Fatal(err)
		}
		writer.Flush()
	}

	saveNewNote := func(fileName string, taskStr string, testCsvData [][]string, settings *models.Settings) (expectedData [][]string) {
		testCsvFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		// Closing a file so data gets saved
		defer testCsvFile.Close()

		csvStorage, _ := storage.NewCsvStorage(testCsvFile, settings)
		csvStorage.Save(taskStr)
		noteId, _ := csvStorage.GetLastId()
		noteId++
		moment := time.Now().Format(models.TimeFormat)

		expectedData = append(testCsvData, []string{strconv.FormatInt(int64(noteId), 10), taskStr, moment, "false"})
		return
	}

	outputBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	logBuff := &bytes.Buffer{}
	logger := log.New(logBuff, "TEST: ", log.Ldate|log.Ltime)
	settings := models.InitSettings(outputBuff, errBuff, logger)

	t.Run("ordinary note", func(t *testing.T) {

		//Opens a file and puts [testCsvData] in to it
		prepareFile("test_empty.csv")

		//Saving new note
		expectedData := saveNewNote("test_empty.csv", "test that thing", testCsvData, settings)
		// Assertion
		testCsvFile, err := os.OpenFile("test_empty.csv", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer util.CleanFile(testCsvFile)

		reader := csv.NewReader(testCsvFile)
		gotData, err := reader.ReadAll()
		if err != nil {
			t.Fatal(err)
		}

		if !util.AssertEqualRawData(expectedData, gotData) {
			t.Errorf("\nEXPECTED:\n%s \nGOT:\n%s", rawDataToStr(expectedData), rawDataToStr(gotData))
		}
	})
	t.Run("note has ',' in noteStr", func(t *testing.T) {
		//Opens a file and puts [testCsvData] in to it
		prepareFile("test_empty.csv")

		//Saving new note
		expectedData := saveNewNote("test_empty.csv", "test, this, thing", testCsvData, settings)

		// Assertion
		testCsvFile, err := os.OpenFile("test_empty.csv", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer util.CleanFile(testCsvFile)

		reader := csv.NewReader(testCsvFile)
		gotData, err := reader.ReadAll()
		if err != nil {
			t.Fatal(err)
		}

		if !util.AssertEqualRawData(expectedData, gotData) {
			t.Errorf("\nEXPECTED:\n%s \nGOT:\n%s", rawDataToStr(expectedData), rawDataToStr(gotData))
		}
	})
	t.Run("first note", func(t *testing.T) {
		var data [][]string
		data = append(data, storage.CSV_HEADERS)

		//Saving new note
		expectedData := saveNewNote("test_empty.csv", "test, this, thing", data, settings)

		// Assertion
		testCsvFile, err := os.OpenFile("test_empty.csv", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer util.CleanFile(testCsvFile)

		reader := csv.NewReader(testCsvFile)
		gotData, err := reader.ReadAll()
		if err != nil {
			t.Fatal(err)
		}

		if !util.AssertEqualRawData(expectedData, gotData) {
			t.Errorf("\nEXPECTED:\n%s \nGOT:\n%s", rawDataToStr(expectedData), rawDataToStr(gotData))
		}
	})
}

func TestGetNote(t *testing.T) {
	outputBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}
	logBuff := &bytes.Buffer{}
	logger := log.New(logBuff, "TEST: ", log.Ldate|log.Ltime)
	settings := models.InitSettings(outputBuff, errBuff, logger)
	testCsvFile, err := os.Open("test_3lines.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer testCsvFile.Close()
	csvStorage, _ := storage.NewCsvStorage(testCsvFile, settings)

	t.Run("get note from file", func(t *testing.T) {
		note, err := csvStorage.GetNote(7)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		rawNoteData := strings.Split("7,Find a video editor,2025-04-26 03:05:39 +0700 +07,false", ",")
		wantedNote, err := storage.NewNoteFromRawData(rawNoteData)
		if err != nil {
			t.Fatal(err)
		}
		if !util.AssertEqualNotes(*note, *wantedNote) {
			t.Error("notes doesnt match")
		}
	})
	t.Run("wrong id error", func(t *testing.T) {
		_, err := csvStorage.GetNote(10000)
		wantedErr := models.IdNotFoundError

		if err == nil {
			t.Errorf("expected an error (%s), but got nothing", wantedErr)
		}

		if strings.Compare(err.Error(), wantedErr) != 0 {
			t.Errorf("expected: '%s' got: '%s'", wantedErr, err)
		}
	})
}
