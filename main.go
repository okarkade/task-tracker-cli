package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
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
	arg := parseArg()
	configureStorage()

	switch arg {
	case "create":
		var task task

		task.Name = ask("Input name of the task")
		task.Task = ask("Input description of the task")
		task.CreatedAt = time.Now().Format(time.RFC1123)

		JSON, err := json.Marshal(task)
		check(err)

		f, err := os.Create(storageDir + "tasks/" + task.Name + ".json")
		check(err)
		defer f.Close()

		_, err = f.Write(JSON)
		check(err)

		fmt.Println("\nTask successfully created!")
	default:
		log.Fatal("unknown argument: " + arg)
	}
}

type task struct {
	Name      string `json:"name"`
	Task      string `json:"task"`
	CreatedAt string `json:"createdAt"`
}

func ask(prompt string) string {
	fmt.Print(prompt + ": ")
	return readStdIn()
}

func readStdIn() string {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func configureStorage() {
	configureFolder(storageDir)
	configureFolder(tasksDir)
	configureFolder(configsDir)
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
	}
}

func parseArg() string {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Invalid number of arguments")
	}
	return flag.Arg(0)
}