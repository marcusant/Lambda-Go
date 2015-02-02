package models

import (
	"time"
)

type Paste struct {
	ID          uint      `db:"id"`
	Owner       uint      `db:"owner"`
	Name        string    `db:"name"`
	UploadDate  time.Time `db:"upload_date"`
	ContentJson string    `db:"content_json"`
	IsCode      bool      `db:"is_code"`
}

// Inheriting from Model interface
func (f Paste) TableName() string {
	return "pastes"
}
