package models

import "gorm.io/gorm"

type Table struct {
	gorm.Model
	Name  string
	Email string
}
