package service

import (
	"context"
	"wake_up_backend/internal/common/engines"
	"wake_up_backend/internal/user/adaptor"
	"wake_up_backend/internal/user/app"
	"wake_up_backend/internal/user/app/command"
	"wake_up_backend/internal/user/app/query"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	mysqlEngine, err := engines.CreateMysqlEngine(engines.MysqlConf{
		User:         "root",
		Password:     "123456",
		Host:         "localhost",
		Port:         3306,
		Database:     "wake_up",
		Charset:      "utf8mb4",
		ShowSql:      true,
		MysqlMaxOpen: 10,
		MysqlMaxIdle: 5,
	})
	if err != nil {
		panic(err)
	}
	repoAdaptor := adaptor.NewReportRepository(mysqlEngine)
	return app.Application{
		Commands: app.Commands{
			Token: command.NewUserTokenHandler(repoAdaptor),
		},
		Queries: app.Queries{
			GetUser: query.NewGetUserHandler(repoAdaptor),
		},
	}, func() {}
}
