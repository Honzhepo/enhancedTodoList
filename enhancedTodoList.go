package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	webview "github.com/webview/webview_go"
)

// Task structure for a todo
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	Deadline  time.Time `json:"deadline"`
	TaskType  string    `json:"type"`
}

var taskFile = os.ExpandEnv("$HOME/.local/share/enhancedTodo/data.json")

func main() {
	// Start HTTP server for static files
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./assets")))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Initialize webview
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Enhanced Todo List")
	w.SetSize(800, 600, webview.HintNone)

	// Bind functions for JavaScript
	w.Bind("getTasks", getTasks)
	w.Bind("addTask", addTask)
	w.Bind("toggleTaskDone", toggleTaskDone)

	// Navigate to local server URL
	w.Navigate("http://localhost:8080/index.html")
	w.Run()
}

// Load tasks from JSON file
func loadTasks() []Task {
	file, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Println("No existing data file found, starting with empty list.")
		return []Task{}
	}
	var tasks []Task
	json.Unmarshal(file, &tasks)
	return tasks
}

// Save tasks to JSON file
func saveTasks(tasks []Task) {
	file, _ := json.MarshalIndent(tasks, "", "  ")
	ioutil.WriteFile(taskFile, file, 0644)
}

// Get tasks as JSON string
func getTasks() ([]Task, error) {
	return loadTasks(), nil
}

// Add a new task
func addTask(title string, deadline string, taskType string) error {
	tasks := loadTasks()
	newTask := Task{
		ID:        len(tasks) + 1,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
		Deadline:  parseDate(deadline),
		TaskType:  taskType,
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	return nil
}

// Toggle task completion
func toggleTaskDone(id int) error {
	tasks := loadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = !tasks[i].Done
			break
		}
	}
	saveTasks(tasks)
	return nil
}

// Parse date from string to time.Time
func parseDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return parsedDate
}
