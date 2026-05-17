package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gectou4/rest_api_go/internal/api"
	"github.com/gectou4/rest_api_go/internal/config"
	"github.com/gectou4/rest_api_go/internal/model"
	"github.com/gectou4/rest_api_go/internal/router"
)

var testDB *sql.DB
var testServer *httptest.Server

func TestMain(m *testing.M) {
	dbCfg := config.LoadDBConfig()

	var err error
	testDB, err = sql.Open("mysql", dbCfg.ConnectionString())
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	if err := testDB.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	routes := config.LoadRoutes()
	r := router.NewRouter(routes)
	app := api.NewAPI(testDB)

	mux := http.NewServeMux()
	mux.Handle("/", r.Middleware(app.Handler()))

	testServer = httptest.NewServer(mux)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}

func TestGetUser(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/user/1")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := data["user_id"]; !ok {
		t.Fatal("Response missing 'user_id' key")
	}

	if int(data["user_id"].(float64)) != 1 {
		t.Fatalf("Expected user_id 1, got %v", data["user_id"])
	}
}

func TestGetUserTask(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/user/1/task")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := data["user_id"]; !ok {
		t.Fatal("Response missing 'user_id' key")
	}
	if _, ok := data["tasks"]; !ok {
		t.Fatal("Response missing 'tasks' key")
	}

	if int(data["user_id"].(float64)) != 1 {
		t.Fatalf("Expected user_id 1, got %v", data["user_id"])
	}
}

func TestAddTask(t *testing.T) {
	payload := map[string]interface{}{
		"title":       "Faire le thè",
		"description": "Comme pour le café, mais avec du thé",
		"status":      int(model.StatusBacklog),
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(testServer.URL+"/task", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := data["task_id"]; !ok {
		t.Fatal("Response missing 'task_id' key")
	}
	if _, ok := data["title"]; !ok {
		t.Fatal("Response missing 'title' key")
	}
}

func getLastTaskID(t *testing.T) int {
	tasks, err := model.NewTask(testDB).GetAll()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("No tasks found")
	}
	last := tasks[len(tasks)-1]
	switch v := last["task_id"].(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		t.Fatalf("Unexpected type for task_id: %T", v)
		return 0
	}
}

func TestEditTask(t *testing.T) {
	taskID := getLastTaskID(t)

	payload := map[string]interface{}{
		"title":       "Faire le thè",
		"description": "Comme pour le café, mais avec du thé et en mieux",
		"status":      int(model.StatusBacklog),
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(fmt.Sprintf("%s/task/%d", testServer.URL, taskID), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data int
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data != 1 {
		t.Fatalf("Expected response 1, got %d", data)
	}
}

func TestAddTaskToUser(t *testing.T) {
	taskID := getLastTaskID(t)

	resp, err := http.Post(fmt.Sprintf("%s/user/1/task/%d", testServer.URL, taskID), "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data int
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data != 1 {
		t.Fatalf("Expected response 1, got %d", data)
	}
}

func TestDelTaskToUser(t *testing.T) {
	taskID := getLastTaskID(t)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/user/1/task/%d", testServer.URL, taskID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data int
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data != 1 {
		t.Fatalf("Expected response 1, got %d", data)
	}
}

func TestDelTask(t *testing.T) {
	taskID := getLastTaskID(t)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/task/%d", testServer.URL, taskID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var data int
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data != 1 {
		t.Fatalf("Expected response 1, got %d", data)
	}
}
