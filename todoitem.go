package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/blockloop/scan"
)

var (
	ErrTaskNameEmpty    = errors.New("Task name is empty")
	ErrItemNotFound     = errors.New("Item not found")
	ErrCannotDeleteItem = errors.New("Cannot delete the item")
)

type TodoItem struct {
	ID        int       `db:"id"`
	Task      string    `db:"task"`
	Status    bool      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func ListTodoItems(db *sql.DB) ([]TodoItem, error) {
	var items []TodoItem
	rows, err := db.Query("SELECT id, task, status, created_at FROM " +
		"items ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = scan.Rows(&items, rows)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func FindItem(db *sql.DB, id string) *TodoItem {
	var item TodoItem
	rows, err := db.Query("SELECT id, task, status, created_at FROM items "+
		"WHERE id = $1", id)
	if err != nil {
		return nil
	}
	err = scan.Row(&item, rows)
	if err != nil {
		return nil
	}
	return &item
}

func (ti TodoItem) Create(db *sql.DB) (int64, error) {
	var id int64
	if ti.Task != "" {
		err := db.QueryRow(
			"INSERT INTO items (task, status) VALUES ($1, $2) RETURNING id",
			ti.Task, ti.Status).Scan(&id)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		return id, nil
	}
	return 0, ErrTaskNameEmpty
}

func (ti TodoItem) Save(db *sql.DB) error {
	_, err := db.Exec("UPDATE items SET task=$1, status=$2 WHERE id=$3",
		ti.Task, ti.Status, ti.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ti TodoItem) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM items WHERE id=$1", ti.ID)
	if err != nil {
		return ErrCannotDeleteItem
	}
	return nil
}

// func (ti TodoItem) JSON() ([]byte, error) {
// 	var bs []byte
// 	bs, err := json.Marshal(ti)
// 	return bs, err
// }
