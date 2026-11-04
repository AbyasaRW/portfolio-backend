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

	"portfolio-backend/controllers"
)

func contentApp() *fiber.App {
	app := fiber.New()
	app.Get("/experiences", controllers.GetExperiences)
	app.Post("/experiences", controllers.CreateExperience)
	app.Put("/experiences/:id", controllers.UpdateExperience)
	app.Delete("/experiences/:id", controllers.DeleteExperience)
	app.Get("/tools", controllers.GetTools)
	app.Post("/tools", controllers.CreateTool)
	app.Put("/tools/:id", controllers.UpdateTool)
	app.Delete("/tools/:id", controllers.DeleteTool)
	return app
}

func do(t *testing.T, app *fiber.App, method, path, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := app.Test(req)
	require.NoError(t, err)
	out, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(out)
}

// ---------- GET ----------

func TestGetExperiences_OrderedAndEmptyIsArray(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "experiences" .*ORDER BY sort_order ASC, id ASC`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company"}))

	code, body := do(t, contentApp(), http.MethodGet, "/experiences", "")

	require.Equal(t, http.StatusOK, code)
	require.JSONEq(t, `[]`, body) // kosong harus [] bukan null
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTools_ReturnsRows(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "tools" .*ORDER BY sort_order ASC, id ASC`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tool_group", "name"}).AddRow(1, "frontend", "Next.js"))

	code, body := do(t, contentApp(), http.MethodGet, "/tools", "")

	require.Equal(t, http.StatusOK, code)
	require.Contains(t, body, `"name":"Next.js"`)
	require.Contains(t, body, `"group":"frontend"`) // kolom tool_group, JSON tetap "group"
	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------- POST ----------

func TestCreateExperience_RequiresCompany(t *testing.T) {
	setupMockDB(t)
	code, body := do(t, contentApp(), http.MethodPost, "/experiences", `{"company":"   "}`)
	require.Equal(t, http.StatusBadRequest, code)
	require.Contains(t, body, "Company tidak boleh kosong")
}

func TestCreateExperience_NormalizesAndIgnoresBodyID(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "experiences"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
	mock.ExpectCommit()

	payload := `{"ID":999,"company":"  MNCTV ","categories":"Broadcast, ngawur","tone":"pink","skills":"a, ,b"}`
	code, body := do(t, contentApp(), http.MethodPost, "/experiences", payload)

	require.Equal(t, http.StatusCreated, code)
	require.Contains(t, body, `"company":"MNCTV"`)
	require.Contains(t, body, `"categories":"broadcast"`)
	require.Contains(t, body, `"tone":"blue"`)
	require.Contains(t, body, `"skills":"a, b"`)
	require.NotContains(t, body, `"ID":999`)
}

func TestCreateTool_RequiresName(t *testing.T) {
	setupMockDB(t)
	code, body := do(t, contentApp(), http.MethodPost, "/tools", `{"group":"dev"}`)
	require.Equal(t, http.StatusBadRequest, code)
	require.Contains(t, body, "Name tidak boleh kosong")
}

func TestCreateTool_InvalidGroupBecomesFrontend(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "tools"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	code, body := do(t, contentApp(), http.MethodPost, "/tools", `{"name":"Go","group":"lainnya"}`)

	require.Equal(t, http.StatusCreated, code)
	require.Contains(t, body, `"group":"frontend"`)
}

// ---------- PUT ----------

func TestUpdateTool_CannotOverwriteID(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "tools"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(7, "Lama"))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "tools"`).WillReturnResult(sqlmock.NewResult(7, 1))
	mock.ExpectCommit()

	code, body := do(t, contentApp(), http.MethodPut, "/tools/7", `{"ID":999,"name":"Baru"}`)

	require.Equal(t, http.StatusOK, code)
	require.Contains(t, body, `"ID":7`)
	require.NotContains(t, body, `"ID":999`)
	require.Contains(t, body, `"name":"Baru"`)
}

func TestUpdateExperience_EmptyCompanyRejected(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "experiences"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company"}).AddRow(1, "Lama"))

	code, _ := do(t, contentApp(), http.MethodPut, "/experiences/1", `{"company":""}`)
	require.Equal(t, http.StatusBadRequest, code)
}

func TestUpdateTool_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "tools"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	code, body := do(t, contentApp(), http.MethodPut, "/tools/404", `{"name":"x"}`)
	require.Equal(t, http.StatusNotFound, code)
	require.Contains(t, body, "tidak ditemukan")
}

// ---------- DELETE ----------

func TestDeleteExperience_Success(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "experiences"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "company"}).AddRow(1, "MNCTV"))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "experiences"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	code, body := do(t, contentApp(), http.MethodDelete, "/experiences/1", "")

	require.Equal(t, http.StatusOK, code)
	require.Contains(t, body, "berhasil dihapus")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTool_NotFound(t *testing.T) {
	mock := setupMockDB(t)
	mock.ExpectQuery(`SELECT \* FROM "tools"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	code, _ := do(t, contentApp(), http.MethodDelete, "/tools/404", "")
	require.Equal(t, http.StatusNotFound, code)
}