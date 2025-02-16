package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
    ID        int
    Title     string
    Completed bool
}
type TodoList struct{
    db *sql.DB
}
func NewTodoList(dbPath string) (*TodoList, error){
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %v", err)
    }
    createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "title" TEXT,
        "completed" BOOLEAN
    );`
    if _, err := db.Exec(createTableSQL); err != nil {
        return nil, fmt.Errorf("error creating tasks table: %v", err)
    }

    return &TodoList{db:db}, nil
}

func (t *TodoList) AddTask(title string) error{
    insert_query := `INSERT INTO tasks (title, completed) VALUES (?, ?)`
    _, err := t.db.Exec(insert_query, title, false)
    if err != nil {
        return fmt.Errorf("error adding task: %v", err)
    }
    fmt.Println("Task added successfully")
    return nil
}
func (t *TodoList) ListTasks() error{
    rows, err := t.db.Query("SELECT id, title, completed FROM tasks")
    if err != nil {
        return fmt.Errorf("error listing tasks: %v", err)
    }
    defer rows.Close()
    fmt.Println("\nYour Tasks:")
    fmt.Println("------------------")
    for rows.Next() {
        var task Task
        err = rows.Scan(&task.ID, &task.Title, &task.Completed)
        if err != nil {
            return fmt.Errorf("error scanning task: %v", err)
        }
         status := "[ ]"
        if task.Completed {
            status = "[✓]"
        }
        fmt.Printf("%d. %s %s\n", task.ID, status, task.Title)
    }
     fmt.Println("------------------")
    return nil
}

func (t *TodoList) DeleteTask(id int) error{
    delete_query := `DELETE FROM tasks WHERE id = ?`
    result, err := t.db.Exec(delete_query, id)
    if err != nil {
        return fmt.Errorf("error deleting task: %v", err)
    }
    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return fmt.Errorf("task with id %d not found", id)
    }

    fmt.Printf("Task %d deleted successfully", id)
    return nil
}

func (t *TodoList) ToggleTask(id int) error{
    var completed bool
    err := t.db.QueryRow("SELECT completed FROM tasks WHERE id = ?", id).Scan(&completed)
    if err != nil {
        return fmt.Errorf("error getting task: %v", err)
    }

    _,err=t.db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", !completed, id)
    if err != nil {
        return fmt.Errorf("error updating task: %v", err)
    }
    status:= "incomplete"
    if !completed{
        status="completed"
    }
    fmt.Printf("Task %d marked as %s\n", id, status)
    return nil
}
func (t *TodoList) Close() error{
    return t.db.Close()
}

func main() {
    t, err := NewTodoList("./tasks.db")
    if err != nil {
        log.Fatal(err)
    }
    defer t.Close()

	//CLI commands
    addCmd := flag.NewFlagSet("add", flag.ExitOnError)
    addTitle := addCmd.String("title", "", "Title of the task to add")

    deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
    deleteID := deleteCmd.Int("id", 0, "ID of the task to delete")

    toggleCmd := flag.NewFlagSet("toggle", flag.ExitOnError)
    toggleID := toggleCmd.Int("id", 0, "ID of the task to toggle")
    if len(os.Args) < 2 {
        fmt.Println("Expected 'add', 'list', 'delete', or 'toggle' subcommands")
        os.Exit(1)
    }
    switch os.Args[1] {
    case "add":
        addCmd.Parse(os.Args[2:])
        if *addTitle == "" {
            fmt.Println("Please provide a title for the task using the -title flag")
            os.Exit(1)
        }
        if err:=t.AddTask(*addTitle); err != nil {
            log.Fatalf("error adding task: %v", err)
        }

    case "list":
        if err:=t.ListTasks(); err != nil {
            log.Fatalf("error listing tasks: %v", err)
        }
    case "delete":
        deleteCmd.Parse(os.Args[2:])
        if *deleteID == 0 {
            fmt.Println("Please provide an id for the task using the -id flag")
            os.Exit(1)
        }
        if err:=t.DeleteTask(*deleteID); err != nil {
            log.Fatalf("error deleting task: %v", err)
        }
    case "toggle":
        toggleCmd.Parse(os.Args[2:])
        if *toggleID == 0 {
            fmt.Println("Please provide an id for the task using the -id flag")
            os.Exit(1)
        }
        if err:=t.ToggleTask(*toggleID); err != nil {
            log.Fatalf("error toggling task: %v", err)
        }   

    default:
        fmt.Println("Expected 'add', 'list', 'delete', or 'toggle' subcommands")
        os.Exit(1)
    }


}