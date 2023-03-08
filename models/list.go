package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrListNameEmpty         = errors.New("List name is empty")
	ErrListNotFound          = errors.New("List not found")
	ErrListCannotDelete      = errors.New("Cannot delete the list")
	ErrListCannotDeleteItems = errors.New("Cannot delete the list's items")
)

type List struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	ItemsCount int    `db:"count"` // NOTE sqlx is better here?
}

func ListLists(db *sql.DB) ([]List, error) {
	var lists []List
	rows, err := db.Query(`
		SELECT lists.id, lists.name, COUNT(items.task)
		FROM lists
		LEFT OUTER JOIN items
		ON items.list_id=lists.id
		GROUP BY lists.id, lists.name;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		list := List{}
		err := rows.Scan(&list.ID, &list.Name, &list.ItemsCount)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func FindList(db *sql.DB, id int) (*List, error) {
	list := List{}
	row := db.QueryRow(`
		SELECT id, name
		FROM lists
		WHERE id = $1`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.Scan(&list.ID, &list.Name)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (li List) Create(db *sql.DB) (int64, error) {
	var id int64
	if li.Name != "" {
		// Postgres doesnt support LastInsertedID(), so we user QueryRow,
		// instead of Exec.
		err := db.QueryRow(
			`INSERT INTO lists (name) VALUES ($1) RETURNING id`,
			li.Name).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, ErrListNameEmpty
}

func (li List) Save(db *sql.DB) error {
	_, err := db.Exec(`UPDATE lists SET name=$1 WHERE id=$2`,
		li.Name, li.ID)
	if err != nil {
		return err
	}
	return nil
}

func (li List) Delete(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM items WHERE list_id=$1`, li.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = tx.Exec(`DELETE FROM lists WHERE id=$1`, li.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
