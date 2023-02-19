package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const ConfigPath = "./config"

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

type dbConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Server   string `toml:"server"`
	Port     string `toml:"port"`
	Name     string `toml:"dbname"`
}

type serverConfig struct {
	Addr   string
	Port   string
	Secure bool
}

type Config struct {
	Env    string
	Server serverConfig `toml:"server"`
	DB     dbConfig     `toml:"database"`
}

func NewConfig(env string) Config {
	config := Config{
		Env: env,
	}
	path := fmt.Sprintf("%s/%s.toml", ConfigPath, config.Env)
	bs, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = toml.Decode(string(bs), &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}

func (cfg Config) DBConnectionString() string {
	// TODO deal with ssl
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Server,
		cfg.DB.Port,
		cfg.DB.Name,
	)
}

func (cfg Config) ServerConnectionString() string {
	return fmt.Sprintf("%s:%s",
		cfg.Server.Addr,
		cfg.Server.Port,
	)
}
