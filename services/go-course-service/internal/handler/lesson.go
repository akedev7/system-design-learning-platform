package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-course-service/internal/repository"
)

type LessonHandler struct {
	repo *repository.LessonRepository
}

func NewLessonHandler(repo *repository.LessonRepository) *LessonHandler {
	return &LessonHandler{repo: repo}
}

func (h *LessonHandler) GetLessonsByModuleID(c echo.Context) error {
	moduleID, err := strconv.Atoi(c.Param("moduleId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid module id"})
	}

	lessons, err := h.repo.GetByModuleID(moduleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, lessons)
}

func (h *LessonHandler) GetLessonByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lesson id"})
	}

	lesson, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}

	return c.JSON(http.StatusOK, lesson)
}
