package models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Name   string
	Path   string
	Status bool
}
