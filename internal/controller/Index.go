package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yockii/giserver-express/pkg/server"
	"github.com/yockii/giserver-express/pkg/util"
)

func InitRouter() {
	{
		manage := server.Group("/manage")
		store := manage.Group("/store")
		store.Post("/", StoreController.Add)
		store.Put("/", StoreController.Update)
		store.Get("/list", StoreController.List)

		space := manage.Group("/space")
		space.Post("/", SpaceController.Add)
		space.Put("/", SpaceController.Update)
		space.Get("/list", SpaceController.List)

		scene := manage.Group("/scene")
		scene.Post("/", SceneController.Add)
		scene.Put("/", SceneController.Update)
		scene.Get("/list", SceneController.List)

		data := manage.Group("/data")
		data.Post("/", DataController.Add)
		data.Put("/", DataController.Update)
		data.Get("/list", DataController.List)

		layer := manage.Group("/layer")
		layer.Post("/", LayerController.Add)
		layer.Put("/", LayerController.Update)
		layer.Get("/list", LayerController.List)
	}

	space := server.Group("/giservices/space")
	space.Get("/:spaceName", SpaceController.SpaceInfo)

	scene := server.Group("/giservices/scene")
	scene.Get("/:sceneId.json", SceneController.SceneInfo)
	scene.Get("/:sceneId/layers.json", SceneController.SceneLayers)
	scene.Get("/:sceneId", SceneController.SceneInfo).Name("scene.info")
	scene.Get("/:sceneId/layers/:layerName/extendxml.xml", LayerController.GetLayerExtendXml)

	// iserver适配
	{
		server.Get("/giservices/:spaceName/rest/realspace/scenes.json", SpaceController.SpaceInfo)
		server.Get("/giservices/:spaceName/rest/realspace/datas/:dataName/config", LayerController.LayerConfig)
		server.Get("/giservices/:spaceName/rest/realspace/datas/:dataName/data/path/:fold/:file", LayerController.TileFile)

		server.Get("/giservices/:spaceName/rest/realspace/login.json", func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"random":     util.SnowflakeIdString(),
				"jsessionID": util.SnowflakeIdString(),
			})
		})
		server.Post("/giservices/:spaceName/rest/realspace/login.json", func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"postResultType": "CreateChild",
				"succeed":        true,
			})
		})
		server.Get("/giservices/:spaceName/rest/realspace/_setup.json", func(ctx *fiber.Ctx) error {
			ctx.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
			return ctx.SendString(`{"isCloudLicenseLogin":false,"serviceLanguage":"chinese","isUGODllError":false,"licenseMode":"DefaultLicense","iserverFeaturesPackageType":"ALL","isEduLicense":false,"isLicenseFinished":true,"isServiceNode":false,"isExpress":false,"cloudLicenseSetting":null,"computerName":"iZ25gt8gjcwZ","cloudLicenseValid":false,"isLicenseError":false,"licErrorMsg":"","iserverLicenseInfo":{"cloudLicenseSetting":null,"isCloudLicenseLogin":false,"entryInfos":[{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21002","licenseID":21002,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21007","licenseID":21007,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21004","licenseID":21004,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21006","licenseID":21006,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21003","licenseID":21003,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21005","licenseID":21005,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21010","licenseID":21010,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21009","licenseID":21009,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21011","licenseID":21011,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21014","licenseID":21014,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21012","licenseID":21012,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21013","licenseID":21013,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0},{"expireDateTime":"2026-11-12","useWith":-1,"licenseStatus":true,"licenseModel":null,"expireDate":{"date":12,"hours":23,"seconds":59,"month":10,"year":126,"minutes":59,"time":1794499199704},"hardwareKeyType":null,"licenseName":"21015","licenseID":21015,"isTrial":false,"userTrademark":"SuperMap","watermarkMode":0}],"productVersion":null,"iserverVersion":"SuperMap iServer高级版","summaryInfo":null,"companyName":"SuperMap","isEduLicense":false,"masterServerAddress":null,"user":"SuperMap","licenseServer":null,"isSuperMapStaff":false},"isSuperMapStaff":false,"serviceLanguages":[],"eduLicenseSetting":null,"isExtendModule":false,"javaActualVersion":"1.8.0_332","iserverUGOVersion":"11.0.1.21420","webLicenseValid":true,"eduLicenseValid":true,"isSupportHardwareLicMode":true,"isPortal":false,"stepParam":null,"systemUGOVersion":"11.0.1.21420","isAdminExist":true,"coreLicExtendExist":false,"isAix":false,"isJDKVersionError":false,"isUGOVersionError":false,"javaExpectedVersion":"1.8","coreLicExistButUnavailable":false}`)
		})

		// KML等文件
		server.Get("/giservices/:spaceName/data/:dataName", DataController.DataFile)
	}
}
