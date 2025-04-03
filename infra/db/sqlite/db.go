package sqlite

import (
	database "service-base-go/infra/db"
	"service-base-go/pkg/logger"

	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDatabase struct {
	DB *gorm.DB
}

func NewSqliteDatabase() database.Database {
	return &SqliteDatabase{}
}

func (d *SqliteDatabase) Connect(dsn string) database.Database {
	logData := logger.GetGlobalLogData()
	logData["class"] = "Sqlite Database"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Veritabanına bağlanılamadı.", logData)
	}
	logData["dsn"] = dsn
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Veritabanına bağlanıldı.", logData)
	return &SqliteDatabase{
		DB: db,
	}
}

func (d *SqliteDatabase) Migrate(models ...interface{}) {
	logData := logger.GetGlobalLogData()
	logData["class"] = "Sqlite Database"
	for _, model := range models {
		if err := d.DB.AutoMigrate(model); err != nil {
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.FatalLevel, "Veritabanı migrate işlemi başarısız.", logData)
		}
	}
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Veritabanı migrate işlemi başarıyla tamamlandı", logData)
}

func (d *SqliteDatabase) GetDB() *gorm.DB {
	return d.DB
}
