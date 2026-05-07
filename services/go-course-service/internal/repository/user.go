package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-course-service/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByAuth0ID(auth0ID string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE auth0_id = $1", auth0ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by auth0_id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	var newUser models.User
	err := r.db.Get(&newUser,
		`INSERT INTO users (auth0_id, email, name, role) VALUES ($1, $2, $3, $4)
		 RETURNING id, auth0_id, email, name, role, created_at, updated_at`,
		user.Auth0ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &newUser, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}