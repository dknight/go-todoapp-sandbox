package models

import (
	"database/sql"
	"errors"
	"time"
	// For easier scan can be used, but decreases performance.
	// "github.com/blockloop/scan"
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
	ListId    int       `db:"list_id"`
}

func ListTodoItems(db *sql.DB) ([]TodoItem, error) {
	var items []TodoItem
	rows, err := db.Query(`
		SELECT id, task, status, created_at
		FROM items
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := TodoItem{}
		err := rows.Scan(&item.ID, &item.Task, &item.Status, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// FIXME id - int
func FindItem(db *sql.DB, id string) (*TodoItem, error) {
	item := TodoItem{}
	row := db.QueryRow(`
		SELECT id, task, status, created_at
		FROM items
		WHERE id = $1`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.Scan(&item.ID, &item.Task, &item.Status, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (ti TodoItem) Create(db *sql.DB) (int64, error) {
	var id int64
	if ti.Task != "" {
		// Postgres doesnt support LastInsertedID(), so we user QueryRow,
		// instead of Exec.
		err := db.QueryRow(
			`INSERT INTO items (task, status) VALUES ($1, $2) RETURNING id`,
			ti.Task, ti.Status).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, ErrTaskNameEmpty
}

func (ti TodoItem) Save(db *sql.DB) error {
	_, err := db.Exec(`UPDATE items SET task=$1, status=$2 WHERE id=$3`,
		ti.Task, ti.Status, ti.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ti TodoItem) Delete(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM items WHERE id=$1`, ti.ID)
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
