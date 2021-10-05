package query

import (
	"context"
)

type ReportGroupList struct {
	PageNum  int
	PageSize int
	UserID   string
}

type ReportGroupListHandler struct {
	readModel ReportGroupListReadModel
}

func NewReportGroupListHandler(readModel ReportGroupListReadModel) ReportGroupListHandler {
	return ReportGroupListHandler{readModel: readModel}
}

type ReportGroupListReadModel interface {
	FindReportGroups(ctx context.Context, userID string) (RespReportGroupList, error)
}

func (h ReportGroupListHandler) Handle(ctx context.Context, userID string) (RespReportGroupList, error) {
	return h.readModel.FindReportGroups(ctx, userID)
}
