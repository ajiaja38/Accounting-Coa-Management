package middleware

import (
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func DBTransaction(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()
		if tx.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to begin transaction")
		}

		c.Locals("tx", tx)

		if err := c.Next(); err != nil {
			tx.Rollback()
			return err
		}

		if c.Response().StatusCode() >= 400 {
			tx.Rollback()
			return nil
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
		}

		return nil
	}
}
