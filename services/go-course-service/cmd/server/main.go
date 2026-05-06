package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go-course-service/internal/config"
	"go-course-service/internal/database"
	"go-course-service/internal/handler"
	appMiddleware "go-course-service/internal/middleware"
	"go-course-service/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunMigrations(&cfg.Database, "file://migrations"); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Debug = cfg.Server.Debug
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Health endpoint (public)
	e.GET("/actuator/health", func(c echo.Context) error {
		if err := db.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "DOWN"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})

	// Info endpoint (public)
	e.GET("/actuator/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"service": "course-service",
			"version": "1.0.0",
		})
	})

	// API routes with auth middleware
	api := e.Group("/api/v1")
	api.Use(appMiddleware.AuthMiddleware(&cfg.Auth))

	// Course handlers (stubs)
	courseHandler := handler.NewCourseHandler()
	api.GET("/courses", courseHandler.GetCourses)
	api.GET("/courses/:id", courseHandler.GetCourseByID)

	// Module handlers
	moduleRepo := repository.NewModuleRepository(db)
	moduleHandler := handler.NewModuleHandler(moduleRepo)
	api.GET("/courses/:courseId/modules", moduleHandler.GetModulesByCourseID)
	api.GET("/modules/:id", moduleHandler.GetModuleByID)

	// Lesson handlers
	lessonRepo := repository.NewLessonRepository(db)
	lessonHandler := handler.NewLessonHandler(lessonRepo)
	api.GET("/modules/:moduleId/lessons", lessonHandler.GetLessonsByModuleID)
	api.GET("/lessons/:id", lessonHandler.GetLessonByID)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}
