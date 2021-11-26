package query

import (
	"context"
)

type ReportAllTypeList struct {
	PageNum  int
	PageSize int
	UserID   string
}

type ReportAllGroupListHandler struct {
	readModel ReportAllGroupListReadModel
}

func NewReportAllGroupListHandler(readModel ReportAllGroupListReadModel) ReportAllGroupListHandler {
	return ReportAllGroupListHandler{readModel: readModel}
}

type ReportAllGroupListReadModel interface {
	FindUserReports(ctx context.Context, userID string, pageNo, pageSize int) (RespReportAllGroupList, error)
}

func (h ReportAllGroupListHandler) Handle(ctx context.Context, query ReportAllTypeList) (RespReportAllGroupList, error) {
	return h.readModel.FindUserReports(ctx, query.UserID, query.PageNum, query.PageSize)
}
