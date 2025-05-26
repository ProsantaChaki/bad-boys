package repository

import (
	"bad_boyes/internal/config"
	"bad_boyes/internal/model"
	"time"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	_, err := config.DB.Exec(
		"INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, time.Now(), time.Now(),
	)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := config.DB.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserExists(email string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	return exists, err
}
