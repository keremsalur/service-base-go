package repository

import (
	"service-base-go/domain/model"
	"service-base-go/domain/repository"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	*BaseRepository[model.Project, uint]
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) repository.IProjectRepository[model.Project, uint] {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository[model.Project, uint](db).(*BaseRepository[model.Project, uint]),
		db:             db,
	}
}
