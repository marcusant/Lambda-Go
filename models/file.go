package models

import (
	"time"
)

type File struct {
	ID         uint      `db:"id"`
	Owner      uint      `db:"owner"`
	Name       string    `db:"name"`
	Extension  string    `db:"extension"`
	UploadDate time.Time `db:"upload_date"`
	Encrypted  bool      `db:"encrypted"`
	LocalName  string    `db:"local_name"`
}

// Inheriting from Model interface
func (f File) TableName() string {
	return "files"
}
