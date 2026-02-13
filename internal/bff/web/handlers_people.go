package web

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/RodrigoMattosoSilveira/WorkEarningsExpenses/internal/bff/clients"
)

type PeopleHandlers struct {
	r  *Renderer
	pc clients.PeopleAPI
}

func NewPeopleHandlers(r *Renderer, pc clients.PeopleAPI) *PeopleHandlers {
	return &PeopleHandlers{
		r: r, 
		pc: pc,
	}
}

type PeopleIndexVM struct {
	People []clients.PersonDTO
	Flash  string
}

type PersonFormVM struct {
	Person clients.PersonDTO
	Errors map[string]string
}

type PersonEditRowVM struct {
	Person clients.PersonDTO
	Errors map[string]string
}

func (h *PeopleHandlers) ListPeople(c fiber.Ctx) error {
	people, err := h.pc.ListPeople(c.Context())
	if err != nil {
		return c.Status(502).SendString("people service unavailable")
	}

	if IsHTMX(c) {
		return h.r.Render(c, "people_tbody", fiber.Map{"People": people})
	}

	// Full page
	return h.r.Render(c, "layout", PeopleIndexVM{People: people})
}

func (h *PeopleHandlers) NewPersonForm(c fiber.Ctx) error {
	vm := PersonFormVM{Errors: map[string]string{}}
	return h.r.Render(c, "people_form", vm) // typically loaded into #panel
}

func (h *PeopleHandlers) CreatePerson(c fiber.Ctx) error {
	req := clients.CreatePersonRequest{
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
	}

	person, vErrs, err := h.pc.CreatePerson(c.Context(), req)
	if err != nil {
		return c.Status(502).SendString("people service unavailable")
	}
	if len(vErrs) > 0 {
		// Replace the panel with the form showing errors, NOT the table.
		// If you want to keep hx-target="#people-tbody", then return a <tr> error row instead.
		c.Set("HX-Retarget", "#panel")
		c.Set("HX-Reswap", "innerHTML")
		return h.r.Render(c, "people_form", PersonFormVM{
			Person: clients.PersonDTO{Name: req.Name, Email: req.Email},
			Errors: vErrs,
		})
	}

	// Success: return a new row fragment (htmx swaps afterbegin into tbody)
	return h.r.Render(c, "people_row", person)
}

func (h *PeopleHandlers) EditPersonRow(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	person, err := h.pc.GetPerson(c.Context(), id)
	if err != nil {
		return c.Status(502).SendString("people service unavailable")
	}
	return h.r.Render(c, "people_row_edit", PersonEditRowVM{Person: person, Errors: map[string]string{}})
}

func (h *PeopleHandlers) UpdatePerson(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	req := clients.UpdatePersonRequest{
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
	}

	person, vErrs, err := h.pc.UpdatePerson(c.Context(), id, req)
	if err != nil {
		return c.Status(502).SendString("people service unavailable")
	}
	if len(vErrs) > 0 {
		return h.r.Render(c, "people_row_edit", PersonEditRowVM{
			Person: clients.PersonDTO{ID: id, Name: req.Name, Email: req.Email},
			Errors: vErrs,
		})
	}
	return h.r.Render(c, "people_row", person)
}

func (h *PeopleHandlers) DeletePerson(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := h.pc.DeletePerson(c.Context(), id); err != nil {
		return c.Status(502).SendString("people service unavailable")
	}
	return c.SendString("") // htmx will swap-out the row
}

func (h *PeopleHandlers) PersonRow(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	person, err := h.pc.GetPerson(c.Context(), id)
	if err != nil {
		return c.Status(502).SendString("people service unavailable")
	}
	return h.r.Render(c, "people_row", person)
}
