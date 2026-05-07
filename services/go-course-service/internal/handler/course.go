package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-course-service/internal/models"
	"go-course-service/internal/repository"
	"go-course-service/internal/response"
)

type CourseHandler struct {
	repo *repository.CourseRepository
}

func NewCourseHandler(repo *repository.CourseRepository) *CourseHandler {
	return &CourseHandler{repo: repo}
}

func (h *CourseHandler) GetCourses(c echo.Context) error {
	courses, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Success(courses))
}

func (h *CourseHandler) GetCourseByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid course id"))
	}

	course, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, response.Error("course not found"))
	}
	return c.JSON(http.StatusOK, response.Success(course))
}

type CreateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *CourseHandler) CreateCourse(c echo.Context) error {
	var req CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, response.Error("title is required"))
	}

	course, err := h.repo.Create(&models.Course{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	return c.JSON(http.StatusCreated, response.Success(course))
}

func (h *CourseHandler) UpdateCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid course id"))
	}

	var req CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, response.Error("title is required"))
	}

	course, err := h.repo.Update(&models.Course{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, response.Error("course not found"))
	}
	return c.JSON(http.StatusOK, response.Success(course))
}

func (h *CourseHandler) DeleteCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid course id"))
	}

	err = h.repo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("course not found"))
	}
	return c.NoContent(http.StatusNoContent)
}