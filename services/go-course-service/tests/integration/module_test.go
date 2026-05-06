package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go-course-service/internal/database"
	"go-course-service/internal/handler"
	"go-course-service/internal/repository"
)

func setupTestDB(t *testing.T) (*repository.ModuleRepository, func()) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithDatabase("test"),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := database.NewConnection(&database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := database.RunMigrations(&database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	}, "file://../../migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	repo := repository.NewModuleRepository(db)

	cleanup := func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}

	return repo, cleanup
}

func TestGetModulesByCourseID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

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
	repo, cleanup := setupTestDB(t)
	defer cleanup()

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
