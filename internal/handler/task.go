package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gectou4/rest_api_go/internal/model"
)

type TaskHandler struct {
	DB *sql.DB
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, "Invalid request")
			return
		}
		input.Title = r.FormValue("title")
		input.Description = r.FormValue("description")
		if s := r.FormValue("status"); s != "" {
			var status int
			json.Unmarshal([]byte(s), &status)
			input.Status = &status
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			if err := r.ParseForm(); err == nil {
				input.Title = r.FormValue("title")
				input.Description = r.FormValue("description")
				if s := r.FormValue("status"); s != "" {
					var status int
					json.Unmarshal([]byte(s), &status)
					input.Status = &status
				}
			}
		}
	}

	if input.Title == "" {
		writeJSON(w, http.StatusBadRequest, "Title is required")
		return
	}

	status := model.StatusBacklog
	if input.Status != nil {
		status = model.ParseTaskStatus(*input.Status)
	}

	loc, _ := time.LoadLocation("Europe/Paris")
	task := model.NewTask(h.DB)
	task.Title = input.Title
	task.Description = input.Description
	task.Status = status
	task.CreationDate = time.Now().In(loc)

	saved, err := task.Save()
	if err != nil || !saved {
		writeJSON(w, http.StatusInternalServerError, "Unable to create new Task")
		return
	}

	writeJSON(w, http.StatusCreated, task.ToMap())
}

func (h *TaskHandler) AddTaskToUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userID := parseInt(params["userId"])
	taskID := parseInt(params["taskId"])

	user := model.NewUser(h.DB, userID)
	if !user.IsLoaded() {
		writeJSON(w, http.StatusBadRequest, "User not exists")
		return
	}

	task := model.NewTask(h.DB, taskID)
	if !task.IsLoaded() {
		writeJSON(w, http.StatusBadRequest, "Task not exists")
		return
	}

	userTask := model.NewUserTask(h.DB, userID)
	userTask.AddTaskID(taskID)
	saved, err := userTask.Save()
	if err != nil || !saved {
		writeJSON(w, http.StatusInternalServerError, "Unable to add Task to user")
		return
	}

	writeJSON(w, http.StatusOK, 1)
}

func (h *TaskHandler) EditTask(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := parseInt(params["id"])
	if id <= 0 {
		writeJSON(w, http.StatusBadRequest, "Id of task to edit is required")
		return
	}

	task := model.NewTask(h.DB, id)
	if !task.IsLoaded() {
		writeJSON(w, http.StatusBadRequest, "Task not found")
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *int    `json:"status"`
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			if v := r.FormValue("title"); v != "" {
				input.Title = &v
			}
			if v := r.FormValue("description"); v != "" {
				input.Description = &v
			}
			if v := r.FormValue("status"); v != "" {
				var s int
				json.Unmarshal([]byte(v), &s)
				input.Status = &s
			}
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			if err := r.ParseForm(); err == nil {
				if v := r.FormValue("title"); v != "" {
					input.Title = &v
				}
				if v := r.FormValue("description"); v != "" {
					input.Description = &v
				}
				if v := r.FormValue("status"); v != "" {
					var s int
					json.Unmarshal([]byte(v), &s)
					input.Status = &s
				}
			}
		}
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = model.ParseTaskStatus(*input.Status)
	}

	saved, err := task.Save()
	if err != nil || !saved {
		writeJSON(w, http.StatusInternalServerError, "Unable to update Task")
		return
	}

	writeJSON(w, http.StatusOK, 1)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := parseInt(params["id"])
	task := model.NewTask(h.DB, id)
	if !task.IsLoaded() {
		writeJSON(w, http.StatusBadRequest, "Task not exists")
		return
	}

	deleted, err := task.Delete()
	if err != nil || !deleted {
		writeJSON(w, http.StatusInternalServerError, "Unable to delete Task")
		return
	}

	writeJSON(w, http.StatusOK, 1)
}

func (h *TaskHandler) DeleteUserTask(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userID := parseInt(params["userId"])
	taskID := parseInt(params["taskId"])

	user := model.NewUser(h.DB, userID)
	if !user.IsLoaded() {
		writeJSON(w, http.StatusBadRequest, "User not exists")
		return
	}

	userTask := model.NewUserTask(h.DB, userID)

	if userTask.HasTask(taskID) {
		userTask.RemoveTaskID(taskID)
		saved, err := userTask.Save()
		if err != nil || !saved {
			writeJSON(w, http.StatusInternalServerError, "Unable to delete Task of user")
			return
		}
	}

	writeJSON(w, http.StatusOK, 1)
}
