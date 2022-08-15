package models

import (
	"gorm.io/gorm"
)

type JwtBlacklist struct {
	gorm.Model
	UserId string
	Token  string
}
