package reqmodel

import (
	"movie_backend_go/db/sqlc"
)

type UserRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type UserUpdateRequest struct {
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type UserListResponse struct {
	UserList []sqlc.UserDatum `json:"user_list"`
}
