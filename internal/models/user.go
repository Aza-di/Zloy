package models

import "time"

type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Login        string    `gorm:"unique;not null" json:"login"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Balance      int       `gorm:"not null;default:0" json:"balance"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
