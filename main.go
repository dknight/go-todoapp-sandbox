package main

// Hot reload
// go run github.com/cosmtrek/air

import (
	"log"
	"os"

	"github.com/dknight/go-todoapp-sandbox/controllers"
	"github.com/dknight/go-todoapp-sandbox/lib"
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
	config *lib.Config
	env    *lib.Env
)

func init() {
	envVar := os.Getenv("ENV")
	config = lib.NewConfig(envVar)

	env = &lib.Env{
		Env:    envVar,
		DB:     lib.NewDB(config),
		Logger: lib.NewLogger(logDir),
	}
	env.Logger.Println("App started")

	err := env.DB.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connection to database succeeded.")
	}
}

func main() {
	defer env.DB.Close()

	engine := html.New(viewsDir, ".gohtml")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", publicDir)

	todoItemController := controllers.NewTodoController(env)
	systemController := controllers.NewSystemController(env)

	app.Get("/", todoItemController.Index)
	app.Post("/items", todoItemController.Post)
	app.Put("/items/:id", todoItemController.Put)
	app.Delete("/items/:id", todoItemController.Delete)

	app.Get("/ping", systemController.Ping)
	app.Get("/instance", systemController.Instance)
	log.Fatalln(app.Listen(config.ServerConnectionString()))
}
