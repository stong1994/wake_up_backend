package domain

import (
	"context"
)

type IRepository interface {
	AddReport(ctx context.Context, report Report) error
}
