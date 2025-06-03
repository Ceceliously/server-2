package createuser

import (
	res "server-2/internal/models/response"
)

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FirstName *string `json:"first_name,omitempty"`
	LastName *string `json:"last_name,omitempty"`
	Age *int `json:"age,omitempty"`
}

type UserCreateResponse struct {
	res.Response
	Username string `json:"username,omitempty"`
}

func ToResponse(username string) *UserCreateResponse {
	return &UserCreateResponse {
		Response: res.OK(),
		Username: username,
	}
}