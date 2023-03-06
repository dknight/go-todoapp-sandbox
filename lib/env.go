package lib

import (
	"database/sql"
	"log"
)

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

type Env struct {
	Env    string
	DB     *sql.DB
	Logger *log.Logger
}

func NewEnv() *Env {
	return &Env{}
}
