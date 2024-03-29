package server

import (
	"github.com/gofiber/fiber/v2"
	"net"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/pkg/config"
	"github.com/yockii/giserver-express/pkg/database"
)

type webApp struct {
	app *fiber.App
}

var defaultApp *webApp

func init() {
	initServerDefault()

	initFiberParser()

	defaultApp = InitWebApp(nil)
}

func initServerDefault() {
	config.DefaultInstance.SetDefault("server.port", 13579)
}

func initFiberParser() {
	customDateTime := fiber.ParserType{
		Customtype: database.DateTime{},
		Converter:  database.DateTimeConverter,
	}
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        []fiber.ParserType{customDateTime},
		ZeroEmpty:         true,
	})
}

func InitWebApp(views fiber.Views) *webApp {
	initFiberParser()
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 views,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			logger.Error(e)
		},
	}))
	app.Use(cors.New())

	return &webApp{app}
}

func (a *webApp) Listener(ln net.Listener) error {
	return a.app.Listener(ln)
}
func (a *webApp) Static(dir string) {
	a.app.Static("/", dir, fiber.Static{
		Compress: true,
	})
}

func (a *webApp) Use(args ...interface{}) fiber.Router {
	return a.app.Use(args...)
}
func (a *webApp) Group(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Group(path, handlers...)
}
func (a *webApp) All(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.All(path, handlers...)
}
func (a *webApp) Get(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Get(path, handlers...)
}
func (a *webApp) Put(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Put(path, handlers...)
}
func (a *webApp) Post(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Post(path, handlers...)
}
func (a *webApp) Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return a.app.Delete(path, handlers...)
}
func (a *webApp) Start(addr string) error {
	return a.app.Listen(addr)
}
func (a *webApp) Shutdown() error {
	return a.app.Shutdown()
}

func Listener(ln net.Listener) error {
	return defaultApp.Listener(ln)
}
func Static(dir string) {
	defaultApp.Static(dir)
}

func Use(args ...interface{}) fiber.Router {
	return defaultApp.Use(args...)
}
func Group(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Group(path, handlers...)
}
func All(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.All(path, handlers...)
}
func Get(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Get(path, handlers...)
}
func Put(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Put(path, handlers...)
}
func Post(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Post(path, handlers...)
}
func Delete(path string, handlers ...fiber.Handler) fiber.Router {
	return defaultApp.Delete(path, handlers...)
}
func Start() error {
	return defaultApp.Start(":" + config.GetString("server.port"))
}
func Shutdown() error {
	return defaultApp.Shutdown()
}
