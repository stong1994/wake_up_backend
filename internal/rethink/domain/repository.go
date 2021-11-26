package domain

import (
	"context"
)

type IRepository interface {
	AddReport(ctx context.Context, report Report) (string, error)
	AddReportGroup(ctx context.Context, reportGroup ReportGroup) (string, error)
	CheckGroup(ctx context.Context, userID, groupID string) (bool, error)
}
