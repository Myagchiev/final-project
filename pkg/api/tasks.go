// pkg/api/tasks.go
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Myagchiev/final-project/pkg/db"
	"github.com/Myagchiev/final-project/pkg/utils"
)

const maxTasks = 50

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func checkAndFixDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(utils.DateLayout)
	truncatedNow := now.Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = today
		return nil
	}

	parsed, err := time.Parse(utils.DateLayout, task.Date)
	if err != nil {
		return err
	}

	if !parsed.Before(truncatedNow) {
		return nil
	}

	if task.Repeat == "" {
		task.Date = today
		return nil
	}

	next, err := utils.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	task.Date = next
	return nil
}

func tasksListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tasks := make([]db.Task, 0, maxTasks)
	var err error

	if search == "" {
		tasks, err = db.Tasks(maxTasks)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, TasksResp{Tasks: tasks})
		return
	}

	if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		searchDate := t.Format(utils.DateLayout)
		tasks, err = db.TasksWithFilter(maxTasks, "", searchDate)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, TasksResp{Tasks: tasks})
		return
	}

	tasks, err = db.TasksWithFilter(maxTasks, search, "")
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONError(w, "title is empty", http.StatusBadRequest)
		return
	}
	if err := checkAndFixDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		writeJSONError(w, "id is empty", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}
	if task.ID == 0 {
		writeJSONError(w, "id is empty", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONError(w, "title is empty", http.StatusBadRequest)
		return
	}
	if err := checkAndFixDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]interface{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		writeJSONError(w, "id is empty", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := db.MarkDone(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		writeJSONError(w, "id is empty", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]interface{}{})
}

func taskCRUDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/api/task":
			addTaskHandler(w, r)
		case "/api/task/done":
			doneTaskHandler(w, r)
		default:
			writeJSONError(w, "not found", http.StatusNotFound)
		}

	case http.MethodGet:
		if r.URL.Path == "/api/task" {
			getTaskHandler(w, r)
			return
		}
		writeJSONError(w, "not found", http.StatusNotFound)

	case http.MethodPut:
		if r.URL.Path == "/api/task" {
			updateTaskHandler(w, r)
			return
		}
		writeJSONError(w, "not found", http.StatusNotFound)

	case http.MethodDelete:
		if r.URL.Path == "/api/task" {
			deleteTaskHandler(w, r)
			return
		}
		writeJSONError(w, "not found", http.StatusNotFound)

	default:
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}