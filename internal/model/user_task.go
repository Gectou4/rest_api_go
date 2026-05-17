package model

import (
	"database/sql"
	"fmt"
)

type UserTask struct {
	UserID   int
	TaskList map[int]*Task
	db       *sql.DB
	table    string
	loaded   bool
}

func NewUserTask(db *sql.DB, userID ...int) *UserTask {
	ut := &UserTask{
		db:       db,
		table:    "user_task",
		TaskList: make(map[int]*Task),
	}
	if len(userID) > 0 && userID[0] > 0 {
		ut.Load(userID[0])
	}
	return ut
}

func (ut *UserTask) Load(id int) {
	if ut.loaded {
		return
	}
	ut.LoadByUserID(id)
}

func (ut *UserTask) GetUserID() int {
	return ut.UserID
}

func (ut *UserTask) SetUserID(id int) {
	ut.UserID = id
}

func (ut *UserTask) GetTaskIDs() []int {
	ids := make([]int, 0, len(ut.TaskList))
	for id := range ut.TaskList {
		ids = append(ids, id)
	}
	return ids
}

func (ut *UserTask) AddTaskID(taskID int) *UserTask {
	ut.TaskList[taskID] = NewTask(ut.db, taskID)
	return ut
}

func (ut *UserTask) RemoveTaskID(taskID int) *UserTask {
	delete(ut.TaskList, taskID)
	return ut
}

func (ut *UserTask) HasTask(taskID int) bool {
	_, ok := ut.TaskList[taskID]
	return ok
}

func GetTaskByUser(db *sql.DB, user *User) *UserTask {
	ut := NewUserTask(db)
	ut.LoadByUser(user)
	return ut
}

func (ut *UserTask) Save() (bool, error) {
	tx, err := ut.db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM `user_task` WHERE user_id = ?", ut.UserID)
	if err != nil {
		return false, err
	}

	stmt, prepErr := tx.Prepare("INSERT INTO `user_task` (`user_id`, `task_id`) VALUES (?, ?)")
	if prepErr != nil {
		return false, prepErr
	}
	defer stmt.Close()

	for taskID := range ut.TaskList {
		if _, execErr := stmt.Exec(ut.UserID, taskID); execErr != nil {
			return false, execErr
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return false, commitErr
	}
	return true, nil
}

func (ut *UserTask) DeleteUserTask(taskID int) (bool, error) {
	tx, err := ut.db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM `user_task` WHERE user_id = ? AND task_id = ?", ut.UserID, taskID)
	if err != nil {
		return false, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return false, commitErr
	}
	return true, nil
}

func (ut *UserTask) ToMap() map[string]interface{} {
	tasks := make(map[string]interface{})
	for taskID, task := range ut.TaskList {
		tasks[fmt.Sprintf("%d", taskID)] = task.ToMap()
	}
	return map[string]interface{}{
		"user_id": ut.UserID,
		"tasks":   tasks,
	}
}

func (ut *UserTask) LoadByUser(user *User) {
	ut.LoadByUserID(user.ID)
}

func (ut *UserTask) LoadByUserID(userID int) {
	if ut.loaded {
		return
	}
	ut.UserID = userID
	rows, err := ut.db.Query("SELECT task_id FROM user_task WHERE user_id = ?", ut.UserID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var taskID int
		if scanErr := rows.Scan(&taskID); scanErr == nil {
			ut.AddTaskID(taskID)
		}
	}
	ut.loaded = true
}
