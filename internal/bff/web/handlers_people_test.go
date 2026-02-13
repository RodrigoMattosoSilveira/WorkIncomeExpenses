package web

import (
	"context"
	"flag"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"

	"github.com/RodrigoMattosoSilveira/WorkEarningsExpenses/internal/bff/clients"
)

var update = flag.Bool("update", false, "update golden files")

type fakePeople struct {
	list []clients.PersonDTO
	get  map[int64]clients.PersonDTO
}

func (f fakePeople) ListPeople(ctx context.Context) ([]clients.PersonDTO, error) {
	return f.list, nil
}
func (f fakePeople) GetPerson(ctx context.Context, id int64) (clients.PersonDTO, error) {
	return f.get[id], nil
}
func (f fakePeople) CreatePerson(ctx context.Context, in clients.CreatePersonRequest) (clients.PersonDTO, map[string]string, error) {
	// for fragment tests we can just return deterministic output
	p := clients.PersonDTO{ID: 99, Name: in.Name, Email: in.Email}
	return p, nil, nil
}
func (f fakePeople) UpdatePerson(ctx context.Context, id int64, in clients.UpdatePersonRequest) (clients.PersonDTO, map[string]string, error) {
	p := clients.PersonDTO{ID: id, Name: in.Name, Email: in.Email}
	return p, nil, nil
}
func (f fakePeople) DeletePerson(ctx context.Context, id int64) error { return nil }

func newBFFAppForTest(t *testing.T) *fiber.App {
	t.Helper()

	// Parse templates from disk (same as runtime)
	r, err := NewRenderer("internal/bff/views/*.html")
	if err != nil {
		t.Fatal(err)
	}

	fp := fakePeople{
		list: []clients.PersonDTO{
			{ID: 1, Name: "Ada", Email: "ada@example.com"},
		},
		get: map[int64]clients.PersonDTO{
			1: {ID: 1, Name: "Ada", Email: "ada@example.com"},
		},
	}

	h := NewPeopleHandlers(r, fp) // fp satisfies clients.PeopleAPI now

	app := fiber.New()
	app.Get("/people", h.ListPeople)
	app.Get("/people/:id/row", h.PersonRow)
	app.Get("/people/:id/edit", h.EditPersonRow)
	return app
}

func readBody(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func assertGolden(t *testing.T, goldenName string, got string) {
	t.Helper()
	path := filepath.Join("internal", "testdata", "golden", goldenName)

	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	wantB, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run `go test ./... -update`)", path, err)
	}
	want := string(wantB)

	if got != want {
		t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", goldenName, got, want)
	}
}

func TestPeopleTbodyFragment_Golden(t *testing.T) {
	app := newBFFAppForTest(t)

	req := httptest.NewRequest("GET", "/people", nil)
	req.Header.Set("HX-Request", "true")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	got := readBody(t, resp.Body)
	assertGolden(t, "people_tbody.html", got)
}

func TestPeopleRowFragment_Golden(t *testing.T) {
	app := newBFFAppForTest(t)

	req := httptest.NewRequest("GET", "/people/1/row", nil)
	req.Header.Set("HX-Request", "true")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	got := readBody(t, resp.Body)
	assertGolden(t, "people_row.html", got)
}

func TestPeopleRowEditFragment_Golden(t *testing.T) {
	app := newBFFAppForTest(t)

	req := httptest.NewRequest("GET", "/people/1/edit", nil)
	req.Header.Set("HX-Request", "true")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	got := readBody(t, resp.Body)

	// normalize any trivial whitespace diffs if you want; here we compare exact
	assertGolden(t, "people_row_edit.html", got)
}
