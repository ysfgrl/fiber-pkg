package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
	"strconv"
)

type App struct {
	apps   []Interface
	config *Config
}

func NewApp(config *Config) *App {
	return &App{
		apps:   make([]Interface, 0),
		config: config,
	}
}
func (a *App) AddApp(routeInterface Interface) {
	a.apps = append(a.apps, routeInterface)
}

func (a *App) Listen(host string, port int) error {
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				switch e.Code {
				case 405:
					return response.NotAllowed(ctx, &gerror.Error{
						Code: "method.not_allowed",
					})
				}
			}
			return err
		},
	})
	app.Use(cors.New())
	for _, r := range a.apps {
		r.Routes(app)
	}
	return app.Listen(host + ":" + strconv.Itoa(port))
}
