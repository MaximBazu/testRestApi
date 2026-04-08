package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToUserResponse(u *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Surname:     u.Surname,
		Email:       u.Email,
		TelegramTag: u.TelegramTag,
	}
}
