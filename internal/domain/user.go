package domain

import (
	"context"
	"time"
)

type User struct {
	Id			int			`json:"id" `
	FullName	string		`json:"full_name"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
	CreatedAt	time.Time	`json:"created_at" `
	UpdatedAt	time.Time	`json:"updated_at" `
}

type RegisterRequest struct {
	FullName	string `json:"full_name" binding:"required"`
	Username	string `json:"username" binding:"required,min=3"`
	Email		string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email		 string `json:"email" binding:"required,email"`
	Password	 string `json:"password" binding:"required,min=8"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	Username	string	`json:"username" binding:"required,min=3"`
	Password	string	`json:"password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserRepository interface {
	Create(user *User, ctx context.Context) error
	GetByEmail(user *User, ctx context.Context) error
	GetById(id int) (*User, error)
	FindByEmailOrUsername(email, username string) (*User, error)
	Update(user *User, ctx context.Context) error
}

type UserUsecase interface {
	Register(user RegisterRequest, ctx context.Context) (*User, error)
	Login(user LoginRequest, ctx context.Context) (*User, string, string, error)
	GetProfile(userId int, ctx context.Context) (*User, error)
	Refresh(input RefreshTokenRequest, ctx context.Context) (string, error)
	UpdateProfile(userId int, input UpdateProfileRequest, ctx context.Context) (*User, error)
}