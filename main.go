package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Id          string
	Name        string
	Description string
	Enabled     string
}

func (task Task) String() string {
	return fmt.Sprintf("enabled:%s\nid: %s\nname: %s\ndescription:%s\n", task.Enabled, task.Id, task.Name, task.Description)
}

func (t Task) Print() {
	fmt.Println(t.String())
}

var verbose = flag.Bool("v", false, "enable verbose log")
var vverbose = flag.Bool("vv", false, "enable super verbose log")

func main() {
	id := flag.String("id", "", "a Id to retrieve some task")
	name := flag.String("name", "", "a name to a task")
	description := flag.String("description", "", "a description of a task")
	action := flag.String("action", "list", "a action to run a cli")
	version := flag.Bool("version", false, "the version of cli")
	all := flag.Bool("all", false, "a flag to list all enabled")

	flag.Parse()

	vvlogf("action=%s,id=%s,name=%s,description=%s,verbose=%t,vverbose=%t\n", *action, *id, *name, *description, *verbose, *vverbose)

	if *version == true {
		log.Println(getVersion())
		return
	}

	db, err := openDatabase()
	if err != nil {
		log.Fatalln(err, "failed to open database")
	}
	defer db.Close()

	err = runMigration(db)
	if err != nil {
		log.Fatalln(err)
	}

	switch *action {
	case "create":
		vlogf("lets create")

		if *name != "" {
			createTask(db, Task{Name: *name, Description: *description})

		} else {
			log.Fatalln("a task needs a name")
		}
	case "list":
		vlogf("listing tasks")
		lists, err := listTasks(db, *all)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(lists)
	case "get":
		if *id != "" {
			getTask(db, *id)
		} else {
			log.Println("a id should be provided")
		}
	case "complete":
		if *id != "" {
			completeTask(db, *id)
		} else {
			log.Println("to complete task need a id")
		}
	}

}

func openDatabase() (*sql.DB, error) {
	databaseName := "tasks.db"
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	databasePath := fmt.Sprintf("%s/.local/gotasks", homePath)

	err = os.MkdirAll(databasePath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	databaseUrl := fmt.Sprintf("%s/%s", databasePath, databaseName)

	db, err := sql.Open("sqlite3", databaseUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigration(db *sql.DB) error {
	migration := `CREATE TABLE IF NOT EXISTS tasks 
	(id INTEGER PRIMARY KEY AUTOINCREMENT,
	name text,
	enabled boolean DEFAULT true,
	description text)`

	_, err := db.Exec(migration)
	if err != nil {
		return err
	}

	vlogf("migration runned")

	return nil
}

func listTasks(db *sql.DB, all bool) ([]Task, error) {
	query := "SELECT id, enabled, name, description FROM tasks WHERE enabled = true"
	if all {
		query = "SELECT id, enabled, name, description FROM tasks"
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("failed to query")
	}

	var result []Task

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Enabled, &task.Name, &task.Description)
		if err != nil {
			log.Fatalln(err)
		}

		result = append(result, task)
	}

	return result, nil
}

func createTask(db *sql.DB, inputTask Task) (*Task, error) {
	var task Task
	err := db.QueryRow(`INSERT INTO tasks (name, description)
	VALUES (?, ?)
	RETURNING id, name, description, enabled;`, inputTask.Name, inputTask.Description).Scan(&task.Id, &task.Name, &task.Description, &task.Enabled)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("task inserted")
	fmt.Printf("id: %s\nenabled: %s\nname: %s\ndescription:%s\n", task.Id, task.Enabled, task.Name, task.Description)

	return &task, nil
}

func getTask(db *sql.DB, id string) (*Task, error) {
	var task Task

	err := db.QueryRow(`SELECT id, enabled, name, description FROM tasks WHERE id = ?`, id).Scan(&task.Id, &task.Enabled, &task.Name, &task.Description)

	if err != nil {
		return nil, errors.New("failed to scan result")
	}

	task.Print()

	return &task, nil
}

func completeTask(db *sql.DB, id string) (*Task, error) {
	var task Task
	err := db.QueryRow(`UPDATE tasks SET enabled = false WHERE id = ? RETURNING id, enabled, name, description`, id).Scan(&task.Id, &task.Enabled, &task.Name, &task.Description)
	if err != nil {
		return nil, errors.New("failed to update")
	}
	log.Println("task completed", task)

	return &task, nil
}

func vlogf(format string, v ...any) {
	if *verbose {
		log.Printf(format, v...)
	}
}

func vvlogf(format string, v ...any) {
	if *vverbose {
		log.Printf(format, v...)
	}
}

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	version := info.Main.Version
	return version
}
