package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-course-service/internal/handler"
	"go-course-service/internal/repository"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := runMigrationsWithConnStr(connStr, "file://../../migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func runMigrationsWithConnStr(connStr, migrationsPath string) error {
	src, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to open migrations path: %w", err)
	}
	defer src.Close()

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	driver, err := postgresMigrate.WithInstance(db.DB, &postgresMigrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}
	defer driver.Close()

	m, err := migrate.NewWithInstance("file", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func TestGetModulesByCourseID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewModuleRepository(db)

	e := echo.New()
	h := handler.NewModuleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/1/modules", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("courseId")
	c.SetParamValues("1")

	err := h.GetModulesByCourseID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetModuleByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewModuleRepository(db)

	e := echo.New()
	h := handler.NewModuleHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/modules/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := h.GetModuleByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
