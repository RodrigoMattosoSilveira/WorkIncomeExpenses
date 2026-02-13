package people_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/RodrigoMattosoSilveira/WorkEarningsExpenses/internal/people"
)

func newTestApp(t *testing.T) *fiber.App {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&people.Person{}); err != nil {
		t.Fatal(err)
	}

	repo := people.NewGormRepo(db)
	svc := people.NewService(repo)
	api := people.NewAPI(svc)

	app := fiber.New()
	api.Register(app)
	return app
}

func TestCreateAndListPeople(t *testing.T) {
	app := newTestApp(t)

	// Create
	body := []byte(`{"name":"Ada","email":"ada@example.com"}`)
	req := httptest.NewRequest("POST", "/api/people", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// List
	req2 := httptest.NewRequest("GET", "/api/people", nil)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}

	var got []people.Person
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Email != "ada@example.com" {
		t.Fatalf("unexpected list: %+v", got)
	}
}

func TestCreateValidation(t *testing.T) {
	app := newTestApp(t)

	body := []byte(`{"name":"","email":"not-an-email"}`)
	req := httptest.NewRequest("POST", "/api/people", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}

	var v struct {
		Errors map[string]string `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatal(err)
	}
	if v.Errors["name"] == "" || v.Errors["email"] == "" {
		t.Fatalf("expected field errors, got %+v", v.Errors)
	}
}
