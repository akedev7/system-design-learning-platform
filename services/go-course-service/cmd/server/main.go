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
	"go-course-service/internal/storage"
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
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
	}))

	e.GET("/actuator/health", func(c echo.Context) error {
		if err := db.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "DOWN"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})

	e.GET("/actuator/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"service": "course-service",
			"version": "1.0.0",
		})
	})

	api := e.Group("/api/v1")
	// Bypass auth for local development if DISABLE_AUTH=true
	if os.Getenv("DISABLE_AUTH") != "true" {
		api.Use(appMiddleware.AuthMiddleware(&cfg.Auth))
	}

	courseRepo := repository.NewCourseRepository(db)
	courseHandler := handler.NewCourseHandler(courseRepo)
	api.GET("/courses", courseHandler.GetCourses)
	api.GET("/courses/:id", courseHandler.GetCourseByID)

	admin := api.Group("")
	// Bypass admin check for local development if DISABLE_AUTH=true
	if os.Getenv("DISABLE_AUTH") != "true" {
		admin.Use(appMiddleware.AdminMiddleware(db))
	}
	admin.POST("/courses", courseHandler.CreateCourse)
	admin.PUT("/courses/:id", courseHandler.UpdateCourse)
	admin.DELETE("/courses/:id", courseHandler.DeleteCourse)

	moduleRepo := repository.NewModuleRepository(db)
	moduleHandler := handler.NewModuleHandler(moduleRepo)
	api.GET("/courses/:courseId/modules", moduleHandler.GetModulesByCourseID)
	api.GET("/modules/:id", moduleHandler.GetModuleByID)

	admin.POST("/modules", moduleHandler.CreateModule)
	admin.PUT("/modules/:id", moduleHandler.UpdateModule)
	admin.DELETE("/modules/:id", moduleHandler.DeleteModule)

	lessonRepo := repository.NewLessonRepository(db)
	lessonHandler := handler.NewLessonHandler(lessonRepo)
	api.GET("/modules/:moduleId/lessons", lessonHandler.GetLessonsByModuleID)
	api.GET("/lessons/:id", lessonHandler.GetLessonByID)

	admin.POST("/lessons", lessonHandler.CreateLesson)
	admin.PUT("/lessons/:id", lessonHandler.UpdateLesson)
	admin.PUT("/lessons/:id/content", lessonHandler.UpdateContent)
	admin.DELETE("/lessons/:id", lessonHandler.DeleteLesson)

	progressRepo := repository.NewProgressRepository(db)
	quizHandler := handler.NewQuizHandler(lessonRepo, progressRepo)
	api.POST("/lessons/:id/submit-quiz", quizHandler.SubmitQuiz)

	diagramHandler := handler.NewDiagramHandler(lessonRepo, progressRepo)
	api.POST("/lessons/:id/validate-diagram", diagramHandler.ValidateDiagram)

	progressHandler := handler.NewProgressHandler(progressRepo, lessonRepo, moduleRepo, courseRepo)
	api.GET("/courses/:id/progress", progressHandler.GetCourseProgress)
	api.GET("/modules/:id/progress", progressHandler.GetModuleProgress)
	api.GET("/lessons/:id/progress", progressHandler.GetLessonProgress)
	api.GET("/courses/:id/resume", progressHandler.GetResumeLesson)
	api.POST("/courses/:id/enroll", progressHandler.EnrollCourse)

	var uploadHandler *handler.UploadHandler
	if cfg.S3.Endpoint != "" {
		s3Client, err := storage.NewS3Client(&cfg.S3)
		if err != nil {
			fmt.Printf("Warning: Failed to initialize S3 client: %v\n", err)
		} else {
			uploadHandler = handler.NewUploadHandler(s3Client)
			admin.POST("/uploads/generate-url", uploadHandler.GenerateUploadURL)
		}
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}