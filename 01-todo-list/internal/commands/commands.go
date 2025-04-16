package commands

import (
	"fmt"
	"os"
)

func AddCommand(task string) {
	fmt.Fprintf(os.Stdout, "Hallo from Add")
}

func ListCommand() {
	fmt.Fprintf(os.Stdout, "Hallo from List")
}

func CompleteCommand(id int) {
	fmt.Fprintf(os.Stdout, "Hallo from Complete")
}

func DeleteCommand(id int) {
	fmt.Fprintf(os.Stdout, "Hallo from Delete")
}
