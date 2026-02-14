package web

import (
	"html/template"
	"io"
	"github.com/gofiber/fiber/v3"
)

type Renderer struct {
	t *template.Template
}

func NewRenderer(pattern string) (*Renderer, error) {
	t, err := template.ParseGlob(pattern)
	if err != nil {
	    return nil, err
	}
	return &Renderer{t: t}, nil
}

func (r *Renderer) Render(c fiber.Ctx, name string, data any) error {
	c.Type("html", "utf-8")
	return r.t.ExecuteTemplate(c.Response().BodyWriter(), name, data)
}

func (r *Renderer) RenderTo(w io.Writer, name string, data any) error {
	return r.t.ExecuteTemplate(w, name, data)
}
