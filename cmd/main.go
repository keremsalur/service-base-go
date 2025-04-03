package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"service-base-go/pkg/config"
	"service-base-go/pkg/logger"
	"service-base-go/pkg/otel"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog"
)

func main() {

	//configs
	cfg := config.LoadConfig()

	//logger
	logger.InitLogger(cfg.AppName)

	// db (sqlite)
	/*
		db := sqlite.NewSqliteDatabase().Connect(cfg.DatabaseUrl)
		db.Migrate() // models interface{}
		db.GetDB().Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	*/

	// Zipkin OpenTelemetry
	// OpenTelemetry provider'ı başlat
	tp, err := otel.InitTracer("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,

		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,

		Concurrency: 256 * 1024,
	})

	app.Use(otelfiber.Middleware(
		otelfiber.WithServerName(cfg.AppName),
	))

	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Uygulama başlatılıyor", nil)

	// metrikler, loglar, hata takibi, vs.
	// Prometheus, Grafana
	//app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/metrics", monitor.New())

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			logger.GetLogger().Fatal().Err(err).Msg("Uygulama başlatılamadı")
			os.Exit(1)
		}
	}()
	greacefulShutdown(app)
}

func greacefulShutdown(app *fiber.App) {
	// Greaceful shutdown
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Uygulama başlatıldı", nil)
	<-signChan
	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Uygulama kapatılıyor", nil)

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Uygulama kapatılırken hata oluştu", map[string]interface{}{"error": err})
	}

	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Uygulama kapatıldı", nil)
}
