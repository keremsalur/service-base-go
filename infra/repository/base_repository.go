package repository

import (
	"context"
	"reflect"
	"service-base-go/domain/repository"
	"service-base-go/pkg/logger"
	"service-base-go/pkg/otel"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

type BaseRepository[T any, ID any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any, ID any](db *gorm.DB) repository.IBaseRepository[T, ID] {
	return &BaseRepository[T, ID]{
		db: db,
	}
}

func (r *BaseRepository[T, ID]) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *BaseRepository[T, ID]) Create(ctx context.Context, entity *T, tx *gorm.DB) error {

	// Tracing
	ctx, span := otel.GetTracer().Start(ctx, "BaseRepository.Create")
	defer span.End()

	logData := logger.GetGlobalLogData()
	logData["class"] = "BaseRepository"
	logData["request_id"] = ctx.Value("X-Request-ID").(string)
	entityName := reflect.TypeOf(entity).Elem().Name()
	err := tx.Save(entity).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "database error")
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Yeni "+entityName+" veritabanına kaydedilirken hata.", logData)
		return err
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Yeni "+entityName+" veritabanına kaydedildi.", logData)
	return err
}

func (r *BaseRepository[T, ID]) GetAll(ctx context.Context, entities *[]T) error {
	logData := logger.GetGlobalLogData()
	logData["class"] = "BaseRepository"
	logData["request_id"] = ctx.Value("X-Request-ID").(string)
	entityName := reflect.TypeOf(entities).Elem().Name()
	err := r.db.Find(entities).Error
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, ""+entityName+"(lar/ler) alınırken hata", logData)
		return err
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, ""+entityName+"(lar/ler) alındı.", logData)
	return err
}

func (r *BaseRepository[T, ID]) GetByID(ctx context.Context, entity *T, id ID, preloads ...string) error {
	logData := logger.GetGlobalLogData()
	logData["class"] = "BaseRepository"
	logData["request_id"] = ctx.Value("X-Request-ID").(string)
	entityName := reflect.TypeOf(entity).Elem().Name()

	// Preload'ları ekle
	for _, preload := range preloads {
		r.db = r.db.Preload(preload)
	}

	err := r.db.Where("id = ?", id).First(entity).Error
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, ""+entityName+" bulunamadı.", logData)
		return err
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, ""+entityName+" bulundu.", logData)
	return err
}

func (r *BaseRepository[T, ID]) Delete(ctx context.Context, entity *T, id ID, tx *gorm.DB) error {
	logData := logger.GetGlobalLogData()
	logData["class"] = "BaseRepository"
	logData["request_id"] = ctx.Value("X-Request-ID").(string)
	entityName := reflect.TypeOf(entity).Elem().Name()
	err := tx.Delete(entity, id).Error
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, ""+entityName+" veritabanından silinirken hata.", logData)
		return err
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, ""+entityName+" veritabanından silindi.", logData)
	return err
}

func (r *BaseRepository[T, ID]) Update(ctx context.Context, entity *T, tx *gorm.DB) error {
	logData := logger.GetGlobalLogData()
	logData["class"] = "BaseRepository"
	logData["request_id"] = ctx.Value("X-Request-ID").(string)
	entityName := reflect.TypeOf(entity).Elem().Name()
	err := tx.Save(entity).Error
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Yeni "+entityName+" veritabanında güncellenirken hata.", logData)
		return err
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Yeni "+entityName+" veritabanında güncellendi.", logData)
	return err
}
