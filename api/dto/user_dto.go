package dto

import "building-distributed-app-in-gin-chapter06/api/models"

type UserDto struct {
	Username  string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user models.User) UserDto {
	return UserDto{
		Username:  user.Username,
		Telephone: user.Telephone,
	}
}
