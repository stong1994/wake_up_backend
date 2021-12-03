package command

import (
	"context"
	"wake_up_backend/internal/rethink/domain"
)

type CompleteRethink struct {
	RethinkID      string
	UserID         string
	ReportContent  string
	RethinkContent string
}

type CompleteRethinkHandler struct {
	repo domain.IRepository
}

func NewCompleteRethinkHandler(repo domain.IRepository) CompleteRethinkHandler {
	return CompleteRethinkHandler{repo: repo}
}

func (h CompleteRethinkHandler) Handle(ctx context.Context, cmd CompleteRethink) error {
	completeRethink, err := domain.NewCompleteRethink(cmd.RethinkID, cmd.UserID, cmd.ReportContent, cmd.RethinkContent)
	if err != nil {
		return err
	}
	rethink, err := h.repo.FindRethinkByID(ctx, cmd.RethinkID)
	if err != nil {
		return err
	}
	newRethink, err := completeRethink.Complete(rethink)
	if err != nil {
		return err
	}
	return h.repo.SaveRethink(ctx, newRethink)
}
