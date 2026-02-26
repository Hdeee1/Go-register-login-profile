package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
)

type mySQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (domain.UserRepository, error) {
	return &mySQLUserRepository{db: db}, nil
}

func (m *mySQLUserRepository) Create(user *domain.User, ctx context.Context) error {
	query := "INSERT INTO users (full_name, username, email, password) VALUES (?, ?, ?, ?)"
	res, err := m.db.Exec(query, user.FullName, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.Id = int(id)	
	return nil
}

func (m *mySQLUserRepository) GetByEmail(user *domain.User, ctx context.Context) error {
	query := "SELECT id, full_name, username, email, password, created_at, updated_at FROM users WHERE email = ?"
	row := m.db.QueryRow(query, user.Email)
	
	if err := row.Scan(
			&user.Id, 
			&user.FullName, 
			&user.Username, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
		); err != nil {
		return err
	}
	
	return nil
}

func (m *mySQLUserRepository) GetById(id int) (*domain.User, error) {
	query := "SELECT id, full_name, username, email, password, created_at, updated_at FROM users WHERE id = ?"
	row := m.db.QueryRow(query, id)

	var user domain.User
	
	if err := row.Scan(
			&user.Id, 
			&user.FullName, 
			&user.Username, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
	); err != nil {
		return  nil, err
	}

	return &user, nil
}

func (m *mySQLUserRepository) FindByEmailOrUsername(email, username string) (*domain.User, error) {
	query := "SELECT id, full_name, username, email, password, created_at, updated_at FROM users WHERE email = ? OR username = ?"
	row := m.db.QueryRow(query, email, username)

	var user domain.User

	if err := row.Scan(
			&user.Id, 
			&user.FullName, 
			&user.Username, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
	); err != nil {
		return  nil, err
	}

	return &user, nil
}

func (m *mySQLUserRepository) Update(user *domain.User, ctx context.Context) error {
	fields := []string{}
	args := []any{}

	if user.Username != "" {
		fields = append(fields, "username = ?")
		args = append(args, user.Username)
	}

	if user.Password != "" {
		fields = append(fields, "password = ?")
		args = append(args, user.Password)
	}

	if len(fields) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, user.Id)
	query := "UPDATE users SET " + strings.Join(fields, ", ") + " WHERE id = ?"

	_, err := m.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
} 

func (m *mySQLUserRepository) SaveOTP(email, otp string, expiresAt time.Time, ctx context.Context) error {
	query := "INSERT INTO password_resets (email, otp, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE otp = ?, expires_at = ?"
	_, err := m.db.Exec(query, email, otp, expiresAt, otp, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *mySQLUserRepository) FindOTP(email string, ctx context.Context) (string, time.Time, error) {
	query := "SELECT otp, expires_at FROM password_resets WHERE email = ?"
	row := m.db.QueryRow(query, email)

	var otp string 
	var expires time.Time
	
	if err := row.Scan(
		&otp,
		&expires,
	); err != nil {
		return "", time.Time{}, err
	}

	return otp, expires, nil
}