package response

type ApiErrors struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
