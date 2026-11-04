package controllers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"portfolio-backend/controllers"
)

func siteInfoBody(t *testing.T) string {
	t.Helper()
	app := fiber.New()
	app.Get("/site", controllers.GetSiteInfo)

	res, err := app.Test(httptest.NewRequest(http.MethodGet, "/site", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	return string(body)
}

func TestGetSiteInfo_Default(t *testing.T) {
	t.Setenv("SITE_URL", "")
	require.JSONEq(t, `{"site_url":"http://localhost:3000"}`, siteInfoBody(t))
}

func TestGetSiteInfo_FromEnvTrimsTrailingSlash(t *testing.T) {
	t.Setenv("SITE_URL", "  https://abyasa.example.com/  ")
	require.JSONEq(t, `{"site_url":"https://abyasa.example.com"}`, siteInfoBody(t))
}
