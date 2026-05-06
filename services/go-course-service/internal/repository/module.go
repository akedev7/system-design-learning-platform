package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-course-service/internal/models"
)

type ModuleRepository struct {
	db *sqlx.DB
}

func NewModuleRepository(db *sqlx.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) GetByID(id int) (*models.Module, error) {
	var module models.Module
	err := r.db.Get(&module, "SELECT * FROM modules WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get module by id: %w", err)
	}
	return &module, nil
}

func (r *ModuleRepository) GetByCourseID(courseID int) ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Select(&modules, "SELECT * FROM modules WHERE course_id = $1 ORDER BY order_index ASC", courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules by course id: %w", err)
	}
	return modules, nil
}
