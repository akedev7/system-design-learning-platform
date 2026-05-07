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

func (r *ProgressRepository) GetLessonProgressByUser(userID string, lessonIDs []int) ([]models.UserLessonProgress, error) {
	if len(lessonIDs) == 0 {
		return []models.UserLessonProgress{}, nil
	}

	var progress []models.UserLessonProgress
	query := `
		SELECT id, user_id, lesson_id, score, passed, completed_at, created_at, updated_at
		FROM user_lesson_progress
		WHERE user_id = $1 AND lesson_id = ANY($2)
	`

	err := r.db.Select(&progress, query, userID, lessonIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson progress: %w", err)
	}

	return progress, nil
}

func (r *ProgressRepository) UpsertCourseProgress(userID string, courseID int, percentage int) (*models.UserCourseProgress, error) {
	var progress models.UserCourseProgress

	query := `
		INSERT INTO user_course_progress (user_id, course_id, completion_percentage, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, course_id)
		DO UPDATE SET completion_percentage = $3, updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, course_id, completion_percentage, completed_at, created_at, updated_at
	`

	err := r.db.Get(&progress, query, userID, courseID, percentage)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert course progress: %w", err)
	}

	return &progress, nil
}

func (r *ProgressRepository) GetCourseProgress(userID string, courseID int) (*models.UserCourseProgress, error) {
	var progress models.UserCourseProgress

	query := `
		SELECT id, user_id, course_id, completion_percentage, completed_at, created_at, updated_at
		FROM user_course_progress
		WHERE user_id = $1 AND course_id = $2
	`

	err := r.db.Get(&progress, query, userID, courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get course progress: %w", err)
	}

	return &progress, nil
}

func (r *ProgressRepository) GetCourseProgressByUser(userID string, courseIDs []int) ([]models.UserCourseProgress, error) {
	if len(courseIDs) == 0 {
		return []models.UserCourseProgress{}, nil
	}

	var progress []models.UserCourseProgress
	query := `
		SELECT id, user_id, course_id, completion_percentage, completed_at, created_at, updated_at
		FROM user_course_progress
		WHERE user_id = $1 AND course_id = ANY($2)
	`

	err := r.db.Select(&progress, query, userID, courseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get course progress: %w", err)
	}

	return progress, nil
}