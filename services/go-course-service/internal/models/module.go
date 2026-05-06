package models

import (
	"time"
)

type Module struct {
	ID          int        `db:"id" json:"id"`
	CourseID    int        `db:"course_id" json:"courseId"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	OrderIndex  int        `db:"order_index" json:"orderIndex"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}
