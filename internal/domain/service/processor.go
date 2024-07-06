package service

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"log/slog"
)

type Processor struct {
	log     *slog.Logger
	ansRepo interfaces.AnswerRepository
	blcRepo interfaces.BlockRepository
	botRepo interfaces.BotRepository
	prtRepo interfaces.ParticipantRepository
}

type Option func(*Processor) error

func NewProcessor(
	log *slog.Logger,
	ansRepo interfaces.AnswerRepository,
	blcRepo interfaces.BlockRepository,
	botRepo interfaces.BotRepository,
	prtRepo interfaces.ParticipantRepository,
) Processor {

	return Processor{
		log:     log,
		ansRepo: ansRepo,
		blcRepo: blcRepo,
		botRepo: botRepo,
		prtRepo: prtRepo,
	}
}

type Message struct {
	Text    string
	Options []string
}

func (p *Processor) Process(
	ctx context.Context,
	botId value.BotId,
	userId value.UserId,
	ans string,
) (res []Message, err error) {
	const op = "processor.Service.Process"

	log := p.log.With(
		slog.String("op", op),
		slog.Uint64("bot_id", uint64(botId)), // Make new slog Attr for IDs?
		slog.Uint64("user_id", uint64(userId)),
		slog.String("answer", ans))

	log.Info("processing answer")

	prtId := value.ParticipantId{
		BotId:  botId,
		UserId: userId,
	}

	var state value.State

	prt, err := p.prtRepo.Participant(ctx, prtId)
	if err != nil {
		if errors.Is(err, interfaces.ErrParticipantNotFound) {
			b, err := p.botRepo.Bot(ctx, botId)
			if err != nil {
				log.Error("failed to get bot", "err", err)
				return
			}

			prt, err := entity.NewParticipant(prtId, b.Start)
			if err != nil {
				log.Error("failed to create participant", "err", err)
				return
			}

			err = p.prtRepo.Save(ctx, prt)
			if err != nil {
				log.Error("failed to save participant", "err", err)
				return
			}

			state = b.Start
		} else {
			return
		}
	} else {
		block, err := p.blcRepo.Block(ctx, botId, prt.Current)
		if err != nil {
			log.Error("failed to get module", "err", err)
			return
		}

		state, err = block.Node.Next(ans)
		if errors.Is(err, value.ErrIncorrectAnswer) {
			res = append(res, Message{
				Text:    "Некорректный ввод!",
				Options: nil,
			})
		}
	}

	log.Info("processing state", "state", state)

	if state.IsNone() {
		log.Info("end of script")
		return
	}

	next, err := p.blcRepo.Block(ctx, botId, state)
	if err != nil {
		log.Error("failed to get module", "err", err)
		return
	}

	log.Info("got next state", "next", next.Node.State)

	res = append(res, mapBlockToMessage(next))

	prt.Current = state
	err = p.prtRepo.UpdateState(ctx, prtId, prt.Current)
	if err != nil {
		log.Error("failed to update participant's state", "err", err)
		return
	}

	if next.Node.Type == value.Message {
		log.Info("auto process message node", "next", next.Node.State)
		return p.Process(ctx, botId, userId, ans)
	}

	log.Info("end processing module", "state", state)

	return
}

func mapBlockToMessage(block *entity.Block) Message {
	dtoOptions := make([]string, len(block.Node.Options))

	for i, option := range block.Node.Options {
		dtoOptions[i] = option.Text
	}

	return Message{
		Text:    block.Text,
		Options: dtoOptions,
	}
}
