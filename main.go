package main

// Hot reload
// go run github.com/cosmtrek/air

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dknight/go-todoapp-sandbox/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	_ "github.com/lib/pq"
)

const (
	viewsDir  = "./views"
	publicDir = "./public"
	logDir    = "./logs"
)

var (
	dbConnString     string
	serverConnString string
	port             string
	db               *sql.DB
	config           Config
	logger           *log.Logger // TODO make err, info, warn, err logger
)

func init() {
	var err error
	env := os.Getenv("ENV")
	if env == "" {
		env = EnvDev
	}
	config = NewConfig(env)
	dbConnString = config.DBConnectionString()
	serverConnString = config.ServerConnectionString()
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	logfile := fmt.Sprintf(
		"%s/%s.log", logDir, time.Now().Format("2006-01-02"),
	)
	fp, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(fp, "", log.LstdFlags)
	logger.Println("App started")
}

func main() {
	engine := html.New(viewsDir, ".gohtml")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", publicDir)

	todoItemController := controllers.NewTodoController(db, logger)
	systemController := controllers.NewSystemController()

	app.Get("/", todoItemController.Index)
	app.Post("/items", todoItemController.Post)
	app.Put("/items/:id", todoItemController.Put)
	app.Delete("/items/:id", todoItemController.Delete)

	app.Get("/ping", systemController.Ping)
	app.Get("/instance", systemController.Instance)
	log.Fatalln(app.Listen(serverConnString))
}
