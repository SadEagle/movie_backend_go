package models

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UserData) ToResponse() UserDataResponse {
	return UserDataResponse{
		ID:        u.ID.String(),
		Login:     u.Login,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

// Request structs
type CreateUserDataRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateUserDataRequest struct {
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

// Response structs
type UserDataResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
