package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("invalid number of arguments")
		return
	}

	// storage configuration
	homeDir, err := os.UserHomeDir()
	check(err)
	storageDir := homeDir + "/.task-tracker-cli/"

	_, err = os.Stat(storageDir)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(storageDir, 0755)
		check(err)
		fmt.Println("Storage created at", storageDir, "\n")

		err = os.Mkdir(storageDir + "tasks/", 0755)
		check(err)
		err = os.Mkdir(storageDir + "configs/", 0755)
		check(err)
	}
	// --------------------------

	switch flag.Arg(0) {
	case "create":
		var createTask task

		fmt.Print("Input name of the task: ")
		createTask.Name = readStdIn()

		fmt.Print("Input description of the task: ")
		createTask.Task = readStdIn()

		createTask.CreatedAt = time.Now().Format(time.RFC1123)

		createJSON, err := json.Marshal(createTask)
		check(err)

		f, err := os.Create(storageDir + "tasks/" + createTask.Name + ".json")
		check(err)
		defer f.Close()

		_, err = f.Write(createJSON)
		check(err)
	default:
		panic("unknown argument: " + flag.Arg(0))
	}
}

type task struct {
	Name      string `json:"name"`
	Task      string `json:"task"`
	CreatedAt string `json:"createdAt"`
}

func readStdIn() string {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
