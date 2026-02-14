package web

import "github.com/gofiber/fiber/v3"

type HomeHandlers struct {
	r *Renderer
}

func NewHomeHandlers(r *Renderer) *HomeHandlers {
	return &HomeHandlers{r: r}
}

func (h *HomeHandlers) Home(c fiber.Ctx) error {
	return h.r.Render(c, "layout", nil)
}
