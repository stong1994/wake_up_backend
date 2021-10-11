package command

import (
	"context"
	"wake_up_backend/internal/user/domain"
)

type UserTokenHandler struct {
	repo domain.IRepository
}

func NewUserTokenHandler(repo domain.IRepository) UserTokenHandler {
	return UserTokenHandler{repo: repo}
}

func (h UserTokenHandler) Handle(ctx context.Context, account, password string) (string, error) {
	user, err := h.repo.GetLoginUser(ctx, account, password)
	if err != nil {
		return "", err
	}
	return user.GenToken()
}
