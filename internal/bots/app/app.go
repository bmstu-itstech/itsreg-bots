package app

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBot command.CreateBotHandler
	Entry     command.EntryHandler
	Process   command.ProcessHandler
}

type Queries struct {
	AllAnswers query.GetAnswersTableHandler
	GetBot     query.GetBotHandler
}
