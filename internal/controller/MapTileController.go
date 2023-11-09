package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
	"github.com/yockii/giserver-express/pkg/util"
	"strings"
	"sync"
)

var MapTileController = &mapTileController{
	mapRand: make(map[string]string),
}

type mapTileController struct {
	mapRand map[string]string
	lock    sync.Mutex
}

func (*mapTileController) Add(ctx *fiber.Ctx) error {
	vt := new(model.MapTile)
	if err := ctx.BodyParser(vt); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if vt.Name == "" || vt.StoreId == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "store id & name must be provided",
		})
	}
	d, err := service.MapTileService.Add(vt)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "sceneId & name duplicated",
		})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: vt.Id,
		})
	}
}
func (*mapTileController) Update(ctx *fiber.Ctx) error {
	data := new(model.MapTile)
	if err := ctx.BodyParser(data); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if data.Id == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "id must be provided",
		})
	}
	err := service.MapTileService.Update(data)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(&server.CommonResponse{})
}

func (*mapTileController) List(ctx *fiber.Ctx) error {
	condition := new(model.MapTile)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	limit, offset, orderBy, err := server.ParsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	total, list, err := service.MapTileService.List(condition, offset, limit, orderBy)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Items:  list,
		},
	})
}

func (c *mapTileController) GetFile(ctx *fiber.Ctx) error {
	df := ctx.Params("*")
	dfArray := strings.Split(df, "/")
	var dp []string
	for _, s := range dfArray {
		dp = append(dp, s)
	}
	vtName := ctx.Params("name")

	requestFileName := dp[len(dp)-1]
	ct := fiber.MIMEOctetStream
	if strings.HasSuffix(requestFileName, ".mvt") {
		ct = "application/mvt"
	} else if strings.HasSuffix(requestFileName, ".png") {
		ct = "image/png"
	} else if strings.HasSuffix(requestFileName, ".json") {
		ct = fiber.MIMEApplicationJSONCharsetUTF8
	} else {
		// 检查etag
		fileTagInfo := util.HashHex(c.getRand(vtName) + "|" + ctx.Path())
		if ctx.Get(fiber.HeaderIfNoneMatch) == fileTagInfo {
			return ctx.SendStatus(fiber.StatusNotModified)
		}
		ctx.Set(fiber.HeaderETag, fileTagInfo)
	}
	ctx.Set(fiber.HeaderContentType, ct)
	reader, err := service.MapTileService.ReadFile(vtName, dp...)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if reader == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	ctx.Set(fiber.HeaderCacheControl, "max-age=86400")
	return ctx.SendStream(reader)
}

func (c *mapTileController) DeleteCache(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	c.clearCache(name)
	service.MapTileService.ClearCache(name)
	return ctx.JSON(server.CommonResponse{Data: "OK"})
}

func (c *mapTileController) getRand(name string) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	if rand, ok := c.mapRand[name]; ok && rand != "" {
		return rand
	}

	rand := ksuid.New().String()
	c.mapRand[name] = rand
	return rand
}

func (c *mapTileController) clearCache(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.mapRand, name)
}
