package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/dknight/go-todoapp-sandbox/models"
	"github.com/gofiber/fiber/v2"
)

type TodoItemController struct {
	db     *sql.DB
	logger *log.Logger
}

func NewTodoController(db *sql.DB, logger *log.Logger) *TodoItemController {
	return &TodoItemController{db, logger}
}

func (ctrl TodoItemController) Index(ctx *fiber.Ctx) error {
	items, err := models.ListTodoItems(ctrl.db)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get todo items")
	}
	ctrl.logger.Println("Listing items")
	return ctx.Render("index", fiber.Map{
		"Items": items,
	})
}

func (ctrl TodoItemController) Post(ctx *fiber.Ctx) error {
	item := models.TodoItem{}
	if err := ctx.BodyParser(&item); err != nil {
		ctrl.logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	id, err := item.Create(ctrl.db)
	if err != nil {
		ctrl.logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.logger.Printf("Item created: %+v\n", item)
	idStr := strconv.FormatInt(id, 10)
	return ctx.Status(http.StatusCreated).SendString(idStr)
}

func (ctrl TodoItemController) Put(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	item := models.FindItem(ctrl.db, id) // pointer,
	if item == nil {
		ctrl.logger.Println(models.ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).
			SendString(models.ErrItemNotFound.Error())
	}

	if err := ctx.BodyParser(item); err != nil {
		ctrl.logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}

	if ctx.FormValue("Status") == "" {
		item.Status = false
	}
	err := item.Save(ctrl.db)
	if err != nil {
		ctrl.logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.logger.Printf("Item updated: %+v\n", item)
	return ctx.Status(http.StatusOK).JSON(item)
}

func (ctrl TodoItemController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	item := models.FindItem(ctrl.db, id) // pointer
	if item == nil {
		ctrl.logger.Println(models.ErrItemNotFound)
		return ctx.Status(http.StatusNotFound).
			SendString(models.ErrItemNotFound.Error())
	}

	err := item.Delete(ctrl.db)
	if err != nil {
		ctrl.logger.Println(models.ErrCannotDeleteItem)
		return ctx.Status(http.StatusInternalServerError).
			SendString(models.ErrCannotDeleteItem.Error())
	}
	ctrl.logger.Printf("Item deleted: %v\n", *item)
	return ctx.Status(http.StatusNoContent).SendString("")
}
