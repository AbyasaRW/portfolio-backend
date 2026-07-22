package controllers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"portfolio-backend/config"
	"portfolio-backend/controllers"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	config.DB = gormDB
	return mock
}

// ✅ GET /projects → berhasil ambil semua project
func TestGetProjects(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "tech_stack"}).
		AddRow(1, "Project A", "Deskripsi A", "Go")
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)

	app := fiber.New()
	app.Get("/projects", controllers.GetProjects)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	res, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), "Project A")
}

// ✅ GET /projects/:id → project ditemukan
func TestGetProjectByID_Found(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "tech_stack"}).
		AddRow(1, "Project A", "Deskripsi A", "Go")
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)

	app := fiber.New()
	app.Get("/projects/:id", controllers.GetProjectByID)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	res, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), "Project A")
}

// ❌ GET /projects/:id → project tidak ditemukan
func TestGetProjectByID_NotFound(t *testing.T) {
	mock := setupMockDB(t)

	// DB tidak kembalikan baris apapun → Not Found
	rows := sqlmock.NewRows([]string{"id", "title", "description", "tech_stack"})
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)

	app := fiber.New()
	app.Get("/projects/:id", controllers.GetProjectByID)

	req := httptest.NewRequest(http.MethodGet, "/projects/999", nil)
	res, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), "tidak ditemukan")
}

// ❌ POST /projects → payload kosong, title tidak ada
func TestCreateProject_BadPayload(t *testing.T) {
	setupMockDB(t)

	app := fiber.New()
	app.Post("/projects", controllers.CreateProject)

	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), "Title tidak boleh kosong")
}

// ✅ POST /projects → berhasil buat project baru
func TestCreateProject_Success(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "projects"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := fiber.New()
	app.Post("/projects", controllers.CreateProject)

	payload := `{"title":"My App","description":"Keren","tech_stack":"Go"}`
	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	require.Equal(t, http.StatusCreated, res.StatusCode)
}

// ❌ DELETE /projects/:id → project tidak ditemukan
func TestDeleteProject_NotFound(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title"})
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)

	app := fiber.New()
	app.Delete("/projects/:id", controllers.DeleteProject)

	req := httptest.NewRequest(http.MethodDelete, "/projects/999", nil)
	res, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

// ✅ DELETE /projects/:id → berhasil hapus
func TestDeleteProject_Success(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "tech_stack"}).
		AddRow(1, "Project A", "Deskripsi A", "Go")
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	app := fiber.New()
	app.Delete("/projects/:id", controllers.DeleteProject)

	req := httptest.NewRequest(http.MethodDelete, "/projects/1", nil)
	res, err := app.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), "berhasil dihapus")
}
