package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/itsreg-bots/internal/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs"
)

type Port struct {
	bots map[string]*telegramBot

	msgCons messagesConsumer
	runCons runnerConsumer

	app *app.Application
	log *slog.Logger
}

func RunTelegramPort(
	app *app.Application,
	msgCh <-chan *message.Message,
	runCh <-chan *message.Message,
) {
	log := logs.DefaultLogger()

	p := Port{
		bots:    make(map[string]*telegramBot),
		msgCons: messagesConsumer{},
		app:     app,
		log:     log,
	}

	wg := sync.WaitGroup{}

	msgCon := newMessagesConsumer(msgCh, p.handleBotMessage, log)
	wg.Add(1)
	go func() {
		msgCon.Process()
		wg.Done()
	}()

	runCon := newRunnerConsumer(runCh, p.handleRunnerMessage, log)
	wg.Add(1)
	go func() {
		runCon.Process()
		wg.Done()
	}()

	wg.Wait()
}

func (p *Port) handleBotMessage(_ context.Context, msg botMessage) error {
	tgBot, ok := p.bots[msg.BotUUID]
	if !ok {
		err := fmt.Errorf("bot not found: %s", msg.BotUUID)
		return err
	}

	return tgBot.SendMessage(msg.UserID, msg.Text, msg.Buttons)
}

func (p *Port) handleRunnerMessage(ctx context.Context, msg runnerMessage) error {
	switch msg.Command {
	case "start":
		return p.startBot(ctx, msg.BotUUID)
	case "stop":
		return p.stopBot(ctx, msg.BotUUID)
	}
	return fmt.Errorf("invalid command: %s", msg.Command)
}

func (p *Port) startBot(ctx context.Context, botUUID string) error {
	if _, ok := p.bots[botUUID]; ok {
		return nil
	}

	tgBot, err := newTelegramBot(ctx, botUUID, p.app, p.log)
	if err != nil {
		return err
	}

	err = tgBot.Start(ctx)
	if err != nil {
		return err
	}

	p.bots[botUUID] = tgBot

	return nil
}

func (p *Port) stopBot(ctx context.Context, botUUID string) error {
	tgBot, ok := p.bots[botUUID]
	if !ok {
		return fmt.Errorf("bot not found: %s", botUUID)
	}

	delete(p.bots, botUUID)

	return tgBot.Stop(ctx)
}
