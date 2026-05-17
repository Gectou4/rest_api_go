package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID           int
	Title        string
	Description  string
	CreationDate time.Time
	Status       TaskStatus
	db           *sql.DB
	table        string
	loaded       bool
}

func NewTask(db *sql.DB, id ...int) *Task {
	t := &Task{
		db:    db,
		table: "task",
	}
	if len(id) > 0 && id[0] > 0 {
		t.Load(id[0])
	}
	return t
}

func (t *Task) IsLoaded() bool {
	return t.loaded
}

func (t *Task) Load(id int) {
	if t.loaded {
		return
	}
	t.ID = id
	query := fmt.Sprintf("SELECT status, title, description, creation_date FROM `%s` WHERE task_id = ?", t.table)
	var status int
	var creationDate time.Time
	err := t.db.QueryRow(query, t.ID).Scan(&status, &t.Title, &t.Description, &creationDate)
	if err == nil {
		t.Status = ParseTaskStatus(status)
		t.CreationDate = creationDate
		t.loaded = true
	}
}

func (t *Task) Save() (bool, error) {
	var err error
	if t.ID <= 0 {
		query := fmt.Sprintf("INSERT INTO `%s` (status, title, description, creation_date) VALUES (?, ?, ?, ?)", t.table)
		result, execErr := t.db.Exec(query, t.Status, t.Title, t.Description, t.CreationDate.Format("2006-01-02 15:04:05"))
		if execErr != nil {
			return false, execErr
		}
		lastID, idErr := result.LastInsertId()
		if idErr != nil {
			return false, idErr
		}
		t.ID = int(lastID)
		return true, nil
	}
	query := fmt.Sprintf("UPDATE `%s` SET status=?, title=?, description=?, creation_date=? WHERE task_id = ?", t.table)
	_, err = t.db.Exec(query, t.Status, t.Title, t.Description, t.CreationDate.Format("2006-01-02 15:04:05"), t.ID)
	return err == nil, err
}

func (t *Task) Delete() (bool, error) {
	query := fmt.Sprintf("DELETE FROM `%s` WHERE task_id = ?", t.table)
	_, err := t.db.Exec(query, t.ID)
	return err == nil, err
}

func (t *Task) GetAll() ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT task_id, status, title, description, creation_date FROM `%s`", t.table)
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id, status int
		var title, description string
		var creationDate time.Time
		if scanErr := rows.Scan(&id, &status, &title, &description, &creationDate); scanErr != nil {
			continue
		}
		tasks = append(tasks, map[string]interface{}{
			"task_id":       id,
			"status":        status,
			"title":         title,
			"description":   description,
			"creation_date": creationDate.Format("2006-01-02 15:04:05"),
		})
	}
	return tasks, rows.Err()
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"task_id":       t.ID,
		"status":        int(t.Status),
		"title":         t.Title,
		"description":   t.Description,
		"creation_date": t.CreationDate.Format("2006-01-02 15:04:05"),
	}
}
