package database

import (
	"gorm.io/gorm"
	
)

type Database interface {
	ConnectionGetting() *gorm.DB
}