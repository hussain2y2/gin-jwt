package models

import "time"

// User ...
type User struct {
	ID              uint64    `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	CreatedAT       time.Time `json:"created_at"`
	UpdatedAT       time.Time `json:"updated_at"`
}
