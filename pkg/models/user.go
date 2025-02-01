package models

import "time"

type User struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"not null;unique"`
	PasswordHash []byte   `json:"-" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

