package domain

import (
	"context"
)

type IRepository interface {
	AddReport(ctx context.Context, report Report) error
	AddReportGroup(ctx context.Context, reportGroup ReportGroup) error
	CheckGroup(ctx context.Context, userID, groupID string) (bool, error)
}
