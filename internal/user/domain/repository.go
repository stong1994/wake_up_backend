package domain

import (
	"context"
)

type IRepository interface {
	GetLoginUser(ctx context.Context, account, password string) (User, error)
}
