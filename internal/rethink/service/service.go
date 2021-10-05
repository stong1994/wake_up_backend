package service

import (
	"context"
	"wake_up_backend/internal/common/engines"
	"wake_up_backend/internal/rethink/adaptor"
	"wake_up_backend/internal/rethink/app"
	"wake_up_backend/internal/rethink/app/command"
	"wake_up_backend/internal/rethink/app/query"
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
			AddReport:      command.NewAddReportHandler(repoAdaptor),
			AddReportGroup: command.NewAddReportGroupHandler(repoAdaptor),
		},
		Queries: app.Queries{
			ReportAllTypeList: query.NewReportAllGroupListHandler(repoAdaptor),
			ReportGroupList:   query.NewReportGroupListHandler(repoAdaptor),
		},
	}, func() {}
}
