package database

import (
	"gorm.io/gorm"
)

type UploadFile struct {
	gorm.Model
	UserId      string
	FileName    string
	FileNameGen string
}
