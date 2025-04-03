package repository

import (
	"context"

	"gorm.io/gorm"
)

type IBaseRepository[T any, ID any] interface {
	BeginTx() *gorm.DB
	Create(ctx context.Context, entity *T, tx *gorm.DB) error
	GetAll(ctx context.Context, entities *[]T) error
	GetByID(ctx context.Context, entity *T, id ID, preload ...string) error
	Update(ctx context.Context, entity *T, tx *gorm.DB) error
	Delete(ctx context.Context, entity *T, id ID, tx *gorm.DB) error
}
