package request

type UserRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
}
