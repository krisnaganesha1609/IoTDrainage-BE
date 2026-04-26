package utils

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func RespondWithOK(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func RespondWithCreated(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status:  "created",
		Message: message,
		Data:    data,
	})
}

func RespondWithError(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "error",
		Message: message,
	})
}

func RespondWithValidationError(c fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Status:  "error",
		Message: "Validation error",
		Error:   errors,
	})
}
