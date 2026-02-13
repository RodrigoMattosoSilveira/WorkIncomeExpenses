package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/RodrigoMattosoSilveira/WorkEarningsExpenses/internal/bff/clients"
	"github.com/RodrigoMattosoSilveira/WorkEarningsExpenses/internal/bff/web"
)

func main() {
	app := fiber.New()

	renderer, err := web.NewRenderer("internal/bff/views/*.html")
	if err != nil {
		log.Fatal(err)
	}

	peopleClient := clients.NewPeopleClient(clients.PeopleClientConfig{
		BaseURL: "http://localhost:8081", // people-svc
		Timeout: 3 * time.Second,
	})

	h := web.NewPeopleHandlers(renderer, peopleClient)

	app.Get("/people", h.ListPeople)
	app.Get("/people/new", h.NewPersonForm)
	app.Post("/people", h.CreatePerson)
	app.Get("/people/:id/edit", h.EditPersonRow)
	app.Patch("/people/:id", h.UpdatePerson)
	app.Delete("/people/:id", h.DeletePerson)
	app.Get("/people/:id/row", h.PersonRow)

	log.Fatal(app.Listen(":8080"))
}
