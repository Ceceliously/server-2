package user

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	Age *int `json:"age"`
}


// type UserUseCase interface {
// 	CreateUser(user *User) (error)
// 	GetUser(user *User) (*User, error)
// }