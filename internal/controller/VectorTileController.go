package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/internal/service"
	"github.com/yockii/giserver-express/pkg/server"
	"io"
	"strings"
)

var VectorTileController = new(vectorTileController)

type vectorTileController struct{}

func (*vectorTileController) Add(ctx *fiber.Ctx) error {
	vt := new(model.VectorTile)
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
	d, err := service.VectorTileService.Add(vt)
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
func (*vectorTileController) Update(ctx *fiber.Ctx) error {
	data := new(model.VectorTile)
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
	err := service.VectorTileService.Update(data)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.JSON(&server.CommonResponse{})
}

func (*vectorTileController) List(ctx *fiber.Ctx) error {
	condition := new(model.VectorTile)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	limit, offset, orderBy, err := server.ParsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	total, list, err := service.VectorTileService.List(condition, offset, limit, orderBy)
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

func (c *vectorTileController) VectorTileInfo(ctx *fiber.Ctx) error {
	vtName := ctx.Params("name")

	if strings.HasSuffix(vtName, ".json") {
		vtName = strings.TrimSuffix(vtName, ".json")
	}

	// 获取矢量瓦片索引数据
	indexJson, err := service.VectorTileService.GetIndex(vtName)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if indexJson == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	//oj := gjson.Parse(indexJson)
	//if !oj.Get("distanceUnit").Exists() {
	//	indexJson, _ = sjson.SetRaw(indexJson, "distanceUnit", "null")
	//}
	//if !oj.Get("tileversion").Exists() {
	//	indexJson, _ = sjson.SetRaw(indexJson, "tileversion", "null")
	//}
	//
	//indexJson, _ = sjson.SetRaw(indexJson, "layers.0.subLayers", "{}")
	//indexJson, _ = sjson.SetRaw(indexJson, "overlapDisplayedOptions", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "clipRegion", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "colorMode", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "customPrjCoordSysType", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "customEntireBounds", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "backgroundStyle", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "description", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "rasterfunction", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "layers.0.caption", "null")
	//indexJson, _ = sjson.SetRaw(indexJson, "layers.0.description", "null")
	//
	//// 科学计数法？
	//indexJson, _ = sjson.SetRaw(indexJson, "scale", "6.9229099892844785E-6")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.0", "6.9229099892844785E-6")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.1", "1.3845819978568957E-5")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.2", "2.7691639957137914E-5")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.3", "5.538327991427583E-5")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.4", "1.1076655982855166E-4")
	//indexJson, _ = sjson.SetRaw(indexJson, "visibleScales.5", "2.215331196571033E-4")

	return ctx.SendString(indexJson)
}

func (c *vectorTileController) GetSpriteJson(ctx *fiber.Ctx) error {
	vtName := ctx.Params("name")

	reader, err := service.VectorTileService.ReadMvtFile(vtName, "", "sprites", "sprite.json")
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if reader == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.SendStream(reader)
}

func (c *vectorTileController) GetStyleJson(ctx *fiber.Ctx) error {
	vtName := ctx.Params("name")

	reader, err := service.VectorTileService.ReadMvtFile(vtName, "styles", "style.json")
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if reader == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// 修改json内容，将路径改为绝对路径
	hostName := ctx.Hostname()
	protocol := ctx.Protocol()

	var styleJson []byte
	styleJson, err = io.ReadAll(reader)

	sj := gjson.ParseBytes(styleJson)

	sources := sj.Get("sources")
	for key, _ := range sources.Map() {
		styleJson, _ = sjson.SetBytes(styleJson, "sources."+key+".tiles", []string{fmt.Sprintf("%s://%s/giservices/vectortile/maps/%s/tiles/{z}/{x}/{y}.mvt", protocol, hostName, vtName)})
	}
	styleJson, _ = sjson.SetBytes(styleJson, "sprite", fmt.Sprintf("%s://%s/giservices/vectortile/maps/%s/sprites/sprite", protocol, hostName, vtName))
	styleJson, _ = sjson.SetBytes(styleJson, "glyphs", fmt.Sprintf("%s://%s/giservices/vectortile/maps/%s/fonts/{fontstack}/{range}.pbf", protocol, hostName, vtName))

	return ctx.SendString(string(styleJson))
}

func (c *vectorTileController) GetMvtFile(ctx *fiber.Ctx) error {
	df := ctx.Params("*")
	dfArray := strings.Split(df, "/")
	var dp []string
	for _, s := range dfArray {
		dp = append(dp, s)
	}

	vtName := ctx.Params("name")

	requestFileName := dp[len(dp)-1]
	ct := fiber.MIMEApplicationJSONCharsetUTF8
	if strings.HasSuffix(requestFileName, ".mvt") {
		ct = "application/mvt"
	} else if strings.HasSuffix(requestFileName, ".png") {
		ct = "image/png"
	}
	ctx.Set(fiber.HeaderContentType, ct)

	reader, err := service.VectorTileService.ReadMvtFile(vtName, dp...)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if reader == nil {
		if strings.HasSuffix(requestFileName, ".mvt") {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	ctx.Set(fiber.HeaderCacheControl, "max-age=10800")
	return ctx.SendStream(reader)
}

func (c *vectorTileController) DeleteCache(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	service.VectorTileService.ClearCache(name)
	return ctx.JSON(server.CommonResponse{Data: "OK"})
}
