package response

type ApiErrors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
