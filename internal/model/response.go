package model

type Response struct {
	Payload interface{} `json:"payload"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}
