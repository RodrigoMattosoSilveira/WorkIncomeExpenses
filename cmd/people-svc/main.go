package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/RodrigoMattosoSilveira/WorkIncomeExpenses/internal/people"
)

func main() {
	db, err := gorm.Open(sqlite.Open("people.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&people.Person{}); err != nil {
		log.Fatal(err)
	}

	repo := people.NewGormRepo(db)
	svc := people.NewService(repo)
	api := people.NewAPI(svc)

	app := fiber.New()
	api.Register(app)

	log.Fatal(app.Listen(":8081"))
}
