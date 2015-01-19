package models

import (
	"time"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	CreationDate time.Time `db:"creation_date"`
}

// Inheriting from Model interface
func (u User) TableName() string {
	return "users"
}
