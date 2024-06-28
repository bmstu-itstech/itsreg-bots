package processor

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"log/slog"
)

var (
	ErrAlreadyFinished = errors.New("participant already finished")
)

// Provider provides modules.
type Provider interface {
	Module(id int32) (*module.Module, error)
}

// Processor processes participant's messages.
type Processor struct {
	log      *slog.Logger
	provider Provider
}

// New returns new Processor.
func New(logger *slog.Logger, provider Provider) *Processor {
	return &Processor{
		log:      logger,
		provider: provider,
	}
}

// Process message from participant and return answers or error.
// Returns ErrAlreadyFinished if participant's current module is last.
func (p *Processor) Process(participant *participant.Participant, msg string) ([]string, error) {
	const op = "processor.Process"

	log := p.log.With(slog.String("op", op)).With(slog.Int("participant", int(participant.Id)))
	log.Info("process message", "current", participant.Current, "msg", msg)

	// Empty array as default return value.
	ans := make([]string, 0)

	// Get current module of participant.
	mod, err := p.provider.Module(participant.Current)
	if err != nil {
		log.Error("failed to get current module",
			"current", participant.Current,
			"err", err)
		return ans, err
	}

	// If current module is last, then there is nothing to process.
	if mod.IsLast() {
		log.Error("module is already finished",
			"current", participant.Current)
		return ans, ErrAlreadyFinished
	}

	// Check the conditions-buttons.
	for _, btn := range mod.Buttons {
		if btn.Value == msg {
			// Switch the participant to a new module.
			mod, err = p.provider.Module(btn.Next)
			if err != nil {
				log.Error("failed to get next module",
					"current", participant.Current,
					"next", btn.Next,
					"err", err)
				return ans, err
			}
			participant.Current = btn.Next

			// Bot will respond with what is saved in the text of the module.
			ans = append(ans, mod.Text)

			// If the next module is silent, then immediately process it.
			if mod.IsSilent {
				auto, err := p.Process(participant, "")
				if err != nil {
					log.Error("failed to auto-process module",
						"current", participant.Current,
						"err", err)
					return ans, err
				}
				ans = append(ans, auto...)
			}

			// End of button branch processing.
			log.Info("finished processing message", "current", participant.Current, "msg", msg)
			return ans, nil
		}
	}

	// If the module does not have buttons or the message does not match any branch of the button, then the default
	// case is handled.
	mod, err = p.provider.Module(mod.Next)
	if err != nil {
		log.Error("failed to auto-process module",
			"current", participant.Current,
			"err", err)
		return ans, nil
	}
	participant.Current = mod.Id

	// Bot will respond with what is saved in the text of the module.
	ans = append(ans, mod.Text)

	// If the next module is silent, then immediately process it.
	if mod.IsSilent {
		auto, err := p.Process(participant, "")
		if err != nil {
			log.Error("failed to get current module",
				"current", participant.Current,
				"err", err)
			return ans, err
		}
		ans = append(ans, auto...)
	}

	// End of default branch processing.
	log.Info("finished processing message", "current", participant.Current, "msg", msg)
	return ans, nil
}
