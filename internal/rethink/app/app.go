package app

import (
	"wake_up_backend/internal/rethink/app/command"
	"wake_up_backend/internal/rethink/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddReport      command.AddReportHandler
	AddReportGroup command.AddReportGroupHandler
}

type Queries struct {
	ReportAllTypeList query.ReportAllGroupListHandler
}

func NewApplication(commands Commands, queries Queries) Application {
	return Application{
		Commands: commands,
		Queries:  queries,
	}
}
