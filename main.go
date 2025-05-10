package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
)

// const doesn't allow to use functions for initialization so use var
var (
	homeDir, _ = os.UserHomeDir()

	// slashes on the end for the easier concatanation if needed
	storageDir = homeDir + "/.task-tracker-cli/"
	tasksDir   = storageDir + "tasks/"
	idPoolPath = tasksDir + "idPool.json"
)

func main() {
	flag.Parse()

	configureFolder(storageDir)
	configureFolder(tasksDir)
	
	// idPool configuration
	switch isExist(idPoolPath) {
	case false:
		_, err := os.Create(idPoolPath)
		check(err)
	case true:
		readAndUnmarshal(idPoolPath, &idPool)
	}

	switch flag.Arg(0) {
	case "create", "add":
		if len(flag.Args()) != 2 {
			log.Fatal("invalid number of arguments")
		}
		var task task
		
		task.ID = generateID()
		idPool[task.ID] = nil

		task.Task = flag.Arg(1)
		task.CreatedAt = time.Now().Format(time.RFC1123)
		task.Status = statusActive

		_, err := os.Create(tasksDir + strconv.Itoa(task.ID) + ".json")
		check(err)
		marshalAndWrite(task, tasksDir + strconv.Itoa(task.ID) + ".json")

		fmt.Println("Task", task.ID, "successfully created!")

	case "list", "ls":
		if len(flag.Args()) > 2 {
			log.Fatal("invalid number of arguments")
		}

		var filter taskStatus

		switch flag.Arg(1) {
		case "active", "":
			filter = statusActive

		case "done":
			filter = statusDone

		case "inactive":
			filter = statusInactive

		default:
			log.Fatal("unknown task status:", flag.Arg(1))
		}

		for id := range idPool {
			var task task
			readAndUnmarshal(tasksDir + strconv.Itoa(id) + ".json", &task)

			if task.Status != filter {
				continue
			}

			fmt.Println(task)
			fmt.Println()
		}
		
	case "mark":
		if len(flag.Args()) != 3 {
			log.Fatal("invalid number of arguments")
		}

		id, err := strconv.Atoi(flag.Arg(1))
		if errors.Is(err, strconv.ErrRange) {
			log.Fatal(flag.Arg(1), "is not a task id")
		}

		if !isExist(tasksDir + strconv.Itoa(id) + ".json") {
			log.Fatal("there is no task with id", id)
		}

		var newStatus taskStatus

		switch flag.Arg(2) {
		case "active":
			newStatus = statusActive

		case "done":
			newStatus = statusDone

		case "inactive":
			newStatus = statusInactive

		default:
			log.Fatal("unknown task status:", flag.Arg(2))
		}

		var task task
		readAndUnmarshal(tasksDir + strconv.Itoa(id) + ".json", &task)

		if task.Status == newStatus {
			log.Fatal("task is already have status ", newStatus)
		} else {
			task.Status = newStatus
			marshalAndWrite(&task, tasksDir + strconv.Itoa(id) + ".json")
		}

	default:
		log.Fatal("unknown argument:", flag.Arg(0))
	}

	marshalAndWrite(idPool, idPoolPath)
}

type task struct {
	ID        int        `json:"id"`
	Task      string     `json:"task"`
	CreatedAt string     `json:"createdAt"`
	Status    taskStatus `json:"status"`
}

func (t task) String() string {
	return strconv.Itoa(t.ID) + ": " + t.Task +
	"\nstatus: " + t.Status.String() +
	"\ncreated at: " + t.CreatedAt
}

type taskStatus int

const (
	statusActive taskStatus = iota
	statusDone
	statusInactive
)

var statusName = map[taskStatus]string{
	statusActive: "active",
	statusDone:   "done",
	statusInactive: "inactive",
}

func (ts taskStatus) String() string {
	return statusName[ts]
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func configureFolder(path string) {
	if !isExist(path) {
		err := os.Mkdir(path, 0755)
		check(err)
	}
}

// pool for used ids
// use map[int]any for convinient search in it
// also search in hash tables is O(1)
var idPool = make(map[int]any)

// generateID is simple search algorithm because it is not guarantied that idPool is sorted
func generateID() int {
	for i := 1; ; i++ {
		_, present := idPool[i]
		if !present {
			return i
		}
	}
}

func marshalAndWrite(v any, path string) {
	JSON, err := json.Marshal(v)
	check(err)

	os.WriteFile(path, JSON, 0755)
}

func readAndUnmarshal(path string, v any) {
	JSON, err := os.ReadFile(path)
	check(err)

	json.Unmarshal(JSON, v)
}