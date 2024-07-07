package application

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"log/slog"
)

type BotsProcessor struct {
	log     *slog.Logger
	ansRepo interfaces.AnswerRepository
	blcRepo interfaces.BlockRepository
	botRepo interfaces.BotRepository
	prtRepo interfaces.ParticipantRepository
}

type Option func(*BotsProcessor) error

func NewProcessor(
	log *slog.Logger,
	ansRepo interfaces.AnswerRepository,
	blcRepo interfaces.BlockRepository,
	botRepo interfaces.BotRepository,
	prtRepo interfaces.ParticipantRepository,
) *BotsProcessor {
	return &BotsProcessor{
		log:     log,
		ansRepo: ansRepo,
		blcRepo: blcRepo,
		botRepo: botRepo,
		prtRepo: prtRepo,
	}
}

func (p *BotsProcessor) Process(
	ctx context.Context,
	botId uint64,
	userId uint64,
	ans string,
) ([]dto.Message, error) {
	const op = "BotsProcessor.Process"

	var res []dto.Message

	log := p.log.With(slog.String("op", op))

	log.Info("processing answer", "bot_id", botId, "user_id", userId, "ans", ans)

	prtId, err := value.NewParticipantId(value.BotId(botId), value.UserId(userId))
	if err != nil {
		log.Error("invalid participant id", "err", err.Error())
		return nil, err
	}

	var state value.State

	prt, err := p.prtRepo.Participant(ctx, prtId)
	if err != nil {
		if errors.Is(err, interfaces.ErrParticipantNotFound) {
			b, err := p.botRepo.Bot(ctx, prtId.BotId)
			if err != nil {
				log.Error("failed to get bot", "err", err.Error())
				return res, err
			}

			prt, err = entity.NewParticipant(prtId, b.Start)
			if err != nil {
				log.Error("failed to create participant", "err", err.Error())
				return res, err
			}

			err = p.prtRepo.Save(ctx, prt)
			if err != nil {
				log.Error("failed to save participant", "err", err.Error())
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

		block, err := p.blcRepo.Block(ctx, prtId.BotId, prt.Current)
		if err != nil {
			log.Error("failed to get block", "block", prt.Current, "err", err.Error())
			return res, err
		}

		state, err = block.Node.Next(ans)
		if errors.Is(err, value.ErrIncorrectAnswer) {
			res = append(res, dto.Message{
				Text:    "Некорректный ввод!",
				Options: nil,
			})
		} else if block.Node.Type != value.Message {
			log.Info("save participant's answer")
			ansId, err := value.NewAnswerId(prtId, block.Node.State)
			if err != nil {
				log.Error("failed to create answer id", "err", err.Error())
				return res, err
			}

			answer, err := entity.NewAnswer(ansId, ans)
			if err != nil {
				log.Error("failed to create answer", "err", err.Error())
				return res, err
			}

			err = p.ansRepo.Save(ctx, answer)
			if err != nil {
				log.Error("failed to save answer", "err", err.Error())
				return res, err
			}
		}
	}

	prt.Current = state
	err = p.prtRepo.UpdateState(ctx, prtId, prt.Current)
	if err != nil {
		log.Error("failed to update participant's state", "err", err.Error())
		return res, err
	}

	if state.IsNone() {
		log.Info("participant finished script")
		return res, nil
	}

	log.Info("processing state", "state", state)

	next, err := p.blcRepo.Block(ctx, prtId.BotId, state)
	if err != nil {
		log.Error("failed to get block", "err", err.Error())
		return res, err
	}

	log.Info("got next state", "next", next.Node.State)

	res = append(res, mapBlockToMessage(next))

	if next.Node.Type == value.Message {
		log.Info("auto process message node", "next", next.Node.State)
		newRes, err := p.Process(ctx, botId, userId, ans)
		if err != nil {
			log.Error("failed to process message", "err", err.Error())
			return res, err
		}
		res = append(res, newRes...)
	}

	log.Info("end processing block", "state", state)

	return res, err
}

func mapBlockToMessage(block *entity.Block) dto.Message {
	dtoOptions := make([]string, len(block.Node.Options))

	for i, option := range block.Node.Options {
		dtoOptions[i] = option.Text
	}

	return dto.Message{
		Text:    block.Text,
		Options: dtoOptions,
	}
}
