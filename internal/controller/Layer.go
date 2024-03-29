package controller

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/constant"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/server"
	"github.com/yockii/giserver-express/pkg/util"
)

var LayerController = new(layerController)

type layerController struct{}

func (*layerController) Add(ctx *fiber.Ctx) error {
	layer := new(model.SceneLayer)
	if err := ctx.BodyParser(layer); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if layer.SpaceId == 0 || layer.SceneId == 0 || layer.DataId == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "spaceId/sceneId/dataId must be provided",
		})
	}
	d, err := service.LayerService.Add(layer)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if d {
		return ctx.JSON(&server.CommonResponse{
			Code: -1,
			Msg:  "sceneId & dataId duplicated",
		})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: layer.Id,
		})
	}
}
func (*layerController) Update(ctx *fiber.Ctx) error {
	data := new(model.SceneLayer)
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
	err := service.LayerService.Update(data)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(&server.CommonResponse{})
}
func (*layerController) List(ctx *fiber.Ctx) error {
	condition := new(model.SceneLayer)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	limit, offset, orderBy, err := server.ParsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	total, list, err := service.LayerService.List(condition, offset, limit, orderBy)
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

func (c *layerController) GetLayerExtendXml(ctx *fiber.Ctx) error {
	sceneIdStr := ctx.Params("sceneId")
	sceneIdInt64, err := strconv.ParseInt(sceneIdStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	sceneId := database.Int64(sceneIdInt64)
	layerName := ctx.Params("layerName")
	if sceneId == 0 || layerName == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var layer *model.SceneLayer
	layer, err = service.LayerService.GetBySceneIdAndLayerName(sceneId, layerName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	var data *model.Data
	data, err = service.DataService.GetById(layer.DataId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if data == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Response().Header.SetContentType(fiber.MIMEApplicationXML)
	return ctx.SendString(c.generateExtendXml(data.DataName, data.DataConfigPath, layer.CreateTime.String(), layer.Name, data.DataType, "NONE"))
}

// generateExtendXml cacheFileType:OSGB/S3M   renderCullMode: NONE/DEFAULT
func (*layerController) generateExtendXml(caption, dsAlias, createTime, layerName, cacheFileType, renderCullMode string) string {
	return fmt.Sprintf(constant.ExtendXmlTmpl, caption, dsAlias, createTime, layerName, cacheFileType, renderCullMode)
}

func (c *layerController) LayerConfig(ctx *fiber.Ctx) error {
	spaceName := ctx.Params("spaceName")
	layerName := ctx.Params("layerName")
	space, err := service.SpaceService.FindByName(spaceName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	var layer *model.SceneLayer
	layer, err = service.LayerService.GetBySpaceIdAndLayerName(space.Id, layerName)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	var data *model.Data
	data, err = service.DataService.GetById(layer.DataId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// 找到config路径，返回
	if data.DataStoreTypeId == -1 {
		return ctx.SendFile(path.Join(data.DataConfigPath, data.DataName))
	} else {
		store, err := service.StoreService.GetById(data.DataStoreTypeId)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if store != nil {
			reader, err := service.OssService.StreamFromStore(store, store.Path+data.DataConfigPath+"/"+data.DataName)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			return ctx.SendStream(reader)
		}
	}
	return ctx.SendStatus(fiber.StatusNotFound)
}

func (c *layerController) TileFile(ctx *fiber.Ctx) error {
	spaceName := ctx.Params("spaceName")
	dataName := ctx.Params("dataName")
	var data *model.Data
	data = service.DataService.GetFromCache(spaceName, "OSGB", dataName)
	if data == nil {
		data = service.DataService.GetFromCache(spaceName, "S3M", dataName)
		if data == nil {
			space, err := service.SpaceService.FindByName(spaceName)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			if space == nil {
				return ctx.SendStatus(fiber.StatusNotFound)
			}
			data, err = service.DataService.GetBySpaceIdAndDataName(space.Id, "", dataName)
			if err != nil {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}
			if data == nil {
				return ctx.SendStatus(fiber.StatusNotFound)
			}
			service.DataService.Cache(spaceName, data.DataType, data)
		}
	}

	fold := ctx.Params("fold")
	file := ctx.Params("file")

	ctx.Set(fiber.HeaderConnection, "keep-alive")
	ctx.Set(fiber.HeaderKeepAlive, "timeout=12")

	if strings.HasSuffix(file, ".s3mb") {
		ctx.Set(fiber.HeaderContentType, "application/s3mb")
	} else if strings.HasSuffix(file, ".s3m") {
		ctx.Set(fiber.HeaderContentType, "application/s3m")
	}

	if data.DataStoreTypeId == -1 {
		return c.sendFromLocalFile(ctx, data, fold, file)
	} else {
		store, err := service.StoreService.GetById(data.DataStoreTypeId)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if store != nil {
			if store.StoreType == -1 {
				return c.sendFromLocalFile(ctx, data, fold, file)
			} else if store.StoreType == 1 {
				objectKey := store.Path + data.DataConfigPath + "/" + fold + "/" + file
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

func (c *layerController) sendFromLocalFile(ctx *fiber.Ctx, data *model.Data, fold string, file string) error {
	f, err := os.Open(path.Join(data.DataConfigPath, fold, file))
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
