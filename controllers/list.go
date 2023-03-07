package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dknight/go-todoapp-sandbox/lib"
	"github.com/dknight/go-todoapp-sandbox/models"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrListEmptyName = errors.New("List name cannot be empty")
)

type ListController struct {
	env *lib.Env
}

func NewListController(env *lib.Env) *ListController {
	return &ListController{env}
}

func (ctrl ListController) GetLists(ctx *fiber.Ctx) error {
	lists, err := models.ListLists(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(lists)
}

func (ctrl ListController) NewList(ctx *fiber.Ctx) error {
	lists, err := models.ListLists(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	return ctx.Render("list", fiber.Map{
		"Lists": lists,
	})
}

func (ctrl ListController) Post(ctx *fiber.Ctx) error {
	list := models.List{}
	if err := ctx.BodyParser(&list); err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	id, err := list.Create(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.env.Logger.Printf("Item created: %+v\n", list)
	idStr := strconv.FormatInt(id, 10)
	return ctx.Status(http.StatusCreated).SendString(idStr)
}

func (ctrl ListController) Put(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	list, err := models.FindList(ctrl.env.DB, id)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusNotFound).
			SendString(err.Error())
	}

	if err := ctx.BodyParser(list); err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}

	if ctx.FormValue("Name") == "" {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(ErrListEmptyName.Error())
	}
	err = list.Save(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}
	ctrl.env.Logger.Printf("List updated: %+v\n", list)
	return ctx.Status(http.StatusOK).JSON(list)
}

func (ctrl ListController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	list, err := models.FindList(ctrl.env.DB, id)
	if err != nil {
		ctrl.env.Logger.Println(err)
		return ctx.Status(http.StatusBadRequest).
			SendString(err.Error())
	}

	err = list.Delete(ctrl.env.DB)
	if err != nil {
		ctrl.env.Logger.Println(models.ErrCannotDeleteItem)
		return ctx.Status(http.StatusInternalServerError).
			SendString(models.ErrCannotDeleteItem.Error())
	}
	ctrl.env.Logger.Printf("Item deleted: %v\n", *list)
	return ctx.Status(http.StatusNoContent).SendString("")
}
