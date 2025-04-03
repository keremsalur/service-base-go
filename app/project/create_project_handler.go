package project

import (
	"context"
	"service-base-go/domain/model"
	"service-base-go/pkg/logger"
	"service-base-go/pkg/otel"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
)

type CreateProjectRequest struct {
	ProjectName string `json:"name" validate:"required"`
}

type CreateProjectResponse struct {
	ID uint `json:"id"`
}

type CreateProjectHandler struct {
	repository IProjectRepository
}

func NewCreateProductHandler(repository IProjectRepository) *CreateProjectHandler {
	return &CreateProjectHandler{
		repository: repository,
	}
}

func (h *CreateProjectHandler) Handle(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {

	// Tracing
	ctx, span := otel.GetTracer().Start(ctx, "Handle.CreateProjectHandler")
	defer span.End()

	// Spanlarda değer geçebilmek için
	//span.AddEvent("Deneme logic")

	logData := logger.GetGlobalLogData()
	logData["class"] = "Create Project Handler"
	project := &model.Project{
		ProjectName: req.ProjectName,
	}

	tx := h.repository.BeginTx()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	defer tx.Rollback()

	if err := h.repository.Create(ctx, project, tx); err != nil {
		tx.Rollback()
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Kayıt edilemedi.", logData)
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository error")
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "İşlem commit edilemedi.", logData)
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository error")
		return nil, err
	}

	return &CreateProjectResponse{ID: project.ID}, nil

}
