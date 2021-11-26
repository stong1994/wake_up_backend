package service

import (
	"context"
	"time"
	"wake_up_backend/internal/common/engines"
	"wake_up_backend/internal/rethink/adaptor"
	"wake_up_backend/internal/rethink/app"
	"wake_up_backend/internal/rethink/app/command"
	"wake_up_backend/internal/rethink/app/query"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	dbEngine, err := engines.NewMongoClient(engines.MongoConfig{
		User:        "stong",
		Password:    "123456",
		Addr:        "localhost:27017",
		DBName:      "wake_up",
		AuthSource:  "admin",
		PoolMaxSize: 30,
		IdleTime:    time.Minute,
	})
	if err != nil {
		panic(err)
	}
	repoAdaptor := adaptor.NewReportRepository(dbEngine)
	return app.Application{
		Commands: app.Commands{
			AddReport:      command.NewAddReportHandler(repoAdaptor),
			AddReportGroup: command.NewAddReportGroupHandler(repoAdaptor),
		},
		Queries: app.Queries{
			ReportAllTypeList: query.NewReportAllGroupListHandler(repoAdaptor),
			ReportGroupList:   query.NewReportGroupListHandler(repoAdaptor),
			//AllReport:         query.NewAllReportHandler(repoAdaptor),
		},
	}, func() {}
}
