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
) ([]Message, error) {
	const op = "processor.Service.Process"

	var res []Message

	log := p.log.With(
		slog.String("op", op),
		slog.Uint64("bot_id", uint64(botId)), // Make new slog Attr for IDs?
		slog.Uint64("user_id", uint64(userId)),
		slog.String("answer", ans))

	log.Info("processing answer")

	prtId, err := value.NewParticipantId(botId, userId)
	if err != nil {
		log.Error("invalid participant id", "err", err)
		return nil, err
	}

	var state value.State

	prt, err := p.prtRepo.Participant(ctx, prtId)
	if err != nil {
		if errors.Is(err, interfaces.ErrParticipantNotFound) {
			b, err := p.botRepo.Bot(ctx, botId)
			if err != nil {
				log.Error("failed to get bot", "err", err)
				return res, err
			}

			prt, err = entity.NewParticipant(prtId, b.Start)
			if err != nil {
				log.Error("failed to create participant", "err", err)
				return res, err
			}

			err = p.prtRepo.Save(ctx, prt)
			if err != nil {
				log.Error("failed to save participant", "err", err)
				return res, err
			}

			state = b.Start
		} else {
			return res, err
		}
	} else {
		if prt.Current.IsNone() {
			log.Info("participant already finished")
			return res, nil
		}

		block, err := p.blcRepo.Block(ctx, botId, prt.Current)
		if err != nil {
			log.Error("failed to get block", "err", err)
			return res, err
		}

		state, err = block.Node.Next(ans)
		if errors.Is(err, value.ErrIncorrectAnswer) {
			res = append(res, Message{
				Text:    "Некорректный ввод!",
				Options: nil,
			})
		} else if block.Node.Type != value.Message {
			log.Info("save participant's answer")
			ansId, err := value.NewAnswerId(prtId, block.Node.State)
			if err != nil {
				log.Error("failed to create answer id", "err", err)
				return res, err
			}

			answer, err := entity.NewAnswer(ansId, ans)
			if err != nil {
				log.Error("failed to create answer", "err", err)
				return res, err
			}

			err = p.ansRepo.Save(ctx, answer)
			if err != nil {
				log.Error("failed to save answer", "err", err)
				return res, err
			}
		}
	}

	prt.Current = state
	err = p.prtRepo.UpdateState(ctx, prtId, prt.Current)
	if err != nil {
		log.Error("failed to update participant's state", "err", err)
		return res, err
	}

	if state.IsNone() {
		log.Info("participant finished script")
		return res, nil
	}

	log.Info("processing state", "state", state)

	next, err := p.blcRepo.Block(ctx, botId, state)
	if err != nil {
		log.Error("failed to get block", "err", err)
		return res, err
	}

	log.Info("got next state", "next", next.Node.State)

	res = append(res, mapBlockToMessage(next))

	if next.Node.Type == value.Message {
		log.Info("auto process message node", "next", next.Node.State)
		newRes, err := p.Process(ctx, botId, userId, ans)
		if err != nil {
			log.Error("failed to process message", "err", err)
			return res, err
		}
		res = append(res, newRes...)
	}

	log.Info("end processing block", "state", state)

	return res, err
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
