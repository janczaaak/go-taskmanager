package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID          int
	Name        string
	Description string
	Completed   bool
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Task Manager 1.0")
	fmt.Println("----------------")
	fmt.Println("#1 - Show Tasks")
	fmt.Println("#2 - Add Task")
	fmt.Println("#3 - Set task as completed")
	fmt.Println("#4 - Delete task")
	fmt.Println("#5 - Exit taskmanager")

	for {
		fmt.Println("----------------")
		fmt.Println("give me a number:")
		scanner.Scan()
		command := scanner.Text()
		command = strings.TrimSpace(command)

		switch command {
		case "1":
			tasks, err := getTasks(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("----------------")
			for _, task := range tasks {
				status := getTaskStatusString(task.Completed)
				fmt.Printf("ID: %d| Name: %s \t| Description: %s \t| Status: %s\n", task.ID, task.Name, task.Description, status)
			}
		case "2":
			fmt.Print("name the task: ")
			scanner.Scan()
			name := scanner.Text()
			fmt.Print("describe task: ")
			scanner.Scan()
			desc := scanner.Text()
			task := Task{
				Name:        name,
				Description: desc,
				Completed:   false,
			}
			addTask(db, task)
			fmt.Println("Task added, hurray")
		case "3":
			fmt.Println("id please: ")
			scanner.Scan()
			ids := scanner.Text()
			id, err := strconv.Atoi(ids)
			if err != nil {
				fmt.Println("Incorrect id, try again")
				continue
			}
			completed(db, id)
		case "4":
			fmt.Println("id please: ")
			scanner.Scan()
			idstr := scanner.Text()
			id, err := strconv.Atoi(idstr)
			if err != nil {
				fmt.Println("Incorrect id, try again")
				continue
			}
			deleteTask(db, id)
		case "5":
			fmt.Println("see you soon...")
			return
		default:

		}
	}

}
func getTaskStatusString(completed bool) string {
	if completed {
		return "Done"
	}
	return "just do it"
}

func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT * FROM Tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Completed)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tasks, nil
}

func addTask(db *sql.DB, task Task) {
	query := "INSERT INTO `Tasks` (`name`, `description`,`completed`) VALUES (?,?,?)"
	insertResult, err := db.ExecContext(context.Background(), query, task.Name, task.Description, false)
	if err != nil {
		log.Fatalf("impossible to insert task: %s", err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted id: %v\n", id)
}

func completed(db *sql.DB, id int) {
	_, err := db.Exec("UPDATE Tasks SET completed=1 WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good job")
}

func deleteTask(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM Tasks WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
}
