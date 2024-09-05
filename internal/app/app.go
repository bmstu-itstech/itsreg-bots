package app

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBot    command.CreateBotHandler
	StartBot     command.StartBotHandler
	StopBot      command.StopBotHandler
	UpdateStatus command.UpdateStatusHandler
	Entry        command.EntryHandler
	Process      command.ProcessHandler
}

type Queries struct {
	AllAnswers query.GetAnswersTableHandler
	GetBot     query.GetBotHandler
}
