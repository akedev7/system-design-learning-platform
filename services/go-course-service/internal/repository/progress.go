package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go-course-service/internal/models"
)

type ProgressRepository struct {
	db *sqlx.DB
}

func NewProgressRepository(db *sqlx.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) UpsertProgress(userID string, lessonID int, score int, passed bool) (*models.UserLessonProgress, error) {
	var progress models.UserLessonProgress

	query := `
		INSERT INTO user_lesson_progress (user_id, lesson_id, score, passed, completed_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE SET score = $3, passed = $4, completed_at = $5, updated_at = $5
		RETURNING id, user_id, lesson_id, score, passed, completed_at, created_at, updated_at
	`

	var completedAt sql.NullTime
	if passed {
		completedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	err := r.db.Get(&progress, query, userID, lessonID, score, passed, completedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert progress: %w", err)
	}

	return &progress, nil
}

func (r *ProgressRepository) GetProgress(userID string, lessonID int) (*models.UserLessonProgress, error) {
	var progress models.UserLessonProgress

	query := `
		SELECT id, user_id, lesson_id, score, passed, completed_at, created_at, updated_at
		FROM user_lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	err := r.db.Get(&progress, query, userID, lessonID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	return &progress, nil
}