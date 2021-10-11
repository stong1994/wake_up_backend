package app

import (
	"wake_up_backend/internal/user/app/command"
	"wake_up_backend/internal/user/app/query"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
	Token command.UserTokenHandler
}

type Queries struct {
	GetUser query.GetUserHandler
}

func NewApplication(commands Commands, queries Queries) Application {
	return Application{
		Commands: commands,
		Queries:  queries,
	}
}
