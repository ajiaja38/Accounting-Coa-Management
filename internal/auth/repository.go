package auth

import (
	"fiber.com/session-api/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByID(id uuid.UUID) (*domain.User, error)
	CreateUser(user *domain.User) error
	EmailExists(email string) (bool, error)
	UserNameExists(userName string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Raw(
		`SELECT id, user_name, email, password, role, created_at, updated_at
		 FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1`,
		email,
	).Scan(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

func (r *repository) FindUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Raw(
		`SELECT id, user_name, email, role, created_at, updated_at
		 FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1`,
		id,
	).Scan(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *repository) CreateUser(user *domain.User) error {
	return r.db.Exec(
		`INSERT INTO users (id, user_name, email, password, role, created_at, updated_at)
		 VALUES (gen_random_uuid(), ?, ?, ?, ?, NOW(), NOW())`,
		user.UserName, user.Email, user.Password, user.Role,
	).Error
}

func (r *repository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Raw(
		`SELECT COUNT(1) FROM users WHERE email = ? AND deleted_at IS NULL`, email,
	).Scan(&count).Error
	return count > 0, err
}

func (r *repository) UserNameExists(userName string) (bool, error) {
	var count int64
	err := r.db.Raw(
		`SELECT COUNT(1) FROM users WHERE user_name = ? AND deleted_at IS NULL`, userName,
	).Scan(&count).Error
	return count > 0, err
}
