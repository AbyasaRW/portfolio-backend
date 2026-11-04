package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRegisterDisabledByDefault(t *testing.T) {
	t.Setenv("ENABLE_REGISTRATION", "")
	app := fiber.New()
	SetupRoutes(app)

	response, err := app.Test(httptest.NewRequest(http.MethodPost, "/api/register", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestRegisterCanBeEnabled(t *testing.T) {
	t.Setenv("ENABLE_REGISTRATION", "true")
	app := fiber.New()
	SetupRoutes(app)

	response, err := app.Test(httptest.NewRequest(http.MethodPost, "/api/register", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}
