package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
)

var StoreController = new(storeController)

type storeController struct{}

func (*storeController) Add(ctx *fiber.Ctx) error {
	store := new(model.Store)
	if err := ctx.BodyParser(store); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if store.Name == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "name must be provided",
		})
	}
	d, err := service.StoreService.Add(store)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: store.Id,
		})
	}
}
