package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
