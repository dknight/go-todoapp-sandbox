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
	ListID    int       `db:"list_id"`
}

func ListTodoItems(db *sql.DB) ([]TodoItem, error) {
	var items []TodoItem
	rows, err := db.Query(`
		SELECT id, task, status, created_at, list_id
		FROM items
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := TodoItem{}
		err := rows.Scan(
			&item.ID,
			&item.Task,
			&item.Status,
			&item.CreatedAt,
			&item.ListID,
		)
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

func FindItemsByListID(db *sql.DB, listID int) ([]TodoItem, error) {
	var items []TodoItem
	rows, err := db.Query(`
		SELECT items.id, task, status, created_at, list_id
		FROM items
		INNER JOIN lists
		ON lists.id=items.list_id
		WHERE items.list_id=$1
		ORDER BY created_at DESC`, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := TodoItem{}
		err := rows.Scan(
			&item.ID,
			&item.Task,
			&item.Status,
			&item.CreatedAt,
			&item.ListID,
		)
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

func FindItem(db *sql.DB, id int) (*TodoItem, error) {
	item := TodoItem{}
	row := db.QueryRow(`
		SELECT id, task, status, created_at, list_id
		FROM items
		WHERE id = $1`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.Scan(
		&item.ID,
		&item.Task,
		&item.Status,
		&item.CreatedAt,
		&item.ListID,
	)
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
			`INSERT INTO items (task, status, list_id) VALUES
			($1, $2, $3) RETURNING id`,
			ti.Task, ti.Status, ti.ListID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, ErrTaskNameEmpty
}

func (ti TodoItem) Save(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE items
		SET task=$1, status=$2, list_id=$3
		WHERE id=$4`,
		ti.Task, ti.Status, ti.ListID, ti.ID)
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
