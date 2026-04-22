package model

import "time"

type User struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Email       string
	Phone       string
	TelegramTag string
	CreatedAt   time.Time
}

//HTTP → handler → service → repository → DB
//			↓
//			mapper → DTO → JSON
