package auth

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	jwtware.Config
	RoleRequire func(roles ...string) fiber.Handler
	Guest       *jwt.Token
}
