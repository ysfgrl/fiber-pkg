package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ysfgrl/fiber-pkg/auth"
	"github.com/ysfgrl/fiber-pkg/logger"
	"github.com/ysfgrl/gerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Base struct {
	Auth    *auth.Auth
	Logger  *logger.Logger
	Prefix  string
	handler []*Route
}

func (ctrl *Base) AddRoute(route *Route) {
	ctrl.handler = append(ctrl.handler, route)
}

func (ctrl *Base) Routes(app *fiber.App) {
	group := app.Group(ctrl.Prefix)
	if ctrl.Auth != nil {
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user", ctrl.Auth.Guest())
			return c.Next()
		})
	}
	fmt.Printf("--------------------------\n")
	fmt.Printf("%s \n", ctrl.Prefix)
	for _, route := range ctrl.handler {
		var handlers []fiber.Handler
		if ctrl.Logger != nil && route.Log {
			handlers = append(handlers, ctrl.Logger.Log)
		}
		if ctrl.Auth != nil && route.Auth {
			handlers = append(handlers, ctrl.Auth.Required())
			handlers = append(handlers, ctrl.Auth.RoleRequired(route.Role...))
		}
		handlers = append(handlers, route.Handler)
		group.Add(route.Method, route.Path, handlers...)
		fmt.Printf("%6s %s\n", route.Method, route.Path)
	}
}
func (ctrl *Base) GetUser(ctx *fiber.Ctx) jwt.Claims {
	contextUser := ctx.Locals("user").(*jwt.Token)
	return contextUser.Claims
}

func (ctrl *Base) GetToken(ctx *fiber.Ctx) string {
	contextUser := ctx.Locals("user").(*jwt.Token)
	return contextUser.Raw
}

func (ctrl *Base) GetIdParams(c *fiber.Ctx, key string) (primitive.ObjectID, *gerror.Error) {
	value := c.Params(key, key)
	if strings.EqualFold(value, key) {
		return primitive.ObjectID{}, &gerror.Error{
			Code: "api.id_nf",
		}
	}
	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return primitive.ObjectID{}, &gerror.Error{
			Code: "api.id_nf",
		}
	}
	return id, nil
}
func (ctrl *Base) GetParams(c *fiber.Ctx, key string) (string, *gerror.Error) {
	value := c.Params(key, key)
	if strings.EqualFold(value, key) {
		return key, gerror.UserError(key+".required", "required")
	}
	return value, nil
}
