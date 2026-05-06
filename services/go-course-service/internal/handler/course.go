package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CourseHandler struct{}

func NewCourseHandler() *CourseHandler {
	return &CourseHandler{}
}

func (h *CourseHandler) GetCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, []interface{}{})
}

func (h *CourseHandler) GetCourseByID(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
