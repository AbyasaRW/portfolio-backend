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

// GET /projects harus urut berdasarkan sort_order lalu id.
func TestGetProjects_OrderedBySortOrder(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "sort_order"}).
		AddRow(2, "Kedua", 2).
		AddRow(1, "Pertama", 1)
	mock.ExpectQuery(`SELECT \* FROM "projects" .*ORDER BY sort_order ASC, id ASC`).WillReturnRows(rows)

	app := fiber.New()
	app.Get("/projects", controllers.GetProjects)

	res, err := app.Test(httptest.NewRequest(http.MethodGet, "/projects", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, mock.ExpectationsWereMet())
}

// POST /projects: judul berisi spasi saja dianggap kosong.
func TestCreateProject_WhitespaceTitleRejected(t *testing.T) {
	setupMockDB(t)

	app := fiber.New()
	app.Post("/projects", controllers.CreateProject)

	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(`{"title":"   "}`))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// POST /projects: accent/image_fit tidak valid dikembalikan ke default.
func TestCreateProject_InvalidViewFieldsNormalized(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "projects"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := fiber.New()
	app.Post("/projects", controllers.CreateProject)

	payload := `{"title":"App","accent":"hotpink","image_fit":"stretch","category":"Web"}`
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), `"accent":"blue"`)
	require.Contains(t, string(body), `"image_fit":"cover"`)
	require.Contains(t, string(body), `"category":"Web"`)
}

// PUT /projects/:id: body tidak boleh menimpa ID asli.
func TestUpdateProject_CannotOverwriteID(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(7, "Lama")
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "projects"`).WillReturnResult(sqlmock.NewResult(7, 1))
	mock.ExpectCommit()

	app := fiber.New()
	app.Put("/projects/:id", controllers.UpdateProject)

	req := httptest.NewRequest(http.MethodPut, "/projects/7", strings.NewReader(`{"ID":999,"title":"Baru"}`))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	require.Contains(t, string(body), `"ID":7`)
	require.NotContains(t, string(body), `"ID":999`)
	require.Contains(t, string(body), `"title":"Baru"`)
}

// PUT /projects/:id: judul dikosongkan ditolak.
func TestUpdateProject_EmptyTitleRejected(t *testing.T) {
	mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Lama")
	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(rows)

	app := fiber.New()
	app.Put("/projects/:id", controllers.UpdateProject)

	req := httptest.NewRequest(http.MethodPut, "/projects/1", strings.NewReader(`{"title":""}`))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// PUT /projects/:id: project tidak ada -> 404.
func TestUpdateProject_NotFound(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	app := fiber.New()
	app.Put("/projects/:id", controllers.UpdateProject)

	req := httptest.NewRequest(http.MethodPut, "/projects/404", strings.NewReader(`{"title":"x"}`))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}