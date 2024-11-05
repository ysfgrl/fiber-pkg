package auth

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
)

type Auth struct {
	cfg Config
}

func NewAuth(config Config) *Auth {
	return &Auth{
		cfg: config,
	}
}

func (a *Auth) Required() fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup: a.cfg.TokenLookup,
		AuthScheme:  a.cfg.AuthScheme,
		ContextKey:  a.cfg.ContextKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.Unauthorized(c, &gerror.Error{
				Code:   "api.permission_denied",
				Detail: err.Error(),
			})
		},
		SuccessHandler: a.cfg.SuccessHandler,
		Claims:         a.cfg.Claims,
		KeyFunc:        a.cfg.KeyFunc,
	})
}
func (a *Auth) RoleRequired(roles ...string) fiber.Handler {
	if a.cfg.RoleRequire == nil {
		return func(ctx *fiber.Ctx) error {
			return ctx.Next()
		}
	}
	return a.cfg.RoleRequire(roles...)
}

func (a *Auth) Guest() *jwt.Token {
	return a.cfg.Guest
}
