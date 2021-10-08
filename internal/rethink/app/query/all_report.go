package query

import (
	"context"
)

type AllReportHandler struct {
	readModel AllReportReadModel
}

func NewAllReportHandler(readModel AllReportReadModel) AllReportHandler {
	return AllReportHandler{readModel: readModel}
}

type AllReportReadModel interface {
	FindAllReport(ctx context.Context, userID string, pageNo, pageSize int) ([]AllReport, error)
}

func (h AllReportHandler) Handle(ctx context.Context, userID string, pageNo, pageSize int) ([]AllReport, error) {
	return h.readModel.FindAllReport(ctx, userID, pageNo, pageSize)
}
