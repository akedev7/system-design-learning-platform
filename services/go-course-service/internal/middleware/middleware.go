package middleware

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go-course-service/internal/config"
)

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSCache struct {
	keys    map[string]*rsa.PublicKey
	mu      sync.RWMutex
	expire  time.Time
	jwksURL string
}

func NewJWKSCache(jwksURL string) *JWKSCache {
	return &JWKSCache{
		keys:    make(map[string]*rsa.PublicKey),
		jwksURL: jwksURL,
	}
}

func (c *JWKSCache) fetchKeys() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			rsaKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				continue
			}
			newKeys[key.Kid] = rsaKey
		}
	}

	c.mu.Lock()
	c.keys = newKeys
	c.expire = time.Now().Add(1 * time.Hour)
	c.mu.Unlock()

	return nil
}

func (c *JWKSCache) GetKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if key, ok := c.keys[kid]; ok && time.Now().Before(c.expire) {
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if key, ok := c.keys[kid]; ok && time.Now().Before(c.expire) {
		return key, nil
	}

	if err := c.fetchKeys(); err != nil {
		return nil, err
	}

	if key, ok := c.keys[kid]; ok {
		return key, nil
	}

	return nil, fmt.Errorf("key not found for kid: %s", kid)
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}

	// Convert bytes to int for exponent
	var eInt int
	if len(eBytes) == 0 {
		eInt = 0
	} else if len(eBytes) <= 4 {
		eInt = int(binary.BigEndian.Uint32(eBytes))
	} else {
		return nil, fmt.Errorf("exponent too large")
	}

	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}

	return pub, nil
}

func AuthMiddleware(cfg *config.AuthConfig) echo.MiddlewareFunc {
	cache := NewJWKSCache(cfg.JWKSURL)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			}

			tokenString := authHeader[7:]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("missing kid in token header")
				}

				return cache.GetKey(kid)
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
			}

			audience, ok := claims["aud"].([]interface{})
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid audience"})
			}

			audienceStr := make([]string, len(audience))
			for i, a := range audience {
				audienceStr[i] = fmt.Sprintf("%v", a)
			}

			hasAudience := false
			for _, a := range audienceStr {
				if a == cfg.Audience {
					hasAudience = true
					break
				}
			}
			if !hasAudience {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid audience"})
			}

			issuer, ok := claims["iss"].(string)
			if !ok || issuer != cfg.Issuer {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid issuer"})
			}

			auth0ID, ok := claims["sub"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid subject claim"})
			}

			c.Set("auth0_id", auth0ID)
			c.Set("user", token)
			return next(c)
		}
	}
}

func AdminMiddleware(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth0ID := c.Get("auth0_id")
			if auth0ID == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			var user struct {
				Role string `db:"role"`
			}
			err := db.Get(&user, "SELECT role FROM users WHERE auth0_id = $1", auth0ID)
			if err != nil {
				if err == sql.ErrNoRows {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "user not found"})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check role"})
			}

			if user.Role != "Admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "admin access required"})
			}

			return next(c)
		}
	}
}
