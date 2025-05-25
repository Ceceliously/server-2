package response


type Response struct {
    Status string `json:"status"`
    Error  string `json:"error,omitempty"`
}

const (
	StatusOk = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOk,
	} 
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}

type User struct {
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Age int `json:"age"`
}