package controller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/domain"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
)

var SpaceController = new(spaceController)

type spaceController struct{}

func (*spaceController) Add(ctx *fiber.Ctx) error {
	space := new(model.Space)
	if err := ctx.BodyParser(space); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if space.Name == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "name must be provided",
		})
	}
	d, err := service.SpaceService.Add(space)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: space.Id,
		})
	}
}

func (c *spaceController) SpaceInfo(ctx *fiber.Ctx) error {
	spaceName := ctx.Params("spaceName")
	space, err := service.SpaceService.FindByName(spaceName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if space == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var scenes []*model.Scene
	// 获取只需要 resourceconfigid、supportmediatypes、name、resourcetype几个字段
	scenes, err = service.SceneService.FindSpaceScenes(space.Id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	var sd []*domain.SpaceScene
	hostName := ctx.Hostname()
	protocol := ctx.Protocol()
	for _, scene := range scenes {
		s := scene
		sceneLocation, _ := ctx.GetRouteURL("scene.info", fiber.Map{"sceneId": s.Id})
		sDomain := &domain.SpaceScene{
			ResourceConfigId:    s.ResourceConfigId,
			SupportedMediaTypes: strings.Split(s.SupportedMediaTypes, ","),
			Path:                protocol + "://" + hostName + sceneLocation,
			Name:                s.Name,
			ResourceType:        s.ResourceType,
		}
		sd = append(sd, sDomain)
	}

	return ctx.JSON(sd)
}