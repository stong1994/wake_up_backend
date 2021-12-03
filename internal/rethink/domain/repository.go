package domain

import (
	"context"
)

type IRepository interface {
	AddReport(ctx context.Context, report Report) (string, error)
	AddReportGroup(ctx context.Context, reportGroup ReportGroup) (string, error)
	CheckGroup(ctx context.Context, userID, groupID string) (bool, error)
	FindRethinkByID(ctx context.Context, rethinkID string) (Rethink, error)
	SaveRethink(ctx context.Context, rethink Rethink) error
}
