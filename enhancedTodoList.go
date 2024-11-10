package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"os"
	"time"

	webview "github.com/webview/webview_go"
)

// Task structure for a todo
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	Deadline    time.Time  `json:"deadline"`
	TaskType    string     `json:"type"`
	CompletedAt *time.Time `json:"completed_at"` // New field to track completion date
}

var (
	taskFile    = os.ExpandEnv("$HOME/.local/share/enhancedTodo/data.json")
	configDir   = os.ExpandEnv("$HOME/.config/enhancedTodoList") // User config directory
	cssFilePath = filepath.Join(configDir, "style.css")          // Path to user's CSS file
)

func main() {
	// Start HTTP server to serve files
	go func() {
		// Serve assets directory
		http.Handle("/", http.FileServer(http.Dir("./assets")))

		// Serve user's CSS from ~/.config/enhancedTodoList/style.css
		http.HandleFunc("/user-style.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, cssFilePath)
		})

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
	w.Bind("deleteTaskGo", deleteTask)

	// Navigate to local server URL
	w.Navigate("http://localhost:8080/index.html")
	w.Run()
}

// Load tasks from JSON file
func loadTasks() []Task {
	taskDir := filepath.Dir(taskFile)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", taskDir, err)
	}

	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		if err := ioutil.WriteFile(taskFile, []byte("[]"), 0644); err != nil {
			log.Fatalf("Failed to create file %s: %v", taskFile, err)
		}
	}

	file, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", taskFile, err)
	}

	var tasks []Task
	json.Unmarshal(file, &tasks)

	// Filter out tasks completed more than a week ago
	now := time.Now()
	updatedTasks := tasks[:0]
	for _, task := range tasks {
		if task.Done && task.CompletedAt != nil {
			if now.Sub(*task.CompletedAt) > 24*time.Hour {
				continue // Skip tasks older than one week
			}
		}
		updatedTasks = append(updatedTasks, task)
	}

	// Save the updated task list if any tasks were removed
	if len(updatedTasks) != len(tasks) {
		saveTasks(updatedTasks)
	}

	return updatedTasks
}

func deleteTask(id int) error {
	tasks := loadTasks()
	updatedTasks := tasks[:0]
	for _, task := range tasks {
		if task.ID != id {
			updatedTasks = append(updatedTasks, task)
		}
	}
	saveTasks(updatedTasks)
	return nil
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
			tasks[i].Done = !task.Done
			if tasks[i].Done {
				now := time.Now()
				tasks[i].CompletedAt = &now
			} else {
				tasks[i].CompletedAt = nil
			}
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
