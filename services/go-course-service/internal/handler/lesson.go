package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-course-service/internal/models"
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

type CreateLessonRequest struct {
	ModuleID    int                  `json:"moduleId"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	ContentJSON json.RawMessage      `json:"contentJsonb"`
	OrderIndex  int                  `json:"orderIndex"`
}

func (h *LessonHandler) CreateLesson(c echo.Context) error {
	var req CreateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Title == "" || req.ModuleID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title and moduleId are required"})
	}

	lesson, err := h.repo.Create(&models.Lesson{
		ModuleID:    req.ModuleID,
		Title:       req.Title,
		Description: req.Description,
		ContentJSON: req.ContentJSON,
		OrderIndex:  req.OrderIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, lesson)
}

func (h *LessonHandler) UpdateLesson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lesson id"})
	}

	var req CreateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	lesson, err := h.repo.Update(&models.Lesson{
		ID:          id,
		ModuleID:    req.ModuleID,
		Title:       req.Title,
		Description: req.Description,
		ContentJSON: req.ContentJSON,
		OrderIndex:  req.OrderIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}
	return c.JSON(http.StatusOK, lesson)
}

func (h *LessonHandler) DeleteLesson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lesson id"})
	}

	err = h.repo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

type UpdateContentRequest struct {
	ContentJSON json.RawMessage `json:"contentJsonb"`
}

func (h *LessonHandler) UpdateContent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lesson id"})
	}

	var req UpdateContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	lesson, err := h.repo.UpdateContent(id, req.ContentJSON)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}

	return c.JSON(http.StatusOK, lesson)
}
