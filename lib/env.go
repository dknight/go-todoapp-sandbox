package lib

import (
	"database/sql"
	"log"
)

type Env struct {
	Env    string
	DB     *sql.DB
	Logger *log.Logger
}

func NewEnv() *Env {
	return &Env{}
}
