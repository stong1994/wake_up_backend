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
	AddReport command.AddReportHandler
}

type Queries struct {
	ReportAllTypeList     query.ReportAllTypeListHandler
}
