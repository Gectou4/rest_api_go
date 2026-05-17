package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gectou4/rest_api_go/internal/model"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetIndex(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := parseInt(params["id"])
	user := model.NewUser(h.DB, id)
	if !user.IsLoaded() {
		writeJSON(w, http.StatusNotFound, "No user found")
		return
	}
	writeJSON(w, http.StatusOK, user.ToMap())
}

func (h *UserHandler) GetUserTask(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := parseInt(params["id"])
	user := model.NewUser(h.DB, id)
	if !user.IsLoaded() {
		writeJSON(w, http.StatusNotFound, "No user found")
		return
	}
	userTask := user.GetTask(h.DB)
	writeJSON(w, http.StatusOK, userTask.ToMap())
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
