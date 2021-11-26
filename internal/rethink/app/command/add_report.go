package command

import (
	"context"
	"wake_up_backend/internal/rethink/domain"
)

type AddReport struct {
	UserID  string
	GroupID string
}

type AddReportHandler struct {
	repo domain.IRepository
}

func NewAddReportHandler(repo domain.IRepository) AddReportHandler {
	return AddReportHandler{repo: repo}
}

func (h AddReportHandler) Handle(ctx context.Context, cmd AddReport) error {
	report, err := domain.NewReport(ctx, cmd.UserID, cmd.GroupID, h.repo.CheckGroup)
	if err != nil {
		return err
	}
	if _, err = h.repo.AddReport(ctx, report); err != nil {
		return err
	}
	return nil
}
