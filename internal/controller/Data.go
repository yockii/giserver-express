package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
)

var DataController = new(dataController)

type dataController struct{}

func (*dataController) Add(ctx *fiber.Ctx) error {
	data := new(model.Data)
	if err := ctx.BodyParser(data); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if data.SpaceId == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "spaceId must be provided",
		})
	}
	d, err := service.DataService.Add(data)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: data.Id,
		})
	}
}
