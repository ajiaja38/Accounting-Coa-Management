package utils

import (
	"fiber.com/session-api/pkg/model"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse[T any](c *fiber.Ctx, code int, message string, data T) error {
	return c.Status(code).JSON(model.ResponseEntity[T]{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SuccessResponsePaginate[T any](c *fiber.Ctx, code int, message string, data T, meta *model.MetaPagination) error {
	return c.Status(code).JSON(model.ResponseEntityPagination[T]{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
