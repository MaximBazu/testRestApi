package dto

import "time"

type CreateUserRequest struct {
	Name        string
	Surname     string
	Email       string
	TelegramTag string
}

type UserResponse struct {
	ID          int
	Name        string
	Surname     string
	Email       string
	TelegramTag string
	CreatedAt   time.Time
}
