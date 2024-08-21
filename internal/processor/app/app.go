package app

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateBot command.CreateBotHandler
	Process   command.ProcessHandler
	StartBot  command.StartBotHandler
	StopBot   command.StopBotHandler
}

type Queries struct {
	AllAnswers query.AllAnswersHandler
}
