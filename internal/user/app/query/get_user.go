package query

import (
	"context"
)

type GetUserHandler struct {
	readModel GetUserReadModel
}

func NewGetUserHandler(readModel GetUserReadModel) GetUserHandler {
	return GetUserHandler{readModel: readModel}
}

type GetUserReadModel interface {
	GetUserByAccount(ctx context.Context, account, password string) (User, error)
}

func (h GetUserHandler) Handle(ctx context.Context, account, password string) (User, error) {
	return h.readModel.GetUserByAccount(ctx, account, password)
}
