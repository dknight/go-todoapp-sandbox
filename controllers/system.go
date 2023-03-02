package controllers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (ctrl SystemController) Ping(ctx *fiber.Ctx) error {
	return ctx.SendString("PING")
}

func (ctrl SystemController) Instance(ctx *fiber.Ctx) error {
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
