package telegram

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/query"
)

const (
	updatesTimeout = 30
)

type TelegramServer struct {
	app *app.Application
	log *slog.Logger

	bots map[string]*chatHandler
}

func NewTelegramServer(app *app.Application) *TelegramServer {
	return &TelegramServer{app: app}
}

func (s TelegramServer) Start(ctx context.Context, botUUID string) error {
	bot, err := s.app.Queries.GetBot.Handle(ctx, query.GetBot{BotUUID: botUUID})
	if err != nil {
		return err
	}

	api, err := tgbotapi.NewBotAPI(bot.Token)
	if err != nil {
		return err
	}

	conf := tgbotapi.NewUpdate(0)
	conf.Timeout = updatesTimeout
	ch, err := api.GetUpdatesChan(conf)
	if err != nil {
		return err
	}

	handler := &chatHandler{
		bot: bot,
		app: app.Application{},
		ch:  nil,
	}
	s.bots[bot.UUID] = handler

	go handler.listen(ch)

	return nil
}

type chatHandler struct {
	app app.Application
	log *slog.Logger

	bot query.Bot
	ch  tgbotapi.UpdatesChannel
}

func (h chatHandler) listen(ch tgbotapi.UpdatesChannel) {
	for {
		select {
		case u := <-ch:
			h.handle(u)
		}
	}
}

func (h chatHandler) handle(u tgbotapi.Update) {
	switch {
	case u.Message != nil:
		h.handleMessage(u.Message)
	}
}

func (h chatHandler) handleMessage(m *tgbotapi.Message) {
	ctx := context.Background()
	err := h.app.Commands.Process.Handle(ctx, command.Process{
		BotUUID: h.bot.UUID,
		UserID:  m.Chat.ID,
		Text:    m.Text,
	})
	if err != nil {
		h.log.Error("failed to process message", "error", err.Error())
	}
}
