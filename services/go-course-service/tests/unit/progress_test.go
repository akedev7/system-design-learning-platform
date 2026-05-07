package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-course-service/internal/models"
	"go-course-service/internal/service"
)

func TestProgressTrackingService_CalculateModuleCompletion_AllComplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{
		{LessonID: 1, Passed: true},
		{LessonID: 2, Passed: true},
	}

	result := pts.CalculateModuleCompletion(lessonProgress)

	assert.Equal(t, 100, result.CompletionPercentage)
	assert.Equal(t, 2, result.CompletedLessons)
	assert.Equal(t, 2, result.TotalLessons)
}

func TestProgressTrackingService_CalculateModuleCompletion_PartialComplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{
		{LessonID: 1, Passed: true},
		{LessonID: 2, Passed: false},
		{LessonID: 3, Passed: true},
	}

	result := pts.CalculateModuleCompletion(lessonProgress)

	assert.Equal(t, 66, result.CompletionPercentage)
	assert.Equal(t, 2, result.CompletedLessons)
	assert.Equal(t, 3, result.TotalLessons)
}

func TestProgressTrackingService_CalculateModuleCompletion_NoneComplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{
		{LessonID: 1, Passed: false},
		{LessonID: 2, Passed: false},
	}

	result := pts.CalculateModuleCompletion(lessonProgress)

	assert.Equal(t, 0, result.CompletionPercentage)
	assert.Equal(t, 0, result.CompletedLessons)
	assert.Equal(t, 2, result.TotalLessons)
}

func TestProgressTrackingService_CalculateModuleCompletion_Empty(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{}

	result := pts.CalculateModuleCompletion(lessonProgress)

	assert.Equal(t, 0, result.CompletionPercentage)
	assert.Equal(t, 0, result.CompletedLessons)
	assert.Equal(t, 0, result.TotalLessons)
}

func TestProgressTrackingService_CalculateCourseCompletion_AllModulesComplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	moduleProgress := []models.ModuleProgress{
		{ModuleID: 1, TotalLessons: 2, CompletedLessons: 2, CompletionPercentage: 100},
		{ModuleID: 2, TotalLessons: 3, CompletedLessons: 3, CompletionPercentage: 100},
		{ModuleID: 3, TotalLessons: 2, CompletedLessons: 2, CompletionPercentage: 100},
	}

	result := pts.CalculateCourseCompletion(moduleProgress)

	assert.Equal(t, 100, result.CompletionPercentage)
	assert.Equal(t, 3, result.CompletedModules)
	assert.Equal(t, 3, result.TotalModules)
}

func TestProgressTrackingService_CalculateCourseCompletion_PartialModules(t *testing.T) {
	pts := service.NewProgressTrackingService()

	moduleProgress := []models.ModuleProgress{
		{ModuleID: 1, TotalLessons: 2, CompletedLessons: 2, CompletionPercentage: 100},
		{ModuleID: 2, TotalLessons: 3, CompletedLessons: 1, CompletionPercentage: 33},
		{ModuleID: 3, TotalLessons: 2, CompletedLessons: 0, CompletionPercentage: 0},
	}

	result := pts.CalculateCourseCompletion(moduleProgress)

	assert.Equal(t, 42, result.CompletionPercentage)
	assert.Equal(t, 1, result.CompletedModules)
	assert.Equal(t, 3, result.TotalModules)
}

func TestProgressTrackingService_CalculateCourseCompletion_Empty(t *testing.T) {
	pts := service.NewProgressTrackingService()

	moduleProgress := []models.ModuleProgress{}

	result := pts.CalculateCourseCompletion(moduleProgress)

	assert.Equal(t, 0, result.CompletionPercentage)
	assert.Equal(t, 0, result.CompletedModules)
	assert.Equal(t, 0, result.TotalModules)
}

func TestProgressTrackingService_GetLastAccessedLesson_FirstIncomplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{
		{LessonID: 1, Passed: false},
		{LessonID: 2, Passed: true},
		{LessonID: 3, Passed: true},
	}

	result := pts.GetLastAccessedLesson(lessonProgress)

	assert.Equal(t, 1, result.LessonID)
}

func TestProgressTrackingService_GetLastAccessedLesson_AllComplete(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{
		{LessonID: 1, Passed: true},
		{LessonID: 2, Passed: true},
		{LessonID: 3, Passed: true},
	}

	result := pts.GetLastAccessedLesson(lessonProgress)

	assert.Equal(t, 3, result.LessonID)
}

func TestProgressTrackingService_GetLastAccessedLesson_Empty(t *testing.T) {
	pts := service.NewProgressTrackingService()

	lessonProgress := []models.LessonProgress{}

	result := pts.GetLastAccessedLesson(lessonProgress)

	assert.Equal(t, 0, result.LessonID)
}