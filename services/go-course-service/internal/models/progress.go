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