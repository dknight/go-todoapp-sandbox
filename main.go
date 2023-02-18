package main

// Hot reload
// go run github.com/cosmtrek/air

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	_ "github.com/lib/pq"
)

const (
	viewsDir  = "./views"
	publicDir = "./public"
)

var (
	dbConnString string
	port         string
	db           *sql.DB
)

func init() {
	var err error
	dbConnString = "postgresql://postgres:123456@localhost:5432/todoapp?sslmode=disable"
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(err)
	}
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
	log.Fatalln(app.Listen(":" + port))
}

func indexHandler(ctx *fiber.Ctx, db *sql.DB) error {
	items, err := ListTodoItems(db)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get todo items")
	}
	return ctx.Render("index", fiber.Map{
		"Items": items,
	})
}

func postHandler(ctx *fiber.Ctx, db *sql.DB) error {
	item := TodoItem{}
	if err := ctx.BodyParser(&item); err != nil {
		log.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	id, err := item.Create(db)
	if err != nil {
		log.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	idStr := strconv.FormatInt(id, 10)
	return ctx.Status(http.StatusCreated).SendString(idStr)
}

func putHandler(ctx *fiber.Ctx, db *sql.DB) error {
	id := ctx.Params("id")
	item := FindItem(db, id) // pointer
	if item == nil {
		log.Println(ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).SendString(ErrItemNotFound.Error())
	}

	if err := ctx.BodyParser(item); err != nil {
		log.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if ctx.FormValue("Status") == "" {
		item.Status = false
	}
	err := item.Save(db)
	if err != nil {
		log.Println(err)
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(item)
}

func deleteHandler(ctx *fiber.Ctx, db *sql.DB) error {
	id := ctx.Params("id")
	item := FindItem(db, id) // pointer
	if item == nil {
		log.Println(ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).SendString(ErrItemNotFound.Error())
	}

	err := item.Delete(db)
	if err != nil {
		log.Println(ErrCannotDeleteItem)
		return ctx.Status(http.StatusInternalServerError).SendString(ErrCannotDeleteItem.Error())
	}
	return ctx.Status(http.StatusNoContent).SendString("")
}
