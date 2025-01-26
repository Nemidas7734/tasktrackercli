package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "todo", "in-progress", "done"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var tasks []Task
var nextID int

func readTasksFromFile() {
	// Open the tasks.json file (if it exists)
	file, err := os.Open("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, initialize an empty task list
			tasks = []Task{}
			nextID = 1
			return
		}
		fmt.Println("Error opening tasks file:", err)
		return
	}
	defer file.Close()

	// Use json.Decoder to read tasks from the file
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var task Task
		if err := decoder.Decode(&task); err != nil {
			fmt.Println("Error decoding task:", err)
			return
		}
		tasks = append(tasks, task)
	}

	nextID = len(tasks) + 1
}

// Write tasks to the tasks.json file using a JSON Encoder
func writeTasksToFile() {
	// Create or open the tasks.json file
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error opening tasks file for writing:", err)
		return
	}
	defer file.Close()

	// Create a json.Encoder and write tasks to the file
	encoder := json.NewEncoder(file)
	for _, task := range tasks {
		if err := encoder.Encode(task); err != nil {
			fmt.Println("Error encoding task:", err)
			return
		}
	}
}

// Add a new task
func addTask(description string) {
	task := Task{
		ID:          nextID,
		Description: description,
		Status:      "todo", // Default status
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, task)
	nextID++
	writeTasksToFile()
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

func updateTask(id int, description string) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			writeTasksToFile()
			fmt.Printf("Task ID %d updated\n", id)
			return
		}
	}
	fmt.Printf("Task ID %d not found\n", id)
}

func deleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			writeTasksToFile()
			fmt.Printf("Task ID %d deleted\n", id)
			return
		}
	}
	fmt.Printf("Task ID %d not found\n", id)
}

func markTaskStatus(id int, status string) {
	for i := range tasks {
		if tasks[i].ID == id {
			if status == "todo" || status == "in-progress" || status == "done" {
				tasks[i].Status = status
				tasks[i].UpdatedAt = time.Now()
				writeTasksToFile()
				fmt.Printf("Task ID %d marked as %s\n", id, status)
				return
			}
			fmt.Println("Invalid status:", status)
			return
		}
	}
	fmt.Printf("Task ID %d not found\n", id)
}

func listTasks(status string) {
	for _, task := range tasks {
		if status == "" || task.Status == status {
			fmt.Printf("ID: %d, Description: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s\n",
				task.ID, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
		}
	}
}

func main() {
	// Read tasks from file at the start
	readTasksFromFile()

	// Set up flags
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	markCmd := flag.NewFlagSet("mark", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Define flags for each command
	addDescription := addCmd.String("description", "", "Description of the task")
	updateID := updateCmd.Int("id", 0, "ID of the task to update")
	updateDescription := updateCmd.String("description", "", "New task description")
	deleteID := deleteCmd.Int("id", 0, "ID of the task to delete")
	markID := markCmd.Int("id", 0, "ID of the task to mark")
	markStatus := markCmd.String("status", "", "New status for the task (todo, in-progress, done)")
	listStatus := listCmd.String("status", "", "Status of the tasks to list (e.g., todo, in-progress, done)")

	// Parse the command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("expected 'add', 'update', 'delete', 'mark', or 'list' subcommands")
		os.Exit(1)
	}

	// Handle commands
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		addTask(*addDescription)
	case "update":
		updateCmd.Parse(os.Args[2:])
		updateTask(*updateID, *updateDescription)
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		deleteTask(*deleteID)
	case "mark":
		markCmd.Parse(os.Args[2:])
		markTaskStatus(*markID, *markStatus)
	case "list":
		listCmd.Parse(os.Args[2:])
		listTasks(*listStatus)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
