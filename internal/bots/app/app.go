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
	Process   command.ProcessHandler
}

type Queries struct {
	AllAnswers query.AllAnswersHandler
	GetBot     query.GetBotHandler
}
