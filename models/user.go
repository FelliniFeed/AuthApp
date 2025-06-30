package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID    `json:"Id" db:"id"`
	Username string ` json:"UserName" db:"username"`
	Password string ` json:"Password" db:"password"`
}