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

func (r *LessonRepository) Create(lesson *models.Lesson) (*models.Lesson, error) {
	var newLesson models.Lesson
	err := r.db.Get(&newLesson,
		`INSERT INTO lessons (module_id, title, description, content_jsonb, order_index)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, module_id, title, description, content_jsonb, order_index, created_at, updated_at`,
		lesson.ModuleID, lesson.Title, lesson.Description, lesson.ContentJSON, lesson.OrderIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}
	return &newLesson, nil
}

func (r *LessonRepository) Update(lesson *models.Lesson) (*models.Lesson, error) {
	var updatedLesson models.Lesson
	err := r.db.Get(&updatedLesson,
		`UPDATE lessons SET title = $1, description = $2, content_jsonb = $3, order_index = $4, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $5 RETURNING id, module_id, title, description, content_jsonb, order_index, created_at, updated_at`,
		lesson.Title, lesson.Description, lesson.ContentJSON, lesson.OrderIndex, lesson.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update lesson: %w", err)
	}
	return &updatedLesson, nil
}

func (r *LessonRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM lessons WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *LessonRepository) UpdateContent(id int, contentJSON []byte) (*models.Lesson, error) {
	var updatedLesson models.Lesson
	err := r.db.Get(&updatedLesson,
		`UPDATE lessons SET content_jsonb = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2 RETURNING id, module_id, title, description, content_jsonb, order_index, created_at, updated_at`,
		contentJSON, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update lesson content: %w", err)
	}
	return &updatedLesson, nil
}
