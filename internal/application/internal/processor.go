package internal

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
	"github.com/zhikh23/itsreg-bots/internal/domain"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"log/slog"
)

type Processor struct {
	log     *slog.Logger
	ansRepo domain.AnswersRepository
	blcRepo domain.BlockRepository
	botRepo domain.BotRepository
	prtRepo domain.ParticipantRepository
}

type ProcessorOption func(*Processor) error

func NewProcessor(opts ...ProcessorOption) (*Processor, error) {
	p := &Processor{}

	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Processor) Process(ctx context.Context, req dto.ProcessRequest) (res dto.ProcessResponse, err error) {
	const op = "processor.Service.Process"

	log := p.log.With(
		slog.String("op", op),
		slog.Uint64("bot_id", req.BotId),
		slog.Uint64("user_id", req.UserId),
		slog.String("answer", req.Answer))

	log.Info("processing answer")

	prtId := value.ParticipantId{
		BotId:  value.BotId(req.BotId),
		UserId: value.UserId(req.UserId),
	}

	var state value.State

	prt, err := p.prtRepo.Get(prtId)
	if err != nil {
		if errors.Is(err, domain.ErrParticipantNotFound) {
			b, err := p.botRepo.Get(prtId.BotId)
			if err != nil {
				log.Error("failed to get bot", "err", err)
				return
			}

			prt, err := entity.NewParticipant(prtId, b.Start)
			if err != nil {
				log.Error("failed to create participant", "err", err)
				return
			}

			err = p.prtRepo.Save(prt)
			if err != nil {
				log.Error("failed to save participant", "err", err)
				return
			}

			state = b.Start
		} else {
			return
		}
	} else {
		block, err := p.blcRepo.Get(prtId.BotId, prt.Current)
		if err != nil {
			log.Error("failed to get module", "err", err)
			return
		}

		state, err = block.Node.Next(req.Answer)
		if errors.Is(err, value.ErrIncorrectAnswer) {
			res.Messages = append(res.Messages, dto.Message{
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

	next, err := p.blcRepo.Get(prtId.BotId, state)
	if err != nil {
		log.Error("failed to get module", "err", err)
		return
	}

	log.Info("got next state", "next", next.Node.State)

	res.Messages = append(res.Messages, dto.MapBlockToMessage(*next))

	prt.Current = state
	err = p.prtRepo.UpdateState(prtId, prt.Current)
	if err != nil {
		log.Error("failed to update participant's state", "err", err)
		return
	}

	if next.Node.Type == value.Message {
		log.Info("auto process message node", "next", next.Node.State)
		return p.Process(ctx, req)
	}

	log.Info("end processing module", "state", state)

	return
}

func WithLogger(log *slog.Logger) ProcessorOption {
	return func(p *Processor) error {
		p.log = log
		return nil
	}
}

func WithAnswersRepository(ansRepo domain.AnswersRepository) ProcessorOption {
	return func(p *Processor) error {
		p.ansRepo = ansRepo
		return nil
	}
}

func WithBlockRepository(blcRepo domain.BlockRepository) ProcessorOption {
	return func(p *Processor) error {
		p.blcRepo = blcRepo
		return nil
	}
}

func WithParticipantRepository(participantRepo domain.ParticipantRepository) ProcessorOption {
	return func(p *Processor) error {
		p.prtRepo = participantRepo
		return nil
	}
}
