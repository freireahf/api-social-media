package models

type Password struct {
	NewPassword     string `json:"new-password"`
	CurrentPassword string `json:"current-password"`
}
