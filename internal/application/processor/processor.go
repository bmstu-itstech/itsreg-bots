package processor

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/application/dto"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	interfaces2 "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"log/slog"
)

type Processor struct {
	log     *slog.Logger
	ansRepo interfaces2.AnswerRepository
	blcRepo interfaces2.BlockRepository
	botRepo interfaces2.BotRepository
	prtRepo interfaces2.ParticipantRepository
}

type Option func(*Processor) error

func New(
	log *slog.Logger,
	ansRepo interfaces2.AnswerRepository,
	blcRepo interfaces2.BlockRepository,
	botRepo interfaces2.BotRepository,
	prtRepo interfaces2.ParticipantRepository,
) *Processor {
	return &Processor{
		log:     log,
		ansRepo: ansRepo,
		blcRepo: blcRepo,
		botRepo: botRepo,
		prtRepo: prtRepo,
	}
}

func (p *Processor) Process(
	ctx context.Context,
	botId uint64,
	userId uint64,
	ans string,
) ([]dto.Message, error) {
	const op = "Processor.Process"

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
		if errors.Is(err, interfaces2.ErrParticipantNotFound) {
			// Новый участник бота.
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
		// Уже существующий участник бота.

		// Текущий блок.
		block, err := p.blcRepo.Block(ctx, prtId.BotId, prt.Current)
		if err != nil {
			log.Error("failed to get block", "block", prt.Current, "err", err.Error())
			return res, err
		}

		// Следующее состояние для данного ответа.
		state, err = block.Next(ans)
		if errors.Is(err, value.ErrIncorrectAnswer) {
			res = append(res, dto.Message{
				Text:    "Некорректный ввод!",
				Options: nil,
			})
		} else if block.Type != value.Message {
			log.Info("save participant's answer")
			ansId, err := value.NewAnswerId(prtId, block.State)
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

	if state.IsNone() {
		log.Info("participant finished script")
		return res, nil
	}

	prt.Current = state
	err = p.prtRepo.UpdateState(ctx, prtId, prt.Current)
	if err != nil {
		log.Error("failed to update participant's state", "err", err.Error())
		return res, err
	}

	next, err := p.blcRepo.Block(ctx, prtId.BotId, state)
	if err != nil {
		log.Error("failed to get block", "err", err.Error())
		return res, err
	}

	log.Info("got next state", "next", next.State)

	res = append(res, mapBlockToMessage(next))

	if next.Type == value.Message {
		log.Info("auto process message node", "next", next.State)
		newRes, err := p.Process(ctx, botId, userId, ans)
		if err != nil {
			log.Error("failed to process message", "err", err.Error())
			return res, err
		}
		res = append(res, newRes...)
	}

	log.Info("end processing block", "state", state)

	return res, nil
}

func mapBlockToMessage(block *entity.Block) dto.Message {
	dtoOptions := make([]string, len(block.Options))

	for i, option := range block.Options {
		dtoOptions[i] = option.Text
	}

	return dto.Message{
		Text:    block.Text,
		Options: dtoOptions,
	}
}
