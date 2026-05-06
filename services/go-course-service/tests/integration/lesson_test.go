package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go-course-service/internal/handler"
	"go-course-service/internal/repository"
)

func TestGetLessonsByModuleID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := handler.NewLessonHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/modules/1/lessons", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("moduleId")
	c.SetParamValues("1")

	err := h.GetLessonsByModuleID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetLessonByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := handler.NewLessonHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := h.GetLessonByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
