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
	CreateBot     command.CreateBotHandler
	DeleteBot     command.DeleteBotHandler
	StartBot      command.StartBotHandler
	StopBot       command.StopBotHandler
	UpdateStatus  command.UpdateStatusHandler
	Entry         command.EntryHandler
	Process       command.ProcessHandler
	CreateMailing command.CreateMailingHandler
	StartMailing  command.StartMailingHandler
}

type Queries struct {
	AllAnswers  query.GetAnswersTableHandler
	GetBot      query.GetBotHandler
	GetBots     query.GetBotsHandler
	StartedBots query.GetStartedBotsHandler
}
