package infra

import (
	"context"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
)

const (
	updatesTimeoutS = 30
	poolingInterval = time.Millisecond
)

type Processor interface {
	Process(ctx context.Context, botUUID string, userID int64, text string) error
}

type tgBot struct {
	botUUID string
	api     *tgbotapi.BotAPI
	upd     tgbotapi.UpdatesChannel
	stop    chan struct{}

	log  *slog.Logger
	proc Processor
}

func newTgBot(
	botUUID string,
	token string,
	proc Processor,
	log *slog.Logger,
) (*tgBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	conf := tgbotapi.NewUpdate(0)
	conf.Timeout = updatesTimeoutS
	updates, err := api.GetUpdatesChan(conf)
	if err != nil {
		return nil, err
	}

	return &tgBot{
		botUUID: botUUID,
		api:     api,
		upd:     updates,
		stop:    make(chan struct{}),
		proc:    proc,
		log:     log,
	}, nil
}

func (b *tgBot) Stop() {
	b.stop <- struct{}{}
	close(b.stop)
}

func (b *tgBot) Listen() {
	run := true
	for run {
		select {
		case u := <-b.upd:
			go b.handleUpdate(u)
		case <-b.stop:
			run = false
		}
		time.Sleep(poolingInterval)
	}
	b.log.Info("bot stopped")
	b.api.StopReceivingUpdates()
}

func (b *tgBot) handleMessage(m *tgbotapi.Message) {
	ctx := context.Background()

	err := b.proc.Process(ctx, b.botUUID, m.Chat.ID, m.Text)
	if err != nil {
		b.log.Error("failed to process message", "err", err.Error())
	}
}

func (b *tgBot) handleUpdate(u tgbotapi.Update) {
	switch {
	case u.Message != nil:
		b.handleMessage(u.Message)
	}
}

type telegramActorService struct {
	log       *slog.Logger
	bots      map[string]*tgBot
	processor Processor
}

func NewTelegramActorService(
	log *slog.Logger,
	processor Processor,
) interfaces.ActorService {
	return &telegramActorService{
		log:       log,
		bots:      make(map[string]*tgBot),
		processor: processor,
	}
}

func (s *telegramActorService) Start(_ context.Context, bot *bots.Bot) error {
	tgb, err := newTgBot(bot.UUID, bot.Token, s.processor, s.log)
	if err != nil {
		return err
	}

	s.bots[bot.UUID] = tgb

	go tgb.Listen()

	return nil
}

func (s *telegramActorService) Stop(_ context.Context, bot *bots.Bot) error {
	tgb, ok := s.bots[bot.UUID]
	if !ok {
		return interfaces.BotNotFoundError{UUID: bot.UUID}
	}

	tgb.Stop()
	delete(s.bots, bot.UUID)

	return nil
}
