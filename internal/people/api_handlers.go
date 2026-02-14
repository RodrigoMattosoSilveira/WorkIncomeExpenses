package people

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type API struct{ svc *Service }

func NewAPI(svc *Service) *API { return &API{svc: svc} }

type errorResponse struct {
	Errors map[string]string `json:"errors"`
}

type createReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type updateReq = createReq

func (a *API) Register(r fiber.Router) {
	r.Get("/api/people", a.list)
	r.Post("/api/people", a.create)
	r.Get("/api/people/:id", a.get)
	r.Patch("/api/people/:id", a.update)
	r.Delete("/api/people/:id", a.delete)
}

func (a *API) list(c fiber.Ctx) error {
	out, err := a.svc.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	return c.JSON(out)
}

func (a *API) get(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	p, err := a.svc.Get(c.Context(), uint(id))
	if err == ErrNotFound {
		return c.SendStatus(404)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	return c.JSON(p)
}

func (a *API) create(c fiber.Ctx) error {
	var req createReq
	// if err := c.BodyParser(&req); err != nil { https://github.com/gofiber/fiber/issues/2964
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}
	p, verrs, err := a.svc.Create(c.Context(), req.Name, req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	if len(verrs) > 0 {
		return c.Status(422).JSON(errorResponse{Errors: verrs})
	}
	return c.Status(201).JSON(p)
}

func (a *API) update(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req updateReq
	// if err := c.BodyParser(&req); err != nil { https://github.com/gofiber/fiber/issues/2964
	if err := c.Bind().Body(&req); err != nil {
		return c.SendStatus(400)
	}
	p, verrs, err := a.svc.Update(c.Context(), uint(id), req.Name, req.Email)
	if err == ErrNotFound {
		return c.SendStatus(404)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	if len(verrs) > 0 {
		return c.Status(422).JSON(errorResponse{Errors: verrs})
	}
	return c.JSON(p)
}

func (a *API) delete(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := a.svc.Delete(c.Context(), uint(id)); err == ErrNotFound {
		return c.SendStatus(404)
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	return c.SendStatus(204)
}
