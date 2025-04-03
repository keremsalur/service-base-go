package project

import (
	"context"
	"service-base-go/domain/model"

	"gorm.io/gorm"
)

type IProjectRepository interface {
	BeginTx() *gorm.DB
	Create(ctx context.Context, project *model.Project, tx *gorm.DB) error
	//Diğer repository işlemleri
}
