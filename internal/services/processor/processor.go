package processor

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"github.com/zhikh23/itsreg-bots/internal/domain/sender"
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
	sender       sender.Sender
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

func (s *Service) Process(ctx context.Context, botId int64, userId int64, msg string) error {
	prtId := entity.ParticipantId{
		BotId:  botId,
		UserId: userId,
	}

	var state objects.State

	prt, err := s.participants.Get(prtId)
	if err != nil {
		if errors.Is(err, participant.ErrParticipantNotFound) {
			b, err := s.bots.Get(botId)
			if err != nil {
				return err
			}

			prt := entity.Participant{
				State:         b.Start,
				ParticipantId: prtId,
			}

			err = s.participants.Save(prt)
			if err != nil {
				return err
			}

			state = b.Start
		} else {
			return err
		}
	} else {
		mod, err := s.modules.Get(botId, prt.State)
		if err != nil {
			return err
		}
		state = mod.Process(msg)
	}

	if state == objects.EndState {
		return nil // End
	}

	next, err := s.modules.Get(botId, state)
	if err != nil {
		return err
	}

	err = s.sender.SendMessage(prtId, next.Text, next.Buttons)
	if err != nil {
		return err
	}

	prt.State = state
	err = s.participants.UpdateCurrentId(prtId, prt.State)
	if err != nil {
		return err
	}

	if next.IsSilent {
		return s.Process(ctx, botId, userId, msg)
	}

	return nil
}

func WithBotRepository(bots bot.Repository) Configuration {
	return func(s *Service) error {
		s.bots = bots
		return nil
	}
}

func WithParticipantRepository(participants participant.Repository) Configuration {
	return func(s *Service) error {
		s.participants = participants
		return nil
	}
}

func WithModuleRepository(modules module.Repository) Configuration {
	return func(s *Service) error {
		s.modules = modules
		return nil
	}
}

func WithSender(sender sender.Sender) Configuration {
	return func(s *Service) error {
		s.sender = sender
		return nil
	}
}

func WithDiscardLogger() Configuration {
	return func(s *Service) error {
		s.logger = slogdiscard.NewDiscardLogger()
		return nil
	}
}
