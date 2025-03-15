package response

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
}
