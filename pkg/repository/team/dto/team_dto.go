package dto

type TeamsByUserResponse struct {
	Teams []TeamResponse `json:"teams"`
}

type TeamResponse struct {
	Code string `json:"code"`
	Debt int    `json:"debt"`
}
