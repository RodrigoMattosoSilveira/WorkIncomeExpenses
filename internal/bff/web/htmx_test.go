package web
import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

const (
	HTMX_REQUEST   = "HTMX Request"
	NOT_HTMX_REQUEST   = "Not HTMX Request"
)
// Handler we want to test
func myHandler(c fiber.Ctx) error {
	if IsHTMX(c) {
		return c.SendString(HTMX_REQUEST)
	}
	return c.SendString(NOT_HTMX_REQUEST)
}
func TestIsHTMX(t *testing.T) {
	app := fiber.New()
	app.Post("/test", myHandler)

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/test?name=john", strings.NewReader(`{"key":"value"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Header", "expected")

	// Perform the request
	resp, err := app.Test(req) // -1 disables request timeout
	assert.NoError(t, err)

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Assertions
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, NOT_HTMX_REQUEST, bodyStr)


	// Create a test request
	req = httptest.NewRequest(http.MethodPost, "/test?name=john", strings.NewReader(`{"key":"value"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Header", "expected")
	req.Header.Set("HX-Request", "true")

	// Perform the request
	resp, err = app.Test(req) // -1 disables request timeout
	assert.NoError(t, err)

	// Read response body
	bodyBytes, _ = io.ReadAll(resp.Body)
	bodyStr = string(bodyBytes)

	// Assertions
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, HTMX_REQUEST, bodyStr)

}