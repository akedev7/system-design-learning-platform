package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go-course-service/internal/config"
	"go-course-service/internal/middleware"
)

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := middleware.AuthMiddleware(&config.AuthConfig{
		JWKSURL:  "https://test.auth0.com/.well-known/jwks.json",
		Audience: "test-audience",
		Issuer:   "https://test.auth0.com/",
	})

	h := mw(func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
