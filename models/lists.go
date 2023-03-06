package models

import (
	"database/sql"
	"errors"
)

type List struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// TODO
func ListLists(db *sql.DB) ([]Lists, error) {
	return nil, errors.New("")
}

// TODO
func FindList(db *sql.DB, id int) (Lists, error) {
	return nil, errors.New("")
}

// TODO
func (li List) Create(db *sql.DB) (int64, error) {
	return 0, errors.New("")
}

// TODO
func (li List) Save(db *sql.DB) (int64, error) {
	return 0, errors.New("")
}

// TODO
func (li List) Delete(db *sql.DB) (int64, error) {
	return 0, errors.New("")
}
