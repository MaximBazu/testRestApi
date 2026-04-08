package model

import "time"

type User struct {
	ID          int
	Name        string
	Surname     string
	Email       string
	TelegramTag string
	CreatedAt   time.Time
}

//HTTP → handler → service → repository → DB
//			↓
//			mapper → DTO → JSON
