package main

// Hot reload
// go run github.com/cosmtrek/air

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	logfile := fmt.Sprintf("%s/%s.log", logDir, time.Now().Format("2006-01-02"))
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

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})
	app.Post("/items", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})
	app.Put("/items/:id", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})
	app.Delete("/items/:id", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	app.Get("/ping", ping)
	app.Get("/instance", instance)
	log.Fatalln(app.Listen(serverConnString))
}

func indexHandler(ctx *fiber.Ctx, db *sql.DB) error {
	items, err := ListTodoItems(db)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get todo items")
	}
	logger.Println("Listing items")
	return ctx.Render("index", fiber.Map{
		"Items": items,
	})
}

func postHandler(ctx *fiber.Ctx, db *sql.DB) error {
	item := TodoItem{}
	if err := ctx.BodyParser(&item); err != nil {
		logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	id, err := item.Create(db)
	if err != nil {
		logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	logger.Printf("Item created: %+v\n", item)
	idStr := strconv.FormatInt(id, 10)
	return ctx.Status(http.StatusCreated).SendString(idStr)
}

func putHandler(ctx *fiber.Ctx, db *sql.DB) error {
	id := ctx.Params("id")
	item := FindItem(db, id) // pointer
	if item == nil {
		logger.Println(ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).SendString(ErrItemNotFound.Error())
	}

	if err := ctx.BodyParser(item); err != nil {
		logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if ctx.FormValue("Status") == "" {
		item.Status = false
	}
	err := item.Save(db)
	if err != nil {
		logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	logger.Printf("Item updated: %+v\n", item)
	return ctx.Status(http.StatusOK).JSON(item)
}

func deleteHandler(ctx *fiber.Ctx, db *sql.DB) error {
	id := ctx.Params("id")
	item := FindItem(db, id) // pointer
	if item == nil {
		logger.Println(ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).SendString(ErrItemNotFound.Error())
	}

	err := item.Delete(db)
	if err != nil {
		logger.Println(ErrCannotDeleteItem)
		return ctx.Status(http.StatusInternalServerError).SendString(ErrCannotDeleteItem.Error())
	}
	logger.Printf("Item deleted: %v\n", *item)
	return ctx.Status(http.StatusNoContent).SendString("")
}

func ping(ctx *fiber.Ctx) error {
	return ctx.SendString("PING")
}

func instance(ctx *fiber.Ctx) error {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	bs := make([]byte, resp.ContentLength)
	resp.Body.Read(bs)
	resp.Body.Close()
	return ctx.SendString(string(bs))
}
