package adaptor

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
	"wake_up_backend/internal/user/app/query"
	"wake_up_backend/internal/user/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type UserModel struct {
	ID          string
	Name        string
	DisplayName string
	Account     string
	Password    string
	CreateTime  time.Time
	DeleteTime  time.Time
}

func (r UserModel) TableName() string {
	return "user"
}

type ReportRepository struct {
	client *gorm.DB
}

func NewReportRepository(dbClient *gorm.DB) ReportRepository {
	return ReportRepository{client: dbClient}
}

func (r ReportRepository) GetUserByAccount(ctx context.Context, account, password string) (query.User, error) {
	var model UserModel
	if err := r.client.Where("account = ? AND password = ?", account, password).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.User{}, ErrNotFound
		}
		return query.User{}, err
	}
	return query.User{
		ID:          model.ID,
		DisplayName: model.DisplayName,
	}, nil
}

func (r ReportRepository) GetLoginUser(ctx context.Context, account, password string) (domain.User, error) {
	var model UserModel
	if err := r.client.Where("account = ? AND password = ?", account, password).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return domain.UnmarshalFromDB(model.ID, model.Name, model.DisplayName, model.DeleteTime)
}
