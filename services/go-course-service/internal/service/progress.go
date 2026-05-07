package service

import (
	"go-course-service/internal/models"
)

type ProgressTrackingService struct{}

func NewProgressTrackingService() *ProgressTrackingService {
	return &ProgressTrackingService{}
}

func (s *ProgressTrackingService) CalculateModuleCompletion(lessonProgress []models.LessonProgress) models.ModuleProgress {
	if len(lessonProgress) == 0 {
		return models.ModuleProgress{
			CompletionPercentage: 0,
			CompletedLessons:      0,
			TotalLessons:          0,
		}
	}

	completedLessons := 0
	for _, lp := range lessonProgress {
		if lp.Passed {
			completedLessons++
		}
	}

	totalLessons := len(lessonProgress)
	percentage := (completedLessons * 100) / totalLessons

	return models.ModuleProgress{
		ModuleID:              0,
		TotalLessons:          totalLessons,
		CompletedLessons:      completedLessons,
		CompletionPercentage:  percentage,
	}
}

func (s *ProgressTrackingService) CalculateCourseCompletion(moduleProgress []models.ModuleProgress) models.CourseProgress {
	if len(moduleProgress) == 0 {
		return models.CourseProgress{
			CompletionPercentage: 0,
			CompletedModules:     0,
			TotalModules:         0,
		}
	}

	completedModules := 0
	totalLessons := 0
	completedLessonsTotal := 0

	for _, mp := range moduleProgress {
		if mp.CompletionPercentage == 100 {
			completedModules++
		}
		totalLessons += mp.TotalLessons
		completedLessonsTotal += mp.CompletedLessons
	}

	totalModules := len(moduleProgress)
	percentage := 0
	if totalLessons > 0 {
		percentage = (completedLessonsTotal * 100) / totalLessons
	}

	return models.CourseProgress{
		CourseID:              0,
		TotalModules:          totalModules,
		CompletedModules:      completedModules,
		TotalLessons:          totalLessons,
		CompletedLessons:      completedLessonsTotal,
		CompletionPercentage:  percentage,
	}
}

func (s *ProgressTrackingService) GetLastAccessedLesson(lessonProgress []models.LessonProgress) models.LessonProgress {
	if len(lessonProgress) == 0 {
		return models.LessonProgress{LessonID: 0}
	}

	for _, lp := range lessonProgress {
		if !lp.Passed {
			return lp
		}
	}

	return lessonProgress[len(lessonProgress)-1]
}