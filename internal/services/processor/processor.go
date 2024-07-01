package processor

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	botmemory "github.com/zhikh23/itsreg-bots/internal/domain/bot/memory"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	modulememory "github.com/zhikh23/itsreg-bots/internal/domain/module/memory"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	partmemory "github.com/zhikh23/itsreg-bots/internal/domain/participant/memory"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/lib/logger/handlers/slogdiscard"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	"log/slog"
)

type Configuration func(s *Service) error

type Service struct {
	logger       *slog.Logger
	bots         bot.Repository
	modules      module.Repository
	participants participant.Repository
}

func New(cfgs ...Configuration) (*Service, error) {
	s := &Service{}

	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) Process(ctx context.Context, botId int64, userId int64, ans string) ([]objects.Message, error) {
	const op = "processor.Service.Process"

	log := s.logger.With(
		slog.String("op", op),
		slog.Int64("bot_id", botId),
		slog.Int64("user_id", userId),
		slog.String("ans", ans))

	log.Info("processing message", "ans", ans)

	prtId := entity.ParticipantId{
		BotId:  botId,
		UserId: userId,
	}

	var state objects.State
	var messages []objects.Message

	prt, err := s.participants.Get(prtId)
	if err != nil {
		if errors.Is(err, participant.ErrParticipantNotFound) {
			b, err := s.bots.Get(botId)
			if err != nil {
				log.Error("failed to get bot", "err", err)
				return messages, err
			}

			prt := entity.Participant{
				State:         b.Start,
				ParticipantId: prtId,
			}

			err = s.participants.Save(prt)
			if err != nil {
				log.Error("failed to save participant", "err", err)
				return messages, err
			}

			state = b.Start
		} else {
			return messages, err
		}
	} else {
		mod, err := s.modules.Get(botId, prt.State)
		if err != nil {
			log.Error("failed to get module", "err", err)
			return messages, err
		}
		state = mod.Process(ans)
	}

	log.Info("processing state", "state", state)

	if state == objects.StateNone {
		log.Info("end state")
		return messages, nil // End
	}

	next, err := s.modules.Get(botId, state)
	if err != nil {
		log.Error("failed to get module", "err", err)
		return messages, err
	}

	log.Info("got next state", "next", next.Id)

	messages = append(messages, objects.Message{
		Text:    next.Text,
		Buttons: next.Buttons,
	})

	prt.State = state
	err = s.participants.UpdateState(prtId, prt.State)
	if err != nil {
		log.Error("failed to update current participant", "err", err)
		return messages, err
	}

	if next.IsSilent {
		log.Info("auto process silent module", "next", next.Id)
		buf, err := s.Process(ctx, botId, userId, ans)
		if err != nil {
			return messages, err
		}
		messages = append(messages, buf...)
	}

	log.Info("end processing module", "state", state)

	return messages, nil
}

func WithBotRepository(bots bot.Repository) Configuration {
	return func(s *Service) error {
		s.bots = bots
		return nil
	}
}

func WithMemoryBotRepository() Configuration {
	return func(s *Service) error {
		s.bots = botmemory.New()
		return nil
	}
}

func WithParticipantRepository(participants participant.Repository) Configuration {
	return func(s *Service) error {
		s.participants = participants
		return nil
	}
}

func WithMemoryParticipantRepository() Configuration {
	return func(s *Service) error {
		s.participants = partmemory.New()
		return nil
	}
}

func WithModuleRepository(modules module.Repository) Configuration {
	return func(s *Service) error {
		s.modules = modules
		return nil
	}
}

func WithMemoryModuleRepository() Configuration {
	return func(s *Service) error {
		s.modules = modulememory.New()
		return nil
	}
}

func WithLogger(logger *slog.Logger) Configuration {
	return func(s *Service) error {
		s.logger = logger
		return nil
	}
}

func WithDiscardLogger() Configuration {
	return func(s *Service) error {
		s.logger = slogdiscard.NewDiscardLogger()
		return nil
	}
}
