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

// const doesn't allow to use functions for initialization so use var
var (
	homeDir, _ = os.UserHomeDir()

	// slashes on the end for the easier concatanation if needed
	storageDir = homeDir + "/.task-tracker-cli/"
	tasksDir = storageDir + "tasks/"
	configsDir = storageDir + "configs/"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("invalid number of arguments")
		return
	}

	configureStorage()

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

func configureStorage() {
	configureFolder(storageDir)
	configureFolder(tasksDir)
	configureFolder(configsDir)
	fmt.Println()
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func configureFolder(path string) {
	if !isExists(path) {
		err := os.Mkdir(path, 0755)
		check(err)
		fmt.Println(path, "created")
	}
}