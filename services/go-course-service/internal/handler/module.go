package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-course-service/internal/models"
	"go-course-service/internal/repository"
)

type ModuleHandler struct {
	repo *repository.ModuleRepository
}

func NewModuleHandler(repo *repository.ModuleRepository) *ModuleHandler {
	return &ModuleHandler{repo: repo}
}

func (h *ModuleHandler) GetModulesByCourseID(c echo.Context) error {
	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course id"})
	}

	modules, err := h.repo.GetByCourseID(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, modules)
}

func (h *ModuleHandler) GetModuleByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid module id"})
	}

	module, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if module == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "module not found"})
	}

	return c.JSON(http.StatusOK, module)
}

type CreateModuleRequest struct {
	CourseID    int    `json:"courseId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"orderIndex"`
}

func (h *ModuleHandler) CreateModule(c echo.Context) error {
	var req CreateModuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Title == "" || req.CourseID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title and courseId are required"})
	}

	module, err := h.repo.Create(&models.Module{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, module)
}

func (h *ModuleHandler) UpdateModule(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid module id"})
	}

	var req CreateModuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	module, err := h.repo.Update(&models.Module{
		ID:          id,
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if module == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "module not found"})
	}
	return c.JSON(http.StatusOK, module)
}

func (h *ModuleHandler) DeleteModule(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid module id"})
	}

	err = h.repo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "module not found"})
	}
	return c.JSON(http.StatusNoContent, nil)
}
