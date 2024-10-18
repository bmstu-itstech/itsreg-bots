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
	"github.com/bmstu-itstech/itsreg-bots/internal/common/jwtauth"
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
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = s.app.Commands.CreateBot.Handle(r.Context(), command.CreateBot{
		BotUUID:    postBots.BotUUID,
		AuthorUUID: userUUID,
		Name:       postBots.Name,
		Token:      postBots.Token,
		Entries:    convertEntryPointsFromAPI(postBots.Entries),
		Mailings:   convertOptionalMailingsFromAPI(postBots.Mailings),
		Blocks:     convertBlocksFromAPI(postBots.Blocks),
	})
	if errors.As(err, &commonerrs.InvalidInputError{}) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-location", fmt.Sprintf("/bots/%s", postBots.BotUUID))
	w.WriteHeader(http.StatusCreated)
}

func (s Server) DeleteBot(w http.ResponseWriter, r *http.Request, botUUID string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = s.app.Commands.DeleteBot.Handle(r.Context(), command.DeleteBot{
		AuthorUUID: userUUID,
		BotUUID:    botUUID,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) GetBots(w http.ResponseWriter, r *http.Request) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	bs, err := s.app.Queries.GetBots.Handle(r.Context(), query.GetBots{
		UserUUID: userUUID,
	})
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, convertBotsToAPI(bs))
}

func (s Server) StartBot(w http.ResponseWriter, r *http.Request, uuid string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = s.app.Commands.StartBot.Handle(r.Context(), command.StartBot{
		AuthorUUID: userUUID,
		BotUUID:    uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) StopBot(w http.ResponseWriter, r *http.Request, uuid string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = s.app.Commands.StopBot.Handle(r.Context(), command.StopBot{
		AuthorUUID: userUUID,
		BotUUID:    uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) StartMailing(w http.ResponseWriter, r *http.Request, uuid string, entryKey string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = s.app.Commands.StartMailing.Handle(r.Context(), command.StartMailing{
		AuthorUUID: userUUID,
		BotUUID:    uuid,
		EntryKey:   entryKey,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) GetBot(w http.ResponseWriter, r *http.Request, uuid string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	bot, err := s.app.Queries.GetBot.Handle(r.Context(), query.GetBot{
		UserUUID: userUUID,
		BotUUID:  uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, convertBotToAPI(bot))
}

func (s Server) GetAnswers(w http.ResponseWriter, r *http.Request, uuid string) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	answers, err := s.app.Queries.AllAnswers.Handle(r.Context(), query.GetAnswersTable{
		UserUUID: userUUID,
		BotUUID:  uuid,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, bots.ErrPermissionDenied) {
		httpError(w, r, err, http.StatusForbidden)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVAnswers(w, answers)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)
	render.JSON(w, r, Error{Message: err.Error()})
}

func convertOptionToAPI(option types.Option) Option {
	return Option{
		Next: option.Next,
		Text: option.Text,
	}
}

func convertOptionFromAPI(option Option) types.Option {
	return types.Option{
		Text: option.Text,
		Next: option.Next,
	}
}

func convertOptionsToAPI(options []types.Option) *[]Option {
	res := make([]Option, len(options))
	for i, option := range options {
		res[i] = convertOptionToAPI(option)
	}
	return &res
}

func convertOptionsFromAPI(options *[]Option) []types.Option {
	if options == nil {
		return nil
	}
	res := make([]types.Option, len(*options))
	for i, option := range *options {
		res[i] = convertOptionFromAPI(option)
	}
	return res
}

func convertEntryPointToAPI(entry types.EntryPoint) EntryPoint {
	return EntryPoint{
		Key:   entry.Key,
		State: entry.State,
	}
}

func convertEntryPointFromAPI(entry EntryPoint) types.EntryPoint {
	return types.EntryPoint{
		Key:   entry.Key,
		State: entry.State,
	}
}

func convertEntryPointsToAPI(entries []types.EntryPoint) []EntryPoint {
	res := make([]EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = convertEntryPointToAPI(entry)
	}
	return res
}

func convertEntryPointsFromAPI(entries []EntryPoint) []types.EntryPoint {
	res := make([]types.EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = convertEntryPointFromAPI(entry)
	}
	return res
}

func convertMailingToAPI(mailing types.Mailing) Mailing {
	return Mailing{
		EntryKey:      mailing.EntryKey,
		Name:          mailing.Name,
		RequiredState: mailing.RequiredState,
	}
}

func convertOptionalMailingsToAPI(mailings []types.Mailing) *[]Mailing {
	res := make([]Mailing, len(mailings))
	for i, mailing := range mailings {
		res[i] = convertMailingToAPI(mailing)
	}
	return &res
}

func convertMailingFromAPI(mailing Mailing) types.Mailing {
	return types.Mailing{
		Name:          mailing.Name,
		EntryKey:      mailing.EntryKey,
		RequiredState: mailing.RequiredState,
	}
}

func convertMailingsFromAPI(mailings []Mailing) []types.Mailing {
	res := make([]types.Mailing, len(mailings))
	for i, mailing := range mailings {
		res[i] = convertMailingFromAPI(mailing)
	}
	return res
}

func convertOptionalMailingsFromAPI(mailings *[]Mailing) []types.Mailing {
	if mailings == nil {
		return make([]types.Mailing, 0)
	}

	return convertMailingsFromAPI(*mailings)
}

func convertBlockToAPI(block types.Block) Block {
	return Block{
		Type:      BlockType(block.Type),
		NextState: block.NextState,
		Options:   convertOptionsToAPI(block.Options),
		State:     block.State,
		Text:      block.Text,
		Title:     block.Title,
	}
}

func convertBlockFromAPI(block Block) types.Block {
	return types.Block{
		Type:      string(block.Type),
		State:     block.State,
		NextState: block.NextState,
		Options:   convertOptionsFromAPI(block.Options),
		Title:     block.Title,
		Text:      block.Text,
	}
}

func convertBlocksToAPI(blocks []types.Block) []Block {
	res := make([]Block, len(blocks))
	for i, block := range blocks {
		res[i] = convertBlockToAPI(block)
	}
	return res
}

func convertBlocksFromAPI(blocks []Block) []types.Block {
	res := make([]types.Block, len(blocks))
	for i, block := range blocks {
		res[i] = convertBlockFromAPI(block)
	}
	return res
}

func convertBotToAPI(bot types.Bot) Bot {
	return Bot{
		Blocks:    convertBlocksToAPI(bot.Blocks),
		BotUUID:   bot.UUID,
		CreatedAt: bot.CreatedAt,
		Entries:   convertEntryPointsToAPI(bot.Entries),
		Mailings:  convertOptionalMailingsToAPI(bot.Mailings),
		Name:      bot.Name,
		Status:    BotStatus(bot.Status),
		Token:     bot.Token,
		UpdatedAt: bot.UpdatedAt,
	}
}

func convertBotsToAPI(bs []types.Bot) []Bot {
	res := make([]Bot, len(bs))
	for i, b := range bs {
		res[i] = convertBotToAPI(b)
	}
	return res
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
