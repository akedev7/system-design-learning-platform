package handler

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go-course-service/internal/models"
	"go-course-service/internal/repository"
	"go-course-service/internal/service"
	"go-course-service/internal/response"
)

type ProgressHandler struct {
	progressRepo *repository.ProgressRepository
	lessonRepo   *repository.LessonRepository
	moduleRepo   *repository.ModuleRepository
	courseRepo   *repository.CourseRepository
	progressService *service.ProgressTrackingService
}

func NewProgressHandler(progressRepo *repository.ProgressRepository, lessonRepo *repository.LessonRepository, moduleRepo *repository.ModuleRepository, courseRepo *repository.CourseRepository) *ProgressHandler {
	return &ProgressHandler{
		progressRepo: progressRepo,
		lessonRepo:   lessonRepo,
		moduleRepo:   moduleRepo,
		courseRepo:   courseRepo,
		progressService: service.NewProgressTrackingService(),
	}
}

func (h *ProgressHandler) GetCourseProgress(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid course id"))
	}

	modules, err := h.moduleRepo.GetByCourseID(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	var allLessonIDs []int
	moduleProgressData := []models.ModuleProgress{}

	for _, module := range modules {
		lessons, err := h.lessonRepo.GetByModuleID(module.ID)
		if err != nil {
			continue
		}

		var lessonIDs []int
		for _, lesson := range lessons {
			lessonIDs = append(lessonIDs, lesson.ID)
		}
		allLessonIDs = append(allLessonIDs, lessonIDs...)

		lessonProgress, _ := h.progressRepo.GetLessonProgressByUser(userID, lessonIDs)
		var lp []models.LessonProgress
		for _, p := range lessonProgress {
			lp = append(lp, models.LessonProgress{
				LessonID:    p.LessonID,
				Score:       p.Score,
				Passed:      p.Passed,
				CompletedAt: p.CompletedAt,
			})
		}

		mp := h.progressService.CalculateModuleCompletion(lp)
		mp.ModuleID = module.ID
		moduleProgressData = append(moduleProgressData, mp)
	}

	cp := h.progressService.CalculateCourseCompletion(moduleProgressData)
	cp.CourseID = courseID

	return c.JSON(http.StatusOK, response.Success(cp))
}

func (h *ProgressHandler) GetModuleProgress(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
	}

	moduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid module id"))
	}

	lessons, err := h.lessonRepo.GetByModuleID(moduleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	var lessonIDs []int
	for _, lesson := range lessons {
		lessonIDs = append(lessonIDs, lesson.ID)
	}

	lessonProgress, _ := h.progressRepo.GetLessonProgressByUser(userID, lessonIDs)
	var lp []models.LessonProgress
	for _, p := range lessonProgress {
		lp = append(lp, models.LessonProgress{
			LessonID:    p.LessonID,
			Score:       p.Score,
			Passed:      p.Passed,
			CompletedAt: p.CompletedAt,
		})
	}

	result := h.progressService.CalculateModuleCompletion(lp)
	result.ModuleID = moduleID

	return c.JSON(http.StatusOK, response.Success(result))
}

func (h *ProgressHandler) GetLessonProgress(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
	}

	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid lesson id"))
	}

	progress, err := h.progressRepo.GetProgress(userID, lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	if progress == nil {
		return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
			"lessonId": lessonID,
			"score":    0,
			"passed":   false,
		}))
	}

	return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
		"lessonId":    progress.LessonID,
		"score":       progress.Score,
		"passed":      progress.Passed,
		"completedAt": progress.CompletedAt,
	}))
}

func (h *ProgressHandler) GetResumeLesson(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
	}

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid course id"))
	}

	modules, err := h.moduleRepo.GetByCourseID(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	for _, module := range modules {
		lessons, err := h.lessonRepo.GetByModuleID(module.ID)
		if err != nil {
			continue
		}

		var lessonIDs []int
		for _, lesson := range lessons {
			lessonIDs = append(lessonIDs, lesson.ID)
		}

		lessonProgress, _ := h.progressRepo.GetLessonProgressByUser(userID, lessonIDs)
		var lp []models.LessonProgress
		for _, p := range lessonProgress {
			lp = append(lp, models.LessonProgress{
				LessonID:    p.LessonID,
				Score:       p.Score,
				Passed:      p.Passed,
				CompletedAt: p.CompletedAt,
			})
		}

		lastLesson := h.progressService.GetLastAccessedLesson(lp)
		if lastLesson.LessonID > 0 {
			return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
				"moduleId": module.ID,
				"lessonId": lastLesson.LessonID,
				"passed":   lastLesson.Passed,
			}))
		}

		if len(lessons) > 0 {
			return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
				"moduleId": module.ID,
				"lessonId": lessons[0].ID,
				"passed":   false,
			}))
		}
	}

	return c.JSON(http.StatusNotFound, response.Error("no lessons found"))
}