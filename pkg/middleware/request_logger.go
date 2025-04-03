package middleware

import (
	"service-base-go/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestLogger(c *fiber.Ctx) error {
	logData := logger.GetGlobalLogData()
	logData["method"] = c.Method()
	logData["url"] = c.BaseURL() + c.OriginalURL()
	logData["ip"] = c.IP()
	logData["status"] = c.Response().StatusCode()

	log := logger.GetLogger()
	start := time.Now()

	traceID := c.Get("X-Request-ID")
	if traceID == "" {
		traceID = uuid.NewString()
		c.Set("X-Request-ID", traceID)
	}
	c.Locals("X-Request-ID", traceID)
	logData["request_id"] = traceID

	userID := c.Get("X-User-ID")
	if userID != "" {
		c.Locals("userID", userID)
		logData["userID"] = userID
	}

	logger.PushLog(log, zerolog.InfoLevel, "İstek geldi.", logData)

	err := c.Next()

	duration := time.Since(start)
	logData["duration"] = duration
	logData["method"] = c.Method()
	logData["url"] = c.BaseURL() + c.OriginalURL()
	logData["ip"] = c.IP()
	logData["status"] = c.Response().StatusCode()

	if c.Response().StatusCode() >= 400 {
		logData["error"] = string(c.Response().Body())
		logger.PushLog(log, zerolog.ErrorLevel, "İstek hata ile sonuçlandı.", logData)
	} else {

		logger.PushLog(log, zerolog.InfoLevel, "İstek tamamlandı.", logData)
	}

	if err != nil {
		log.Error().Err(err).Msg("Hata meydana geldi")
	}

	return err
}
