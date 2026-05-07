package handler

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-course-service/internal/response"
	"go-course-service/internal/storage"
)

type UploadHandler struct {
	s3Client *storage.S3Client
}

func NewUploadHandler(s3Client *storage.S3Client) *UploadHandler {
	return &UploadHandler{s3Client: s3Client}
}

type GenerateUploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

type GenerateUploadURLResponse struct {
	UploadURL string `json:"uploadURL"`
	Key       string `json:"key"`
	PublicURL string `json:"publicURL"`
}

func (h *UploadHandler) GenerateUploadURL(c echo.Context) error {
	var req GenerateUploadURLRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if req.Filename == "" {
		return c.JSON(http.StatusBadRequest, response.Error("filename is required"))
	}

	contentType := req.ContentType
	if contentType == "" {
		ext := strings.ToLower(filepath.Ext(req.Filename))
		contentType = mimeTypeFromExt(ext)
	}

	ext := strings.ToLower(filepath.Ext(req.Filename))
	key := generateStorageKey(ext)

	uploadURL, err := h.s3Client.GeneratePresignedUpload(c.Request().Context(), key, contentType, 15*time.Minute)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to generate upload URL"))
	}

	publicURL := h.s3Client.GetPublicURL(key)

	return c.JSON(http.StatusOK, response.Success(GenerateUploadURLResponse{
		UploadURL: uploadURL,
		Key:       key,
		PublicURL: publicURL,
	}))
}

func mimeTypeFromExt(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}

func generateStorageKey(ext string) string {
	uuid := uuid.New().String()
	return "uploads/" + uuid + ext
}