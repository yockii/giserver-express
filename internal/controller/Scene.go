package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/domain"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
)

var SceneController = new(sceneController)

type sceneController struct{}

func (*sceneController) Add(ctx *fiber.Ctx) error {
	scene := new(model.Scene)
	if err := ctx.BodyParser(scene); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if scene.SpaceId == 0 || scene.Name == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "id & name must be provided",
		})
	}
	d, err := service.SceneService.Add(scene)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: scene.Id,
		})
	}
}

func (c *sceneController) SceneInfo(ctx *fiber.Ctx) error {
	sceneIdStr := ctx.Params("sceneId")
	sceneId, err := strconv.ParseInt(sceneIdStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	var scene *domain.Scene
	scene, err = service.SceneService.GetRichSceneInfoById(sceneId)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if scene == nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"succeed": false,
			"error": fiber.Map{
				"code":     404,
				"errorMsg": "资源不存在",
			},
		})
	}
	return ctx.JSON(scene)
}

func (c *sceneController) SceneLayers(ctx *fiber.Ctx) error {
	sceneIdStr := ctx.Params("sceneId")
	sceneId, err := strconv.ParseInt(sceneIdStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	var layers []*domain.SceneLayer
	layers, err = service.LayerService.FindSceneDomainLayers(sceneId)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(layers)
}
