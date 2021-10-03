package query

import (
	"context"
)

type ReportAllTypeList struct {
	PageNum int
	PageSize int
	UserID string
}

type ReportAllTypeListHandler struct {
	readModel ReportAllGroupListReadModel
}

type ReportAllGroupListReadModel interface {
	FindReportWithAllGroup(ctx context.Context, userID string, pageNo, pageSize int) (RespReportAllGroupList, error)
}


func (h ReportAllTypeListHandler) Handle(ctx context.Context, query ReportAllTypeList) (RespReportAllGroupList, error) {
	return h.readModel.FindReportWithAllGroup(ctx, query.UserID, query.PageNum, query.PageSize)
}