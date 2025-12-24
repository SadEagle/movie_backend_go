package model

import "github.com/google/uuid"

type UserBase struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID uuid.UUID `json:"id"`
	UserBase
}

type UserCreateRequest struct {
	UserBase
}
type UserUpdateRequest struct {
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type UserResponse struct {
	User
}
