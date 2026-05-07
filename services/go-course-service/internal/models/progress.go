package models

import (
	"time"
)

type UserLessonProgress struct {
	ID          int        `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"userId"`
	LessonID    int        `db:"lesson_id" json:"lessonId"`
	Score       int        `db:"score" json:"score"`
	Passed      bool       `db:"passed" json:"passed"`
	CompletedAt *time.Time `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}

type UserCourseProgress struct {
	ID                    int        `db:"id" json:"id"`
	UserID                string     `db:"user_id" json:"userId"`
	CourseID              int        `db:"course_id" json:"courseId"`
	CompletionPercentage  int        `db:"completion_percentage" json:"completionPercentage"`
	CompletedAt           *time.Time `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updatedAt"`
}

type CourseProgress struct {
	CourseID             int  `json:"courseId"`
	TotalModules         int  `json:"totalModules"`
	CompletedModules     int  `json:"completedModules"`
	TotalLessons         int  `json:"totalLessons"`
	CompletedLessons     int  `json:"completedLessons"`
	CompletionPercentage int  `json:"completionPercentage"`
}

type ModuleProgress struct {
	ModuleID              int  `json:"moduleId"`
	TotalLessons          int  `json:"totalLessons"`
	CompletedLessons      int  `json:"completedLessons"`
	CompletionPercentage  int  `json:"completionPercentage"`
}

type LessonProgress struct {
	LessonID    int         `json:"lessonId"`
	Score       int         `json:"score"`
	Passed      bool        `json:"passed"`
	CompletedAt *time.Time  `json:"completedAt,omitempty"`
}