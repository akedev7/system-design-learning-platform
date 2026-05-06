package models

import (
	"encoding/json"
	"time"
)

type Lesson struct {
	ID          int             `db:"id" json:"id"`
	ModuleID    int             `db:"module_id" json:"moduleId"`
	Title       string          `db:"title" json:"title"`
	Description string          `db:"description" json:"description"`
	ContentJSON json.RawMessage `db:"content_jsonb" json:"contentJsonb"`
	OrderIndex  int             `db:"order_index" json:"orderIndex"`
	CreatedAt   time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updatedAt"`
}
