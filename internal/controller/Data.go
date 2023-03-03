package controller

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
	"github.com/yockii/giserver-express/pkg/util"
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

func (*dataController) Update(ctx *fiber.Ctx) error {
	data := new(model.Data)
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
	err := service.DataService.Update(data)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(&server.CommonResponse{})
}

func (*dataController) List(ctx *fiber.Ctx) error {
	condition := new(model.Data)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	limit, offset, orderBy, err := server.ParsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	total, list, err := service.DataService.List(condition, offset, limit, orderBy)
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

func (c *dataController) DataFile(ctx *fiber.Ctx) error {
	spaceName := ctx.Params("spaceName")
	dataName := ctx.Params("dataName")

	suffix := strings.ToUpper(dataName[strings.LastIndex(dataName, "."):])
	var data *model.Data
	data = service.DataService.GetFromCache(spaceName, suffix, dataName)
	if data == nil {
		space, err := service.SpaceService.FindByName(spaceName)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if space == nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		data, err = service.DataService.GetBySpaceIdAndDataName(space.Id, suffix, dataName)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if data == nil {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		service.DataService.Cache(spaceName, data.DataType, data)
	}

	if data.DataStoreTypeId == 0 {
		return c.sendFromLocalFile(ctx, data, data.DataName)
	} else {
		store, err := service.StoreService.GetById(data.DataStoreTypeId)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if store != nil {
			if store.StoreType == 0 {
				return c.sendFromLocalFile(ctx, data, suffix)
			} else if store.StoreType == 1 {
				objectKey := store.Path + data.DataConfigPath + "/" + data.DataName
				fileInfo := util.HashHex(objectKey)
				if ctx.Get(fiber.HeaderIfNoneMatch) == fileInfo {
					return ctx.SendStatus(fiber.StatusNotModified)
				}
				reader, err := service.OssService.StreamFromStore(store, objectKey)
				if err != nil {
					return ctx.SendStatus(fiber.StatusInternalServerError)
				}
				ctx.Set(fiber.HeaderETag, fileInfo)
				return ctx.SendStream(reader)
			}
		}
	}
	return ctx.SendStatus(fiber.StatusNotFound)
}

func (c *dataController) sendFromLocalFile(ctx *fiber.Ctx, data *model.Data, file string) error {
	f, err := os.Open(path.Join(data.DataConfigPath, file))
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if fs, _ := f.Stat(); fs != nil {
		fileInfo := strconv.FormatInt(fs.ModTime().Unix(), 16)
		if ctx.Get(fiber.HeaderIfNoneMatch) == fileInfo {
			f.Close()
			return ctx.SendStatus(fiber.StatusNotModified)
		}
		ctx.Set(fiber.HeaderETag, fileInfo)
	}
	return ctx.SendStream(f)
}
