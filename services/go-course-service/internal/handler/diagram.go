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
)

type DiagramHandler struct {
	lessonRepo     *repository.LessonRepository
	progressRepo   *repository.ProgressRepository
	diagramService *service.DiagramValidationService
}

func NewDiagramHandler(lessonRepo *repository.LessonRepository, progressRepo *repository.ProgressRepository) *DiagramHandler {
	return &DiagramHandler{
		lessonRepo:     lessonRepo,
		progressRepo:   progressRepo,
		diagramService: service.NewDiagramValidationService(),
	}
}

func (h *DiagramHandler) ValidateDiagram(c echo.Context) error {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid lesson id"})
	}

	lesson, err := h.lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}

	var req models.DiagramValidationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	diagramConfig, err := extractDiagramConfigFromContent(lesson.ContentJSON)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no diagram found in lesson content"})
	}

	if diagramConfig == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no diagram configuration in lesson"})
	}

	result := h.diagramService.ValidateDiagram(*diagramConfig, req.Diagram)

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	score := result.Score
	passed := result.Valid

	_, err = h.progressRepo.UpsertProgress(sub, lessonID, score, passed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save progress"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid":  result.Valid,
		"score":  result.Score,
		"errors": result.Errors,
	})
}

func extractDiagramConfigFromContent(contentJSON json.RawMessage) (*models.DiagramConfig, error) {
	var blocks []models.ContentBlock
	if err := json.Unmarshal(contentJSON, &blocks); err != nil {
		return nil, err
	}

	for _, block := range blocks {
		if block.Type == "ReactFlowDiagram" {
			var config models.DiagramConfig
			if err := json.Unmarshal(block.Config, &config); err != nil {
				return nil, err
			}
			return &config, nil
		}
	}

	return nil, nil
}