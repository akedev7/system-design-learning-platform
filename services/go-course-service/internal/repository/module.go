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

func (r *ModuleRepository) Create(module *models.Module) (*models.Module, error) {
	var newModule models.Module
	err := r.db.Get(&newModule,
		`INSERT INTO modules (course_id, title, description, order_index)
		 VALUES ($1, $2, $3, $4) RETURNING id, course_id, title, description, order_index, created_at, updated_at`,
		module.CourseID, module.Title, module.Description, module.OrderIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}
	return &newModule, nil
}

func (r *ModuleRepository) Update(module *models.Module) (*models.Module, error) {
	var updatedModule models.Module
	err := r.db.Get(&updatedModule,
		`UPDATE modules SET title = $1, description = $2, order_index = $3, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $4 RETURNING id, course_id, title, description, order_index, created_at, updated_at`,
		module.Title, module.Description, module.OrderIndex, module.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update module: %w", err)
	}
	return &updatedModule, nil
}

func (r *ModuleRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM modules WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
