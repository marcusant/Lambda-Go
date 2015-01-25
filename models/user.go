package models

import (
	"time"
)

type User struct {
	ID                uint      `db:"id"`
	Username          string    `db:"username"`
	Password          string    `db:"password"`
	CreationDate      time.Time `db:"creation_date"`
	ApiKey            string    `db:"apikey"`
	EncryptionEnabled bool      `db:"encryption_enabled"`
	ThemeName         string    `db:"theme_name"`
}

// Inheriting from Model interface
func (u User) TableName() string {
	return "users"
}
