package dto

import "time"

type CreateUserRequest struct {
	Name        string
	Surname     string
	Patronymic  string
	Email       string
	Phone       string
	TelegramTag string
}

type UserResponse struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Email       string
	Phone       string
	TelegramTag string
	CreatedAt   time.Time
}
