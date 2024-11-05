package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ysfgrl/gerror"
)

func baseResponse(c *fiber.Ctx, status int, content any, err *gerror.Error) error {
	if err != nil {
		return c.Status(status).JSON(Response{
			Code:    status,
			Content: content,
			Error:   err,
		})
	} else {
		return c.Status(status).JSON(Response{
			Code:    status,
			Content: content,
			Error:   nil,
		})
	}

}

func OK(c *fiber.Ctx, content any) error {
	return baseResponse(c, fiber.StatusOK, content, nil)
}

func Created(c *fiber.Ctx, content any) error {
	return baseResponse(c, fiber.StatusCreated, content, nil)
}

func Unauthorized(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusUnauthorized, nil, err)
}
func Forbidden(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusForbidden, nil, err)
}
func NotAllowed(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusMethodNotAllowed, nil, err)
}

func NotImplemented(c *fiber.Ctx) error {
	return baseResponse(c, fiber.StatusNotImplemented, nil, gerror.UserError("NotImplemented", "NotImplemented"))
}

func NotFound(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusNotFound, nil, err)
}

func BadRequest(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusBadRequest, nil, err)
}

func InternalServerError(c *fiber.Ctx, err *gerror.Error) error {
	return baseResponse(c, fiber.StatusInternalServerError, nil, err)
}
