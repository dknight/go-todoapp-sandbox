package lib

import (
	"database/sql"
	"log"
)

func NewDB(cfg *Config) *sql.DB {
	dbConnString := cfg.DBConnectionString()
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
