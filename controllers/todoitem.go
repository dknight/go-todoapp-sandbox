package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/dknight/go-todoapp-sandbox/lib"
	"github.com/dknight/go-todoapp-sandbox/models"
	"github.com/gofiber/fiber/v2"
)

type TodoItemController struct {
	env *lib.Env
}

func NewTodoController(env *lib.Env) *TodoItemController {
	return &TodoItemController{env}
}

func (ctrl TodoItemController) Index(ctx *fiber.Ctx) error {
	lists, err := models.ListLists(ctrl.env.DB)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get lists")
	}

	return ctx.Render("index", fiber.Map{
		"Lists": lists,
	})
}

func (ctrl TodoItemController) GetItemsByList(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("listid"))
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	items, err := models.FindItemsByListID(ctrl.env.DB, id)
	if err != nil {
		log.Println(err)
		return errors.New("Error: cannot get lists")
	}

	return ctx.Status(http.StatusOK).JSON(&items)
}

func (ctrl TodoItemController) Post(ctx *fiber.Ctx) error {
	item := models.TodoItem{}
	if err := ctx.BodyParser(&item); err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	id, err := item.Create(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.env.Logger.Printf("Item created: %+v\n", item)
	idStr := strconv.FormatInt(id, 10)
	return ctx.Status(http.StatusCreated).SendString(idStr)
}

func (ctrl TodoItemController) Put(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	item, err := models.FindItem(ctrl.env.DB, id)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}

	if err := ctx.BodyParser(item); err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}

	if ctx.FormValue("Status") == "" {
		item.Status = false
	}
	err = item.Save(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.env.Logger.Printf("Item updated: %+v\n", item)
	return ctx.Status(http.StatusOK).JSON(item)
}

func (ctrl TodoItemController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	item, err := models.FindItem(ctrl.env.DB, id)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}

	err = item.Delete(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(models.ErrCannotDeleteItem)
		return ctx.Status(http.StatusInternalServerError).
			SendString(models.ErrCannotDeleteItem.Error())
	}
	ctrl.env.Logger.Printf("Item deleted: %v\n", *item)
	return ctx.Status(http.StatusNoContent).SendString("")
}
