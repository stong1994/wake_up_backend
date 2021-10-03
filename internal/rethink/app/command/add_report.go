package command

import (
	"context"
	"wake_up_backend/internal/rethink/domain"
)

type AddReport struct {
	ReportID string
	UserID string
	GroupID string
	Content string
}

type AddReportHandler struct {
	repo domain.IRepository
}

func (h AddReportHandler) Handle(ctx context.Context, cmd AddReport) error {
	report, err := domain.NewReport(cmd.ReportID, cmd.GroupID, cmd.UserID, cmd.Content)
	if err != nil {
		return err
	}
	if err = h.repo.AddReport(ctx, report); err != nil {
		return err
	}
	return nil
}