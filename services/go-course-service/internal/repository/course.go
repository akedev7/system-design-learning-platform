package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-course-service/internal/models"
)

type CourseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(db *sqlx.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Select(&courses, "SELECT * FROM courses ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

func (r *CourseRepository) GetByID(id int) (*models.Course, error) {
	var course models.Course
	err := r.db.Get(&course, "SELECT * FROM courses WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get course by id: %w", err)
	}
	return &course, nil
}

func (r *CourseRepository) Create(course *models.Course) (*models.Course, error) {
	var newCourse models.Course
	err := r.db.Get(&newCourse,
		`INSERT INTO courses (title, description) VALUES ($1, $2)
		 RETURNING id, title, description, created_at, updated_at`,
		course.Title, course.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}
	return &newCourse, nil
}

func (r *CourseRepository) Update(course *models.Course) (*models.Course, error) {
	var updatedCourse models.Course
	err := r.db.Get(&updatedCourse,
		`UPDATE courses SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $3 RETURNING id, title, description, created_at, updated_at`,
		course.Title, course.Description, course.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update course: %w", err)
	}
	return &updatedCourse, nil
}

func (r *CourseRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}