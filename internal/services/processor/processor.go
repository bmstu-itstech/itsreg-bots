package processor

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	botmemory "github.com/zhikh23/itsreg-bots/internal/domain/bot/memory"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	modulememory "github.com/zhikh23/itsreg-bots/internal/domain/module/memory"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	partmemory "github.com/zhikh23/itsreg-bots/internal/domain/participant/memory"
	"github.com/zhikh23/itsreg-bots/internal/domain/sender"
	"github.com/zhikh23/itsreg-bots/internal/domain/sender/recorder"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/lib/logger/handlers/slogdiscard"
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

func (s *Service) Process(botId int64, userId int64, msg string) error {
	prtId := entity.ParticipantId{
		BotId:  botId,
		UserId: userId,
	}

	var nodeId entity.NodeId

	prt, err := s.participants.Get(prtId)
	if err != nil {
		if errors.Is(err, participant.ErrParticipantNotFound) {
			b, err := s.bots.Get(botId)
			if err != nil {
				return err
			}

			prt := entity.Participant{
				CurrentId:     b.Start,
				ParticipantId: prtId,
			}

			err = s.participants.Save(prt)
			if err != nil {
				return err
			}

			nodeId = b.Start
		} else {
			return err
		}
	} else {
		mod, err := s.modules.Get(botId, prt.CurrentId)
		if err != nil {
			return err
		}
		nodeId = mod.Process(msg)
	}

	if nodeId == entity.NodeNodeId {
		return nil // End
	}

	next, err := s.modules.Get(botId, nodeId)
	if err != nil {
		return err
	}

	err = s.sender.SendMessage(prtId, next.Text, next.Buttons)
	if err != nil {
		return err
	}

	prt.CurrentId = nodeId
	err = s.participants.UpdateCurrentId(prtId, prt.CurrentId)
	if err != nil {
		return err
	}

	if next.IsSilent {
		return s.Process(botId, userId, msg)
	}

	return nil
}

func WithLogger(logger *slog.Logger) Configuration {
	return func(s *Service) error {
		s.logger = logger
		return nil
	}
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

func WithMemoryModuleRepository() Configuration {
	return func(s *Service) error {
		s.modules = modulememory.New()
		return nil
	}
}

func WithMemoryParticipantRepository() Configuration {
	return func(s *Service) error {
		s.participants = partmemory.New()
		return nil
	}
}

func WithMemoryBotRepository() Configuration {
	return func(s *Service) error {
		s.bots = botmemory.New()
		return nil
	}
}

func WithRecorderSender() Configuration {
	return func(s *Service) error {
		s.sender = recorder.New()
		return nil
	}
}
