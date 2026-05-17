package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gectou4/rest_api_go/internal/handler"
	"github.com/gectou4/rest_api_go/internal/router"
)

type API struct {
	DB          *sql.DB
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
}

func NewAPI(db *sql.DB) *API {
	return &API{
		DB:          db,
		UserHandler: &handler.UserHandler{DB: db},
		TaskHandler: &handler.TaskHandler{DB: db},
	}
}

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleRequest)
	return mux
}

func (a *API) handleRequest(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	controller := router.GetController(r.Context())
	action := router.GetAction(r.Context())
	params := router.GetParams(r.Context())

	if controller == "" || action == "" {
		writeJSON(w, http.StatusNotFound, "Not Found")
		return
	}

	switch controller {
	case "User":
		switch action {
		case "GetIndex":
			a.UserHandler.GetIndex(w, r, params)
			return
		case "GetUserTask":
			a.UserHandler.GetUserTask(w, r, params)
			return
		}
	case "Task":
		switch action {
		case "AddTask":
			a.TaskHandler.AddTask(w, r, params)
			return
		case "AddTaskToUser":
			a.TaskHandler.AddTaskToUser(w, r, params)
			return
		case "EditTask":
			a.TaskHandler.EditTask(w, r, params)
			return
		case "DeleteTask":
			a.TaskHandler.DeleteTask(w, r, params)
			return
		case "DeleteUserTask":
			a.TaskHandler.DeleteUserTask(w, r, params)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, "Not Found")
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(data)
}
