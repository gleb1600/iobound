package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string    `json:"id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskStore struct {
	sync.RWMutex
	tasks map[string]*Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*Task),
	}
}

func (ts *TaskStore) CreateTask(input string) *Task {
	ts.Lock()
	defer ts.Unlock()

	task := &Task{
		ID:        uuid.New().String()[:8],
		Input:     input,
		Output:    "<no output yet>",
		Status:    "in work",
		CreatedAt: time.Now().UTC(),
	}

	ts.tasks[task.ID] = task
	return task
}

func (ts *TaskStore) CompleteTask(id string) {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks[id].Status = "completed"
	ts.tasks[id].Output = "some result"
}

func (ts *TaskStore) GetTask(id string) (*Task, bool) {
	ts.RLock()
	defer ts.RUnlock()

	task, exists := ts.tasks[id]
	return task, exists
}

func processTask(task *Task, store *TaskStore) {
	// Эмитация длительной операции
	time.Sleep(30 * time.Second)

	store.CompleteTask(task.ID)
}

func CreateHandler(taskStore *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Input string `json:"input"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		task := taskStore.CreateTask(request.Input)

		go processTask(task, taskStore)

		fmt.Fprintf(w, "Task started.\nTask ID:%s", task.ID)
	}
}

func FindTaskHandler(taskStore *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[1:]
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		task, exists := taskStore.GetTask(id)
		if !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		ShowTask(w, task)
	}
}

func ShowTasksHandler(taskStore *TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		for _, t := range taskStore.tasks {
			ShowTask(w, t)
		}
	}
}

func ShowTask(w http.ResponseWriter, task *Task) {
	fmt.Fprintf(w, "\nTask ID:      %s\nStatus:       %s\nInput:        %s\nOutput:       %s\nCreated at:   %s\n",
		task.ID,
		task.Status,
		task.Input,
		task.Output,
		task.CreatedAt.Format("02.01.2006 15:04:05"),
	)
}

func main() {
	taskStore := NewTaskStore()

	http.HandleFunc("/create", CreateHandler(taskStore))
	http.HandleFunc("/", FindTaskHandler(taskStore))
	http.HandleFunc("/tasks", ShowTasksHandler(taskStore))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
