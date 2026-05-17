package model

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID    int
	Name  string
	Email string
	db    *sql.DB
	table string
	loaded bool
}

func NewUser(db *sql.DB, id ...int) *User {
	u := &User{
		db:    db,
		table: "user",
	}
	if len(id) > 0 && id[0] > 0 {
		u.Load(id[0])
	}
	return u
}

func (u *User) IsLoaded() bool {
	return u.loaded
}

func (u *User) Load(id int) {
	if u.loaded {
		return
	}
	u.ID = id
	query := fmt.Sprintf("SELECT email, name FROM `%s` WHERE user_id = ?", u.table)
	err := u.db.QueryRow(query, u.ID).Scan(&u.Email, &u.Name)
	if err == nil {
		u.loaded = true
	}
}

func (u *User) GetTask(db *sql.DB) *UserTask {
	return GetTaskByUser(db, u)
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id": u.ID,
		"name":    u.Name,
		"email":   u.Email,
	}
}
