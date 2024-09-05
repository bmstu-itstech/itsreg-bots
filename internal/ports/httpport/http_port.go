package httpport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/bmstu-itstech/itsreg-bots/internal/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type Server struct {
	app *app.Application
}

func NewHTTPServer(app *app.Application) *Server {
	return &Server{app: app}
}

func (s Server) CreateBot(w http.ResponseWriter, r *http.Request) {
	postBots := PostBots{}
	if err := render.Decode(r, &postBots); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.app.Commands.CreateBot.Handle(r.Context(), command.CreateBot{
		BotUUID: postBots.BotUUID,
		Name:    postBots.Name,
		Token:   postBots.Token,
		Entries: mapEntryPointsFromAPI(postBots.Entries),
		Blocks:  mapBlocksFromAPI(postBots.Blocks),
	})
	if errors.As(err, &commonerrs.InvalidInputError{}) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-location", fmt.Sprintf("/bots/%s", postBots.BotUUID))
	w.WriteHeader(http.StatusCreated)
}

func (s Server) StartBot(w http.ResponseWriter, r *http.Request, uuid string) {
	err := s.app.Commands.StartBot.Handle(r.Context(), command.StartBot{
		BotUUID: uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s Server) StopBot(w http.ResponseWriter, r *http.Request, uuid string) {
	err := s.app.Commands.StopBot.Handle(r.Context(), command.StopBot{
		BotUUID: uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s Server) GetBot(w http.ResponseWriter, r *http.Request, uuid string) {
	bot, err := s.app.Queries.GetBot.Handle(r.Context(), query.GetBot{
		BotUUID: uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, mapBotToAPI(bot))
}

func (s Server) GetAnswers(w http.ResponseWriter, r *http.Request, uuid string) {
	answers, err := s.app.Queries.AllAnswers.Handle(r.Context(), query.GetAnswersTable{
		BotUUID: uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderCSVAnswers(w, answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func mapOptionToAPI(option types.Option) Option {
	return Option{
		Next: option.Next,
		Text: option.Text,
	}
}

func mapOptionFromAPI(option Option) types.Option {
	return types.Option{
		Text: option.Text,
		Next: option.Next,
	}
}

func mapOptionsToAPI(options []types.Option) *[]Option {
	res := make([]Option, len(options))
	for i, option := range options {
		res[i] = mapOptionToAPI(option)
	}
	return &res
}

func mapOptionsFromAPI(options *[]Option) []types.Option {
	if options == nil {
		return nil
	}
	res := make([]types.Option, len(*options))
	for i, option := range *options {
		res[i] = mapOptionFromAPI(option)
	}
	return res
}

func mapEntryPointToAPI(entry types.EntryPoint) EntryPoint {
	return EntryPoint{
		Key:   entry.Key,
		State: entry.State,
	}
}

func mapEntryPointFromAPI(entry EntryPoint) types.EntryPoint {
	return types.EntryPoint{
		Key:   entry.Key,
		State: entry.State,
	}
}

func mapEntryPointsToAPI(entries []types.EntryPoint) []EntryPoint {
	res := make([]EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = mapEntryPointToAPI(entry)
	}
	return res
}

func mapEntryPointsFromAPI(entries []EntryPoint) []types.EntryPoint {
	res := make([]types.EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = mapEntryPointFromAPI(entry)
	}
	return res
}

func mapBlockToAPI(block types.Block) Block {
	return Block{
		Type:      BlockType(block.Type),
		NextState: block.State,
		Options:   mapOptionsToAPI(block.Options),
		State:     block.State,
		Text:      block.Text,
		Title:     block.Title,
	}
}

func mapBlockFromAPI(block Block) types.Block {
	return types.Block{
		Type:      string(block.Type),
		State:     block.State,
		NextState: block.NextState,
		Options:   mapOptionsFromAPI(block.Options),
		Title:     block.Title,
		Text:      block.Text,
	}
}

func mapBlocksToAPI(blocks []types.Block) []Block {
	res := make([]Block, len(blocks))
	for i, block := range blocks {
		res[i] = mapBlockToAPI(block)
	}
	return res
}

func mapBlocksFromAPI(blocks []Block) []types.Block {
	res := make([]types.Block, len(blocks))
	for i, block := range blocks {
		res[i] = mapBlockFromAPI(block)
	}
	return res
}

func mapBotToAPI(bot types.Bot) Bot {
	return Bot{
		Blocks:    mapBlocksToAPI(bot.Blocks),
		BotUUID:   bot.UUID,
		CreatedAt: bot.CreatedAt,
		Entries:   mapEntryPointsToAPI(bot.Entries),
		Name:      bot.Name,
		Status:    BotStatus(bot.Status),
		Token:     bot.Token,
		UpdatedAt: bot.UpdatedAt,
	}
}

func renderCSVAnswers(w http.ResponseWriter, answers types.AnswersTable) error {
	csvWriter := csv.NewWriter(w)

	err := csvWriter.Write(answers.THead)
	if err != nil {
		return err
	}
	err = csvWriter.WriteAll(answers.TBody)
	if err != nil {
		return err
	}

	return nil
}
