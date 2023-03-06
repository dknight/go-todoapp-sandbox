package models

import (
	"database/sql"
	"errors"
)

var (
	ErrListNameEmpty    = errors.New("List name is empty")
	ErrListNotFound     = errors.New("List not found")
	ErrCannotDeleteList = errors.New("Cannot delete the list")
)

type List struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func ListLists(db *sql.DB) ([]Lists, error) {
	var lists []Lists
	rows, err := db.Query(`SELECT id, name FROM lists;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		list := List{}
		err := rows.Scan(&list.ID, &list.Name)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func FindList(db *sql.DB, id int) (List, error) {
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
	if ti.Name != "" {
		// Postgres doesnt support LastInsertedID(), so we user QueryRow,
		// instead of Exec.
		err := db.QueryRow(
			`INSERT INTO lists (name) VALUES ($1) RETURNING id`,
			li.Task).Scan(&id)
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
	_, err := db.Exec(`DELETE FROM lists WHERE id=$1`, li.ID)
	if err != nil {
		return ErrCannotDeleteList
	}
	return nil
}
