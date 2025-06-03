package userget

import (
	res "server-2/internal/models/response"
	entity "server-2/internal/models/user/user"
)


type UserGetResponse struct {
	res.Response
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Age int `json:"age"`
}

func ToResponse(obtU *entity.User) *UserGetResponse {
	resp := UserGetResponse {
		Response: res.OK(),
		Username: obtU.Username,
	}

	if obtU.FirstName != nil {
		resp.FirstName = *obtU.FirstName
	}

	if obtU.LastName != nil {
		resp.LastName = *obtU.LastName
	}

	if obtU.Age != nil {
		resp.Age = *obtU.Age
	}

	return &resp
}


	