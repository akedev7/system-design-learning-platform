package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-course-service/internal/models"
)

type LessonRepository struct {
	db *sqlx.DB
}

func NewLessonRepository(db *sqlx.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) GetByID(id int) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Get(&lesson, "SELECT * FROM lessons WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lesson by id: %w", err)
	}
	return &lesson, nil
}

func (r *LessonRepository) GetByModuleID(moduleID int) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Select(&lessons, "SELECT * FROM lessons WHERE module_id = $1 ORDER BY order_index ASC", moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons by module id: %w", err)
	}
	return lessons, nil
}
