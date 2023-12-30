package response

type UserResponse struct {
	Name           string `json:"name"`
	SecondName     string `json:"second_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	NickName       string `json:"nick_name"`
}
