package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	IsAdmin   *bool     `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type UserList struct {
	UserList []User `json:"user_list"`
}

// Request structs
type CreateUserRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}
