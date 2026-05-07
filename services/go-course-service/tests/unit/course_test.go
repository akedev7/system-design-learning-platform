package unit

import (
	"testing"

	"go-course-service/internal/models"
)

func TestCourseModel(t *testing.T) {
	course := models.Course{
		ID:          1,
		Title:       "Test Course",
		Description: "Test Description",
	}

	if course.ID != 1 {
		t.Errorf("expected ID 1, got %d", course.ID)
	}
	if course.Title != "Test Course" {
		t.Errorf("expected Title 'Test Course', got '%s'", course.Title)
	}
	if course.Description != "Test Description" {
		t.Errorf("expected Description 'Test Description', got '%s'", course.Description)
	}
}