package model

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ProjectName string
}
