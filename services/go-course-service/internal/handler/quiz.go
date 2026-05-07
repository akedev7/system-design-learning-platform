package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go-course-service/internal/models"
	"go-course-service/internal/repository"
	"go-course-service/internal/service"
	"go-course-service/internal/response"
)

type QuizHandler struct {
	lessonRepo  *repository.LessonRepository
	progressRepo *repository.ProgressRepository
	quizService *service.QuizService
}

func NewQuizHandler(lessonRepo *repository.LessonRepository, progressRepo *repository.ProgressRepository) *QuizHandler {
	return &QuizHandler{
		lessonRepo:   lessonRepo,
		progressRepo: progressRepo,
		quizService:  service.NewQuizService(),
	}
}

func (h *QuizHandler) SubmitQuiz(c echo.Context) error {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid lesson id"))
	}

	lesson, err := h.lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	if lesson == nil {
		return c.JSON(http.StatusNotFound, response.Error("lesson not found"))
	}

	var req models.QuizAnswerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	quiz, err := extractQuizFromContent(lesson.ContentJSON)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("no quiz found in lesson content"))
	}

	result := h.quizService.GradeQuiz(*quiz, req.Answers)

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
	}

	progress, err := h.progressRepo.UpsertProgress(sub, lessonID, result.Score, result.Passed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to save progress"))
	}

	return c.JSON(http.StatusOK, response.Success(map[string]interface{}{
		"result":   result,
		"progress": progress,
	}))
}

func extractQuizFromContent(contentJSON json.RawMessage) (*models.Quiz, error) {
	var blocks []models.ContentBlock
	if err := json.Unmarshal(contentJSON, &blocks); err != nil {
		return nil, err
	}

	for _, block := range blocks {
		if block.Type == "Quiz" {
			var quiz models.Quiz
			if err := json.Unmarshal(block.Config, &quiz); err != nil {
				return nil, err
			}
			return &quiz, nil
		}
	}

	return nil, nil
}