package models

type RequestPayload struct {
	UserID  string `json:"user_id"`
	Payload string `json:"payload"`
}
