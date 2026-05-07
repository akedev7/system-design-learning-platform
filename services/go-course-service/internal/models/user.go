package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Auth0ID   string    `db:"auth0_id" json:"auth0Id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}