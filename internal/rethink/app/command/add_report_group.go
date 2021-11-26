package command

import (
	"context"
	"wake_up_backend/internal/rethink/domain"
)

type AddReportGroup struct {
	UserID string
	Name   string
}

type AddReportGroupHandler struct {
	repo domain.IRepository
}

func NewAddReportGroupHandler(repo domain.IRepository) AddReportGroupHandler {
	return AddReportGroupHandler{repo: repo}
}

func (h AddReportGroupHandler) Handle(ctx context.Context, cmd AddReportGroup) error {
	group, err := domain.NewReportGroup(cmd.UserID, cmd.Name)
	if err != nil {
		return err
	}
	if _, err = h.repo.AddReportGroup(ctx, group); err != nil {
		return err
	}
	return nil
}
