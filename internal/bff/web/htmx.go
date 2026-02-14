package web

import "github.com/gofiber/fiber/v3"

func IsHTMX(c fiber.Ctx) bool {
	hxRequest := c.Get("HX-Request")
	return hxRequest == "true"
}
