package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Message = func(status string, message string, data interface{}) Response {
	return Response{Status: status, Message: message, Data: data}
}

var Respond = func(c *fiber.Ctx, code int, response Response) error {
	c.Status(code)
	return c.JSON(response)
}
