package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	// Initialize webview
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Enhanced Todo List")
	w.SetSize(800, 600, webview.HintNone)

	// Load tasks and setup event handlers
	w.Bind("getTasks", getTasks)
	w.Bind("addTask", addTask)
	w.Bind("toggleTaskDone", toggleTaskDone)

	// Load HTML with inlined CSS and JS
	var content = loadHTML()
	w.Navigate("data:text/html," + content)
	w.Run()
}

// Load and parse tasks from JSON file
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
func loadHTML() string {
	// Load HTML content
	htmlFile, err := ioutil.ReadFile("./assets/index.html")
	if err != nil {
		log.Fatalf("Failed to load HTML file: %v", err)
	}
	htmlContent := string(htmlFile)

	// Load and sanitize CSS content
	cssFile, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.config/enhancedTodoList/style.css"))
	if err != nil {
		log.Fatalf("Failed to load CSS file: %v", err)
	}
	cssContent := string(cssFile)

	// Add debugging logs for CSS length
	log.Printf("CSS content length: %d characters", len(cssContent))

	// Load and sanitize JavaScript content
	jsFile, err := ioutil.ReadFile("./assets/main.js")
	if err != nil {
		log.Fatalf("Failed to load JS file: %v", err)
	}
	jsContent := string(jsFile)

	log.Printf("JS content length: %d characters", len(jsContent))

	// Combine HTML with CSS and JavaScript, using placeholders
	fullContent := fmt.Sprintf(
		`<style>%s</style>%s<script>%s</script>`,
		cssContent,
		htmlContent,
		jsContent,
	)

	// Add debugging to check final HTML length
	log.Printf("Full HTML content length: %d characters", len(fullContent))

	return fullContent
}

// Parse date from string to time.Time
func parseDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return parsedDate
}
