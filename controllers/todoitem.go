package controllers

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dknight/go-todoapp-sandbox/models"
	"github.com/gofiber/fiber"
)

type TodoItemController struct {
	db *sql.DB
}

func (ctrl TodoItemController) Index(c *fiber.Ctx) error {
	items, err := models.ListTodoItems(ctrl.db)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get todo items")
	}
	logger.Println("Listing items")
	return ctx.Render("index", fiber.Map{
		"Items": items,
	})
}
