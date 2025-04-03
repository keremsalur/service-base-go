package handler

import (
	"context"
	"errors"
	"service-base-go/pkg/logger"
	"service-base-go/pkg/otel"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
)

type Request any
type Response any

type HandlerInterface[Req Request, Res Response] interface {
	Handle(ctx context.Context, req *Req) (*Res, error)
}

//type Handler[Req Request, Res Response] func(ctx context.Context, req Req) (Res, error)

func Handle[Req Request, Res Response](handler HandlerInterface[Req, Res]) fiber.Handler {
	// Her isteği işleyen bir handler fonksiyonu oluşturuyoruz.
	// Bu fonksiyon, fiber.Ctx tipinde bir parametre alıyor ve error döndürüyor.
	return func(c *fiber.Ctx) error {

		// Tracing burda kopuyor c.Context ile yeni bir context oluştuğu için
		_, span := otel.GetTracer().Start(c.UserContext(), "Handle.Handler")
		defer span.End()

		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()

		var req Req
		logData := logger.GetGlobalLogData()
		logData["class"] = "Handler"

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			span.RecordError(err)
			span.SetStatus(codes.Error, "handler body parser error")
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Body Parser işlemi sırasında hata", logData)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ParamsParser(&req); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "handler params parser error")
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Param Parser işlemi sırasında hata", logData)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "handler query parser error")
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Query Parser işlemi sırasında hata", logData)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "handler req header parser error")
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Req Header Parser işlemi sırasında hata", logData)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		// middlewarelar
		// 1. request logger
		// 2. request validator
		// 3. request authorizer
		// 4. request handler
		//req = c.Context().Value("validatedDTO").(Req)

		res, err := handler.Handle(ctx, &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	}
}
