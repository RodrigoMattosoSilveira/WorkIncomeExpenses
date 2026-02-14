package web

import (
	"net/http/httptest"
	"testing"
)	

func TestHome_Golden(t *testing.T) {
	app := newBFFAppForTest(t)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	got := readBody(t, resp.Body)
	assertGolden(t, "home.html", got)
}
