package grpcport

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bmstu-itstech/itsreg-bots/internal/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"

	botspb "github.com/bmstu-itstech/itsreg-bots/api/grpc/gen/bots"
)

type grpcPort struct {
	botspb.UnimplementedBotsServiceServer
	app *app.Application
}

func RegisterGRPCServer(
	grpcServer *grpc.Server,
	app *app.Application,
) {
	botspb.RegisterBotsServiceServer(grpcServer, &grpcPort{app: app})
}

func (p *grpcPort) CreateBot(
	ctx context.Context, req *botspb.CreateBotRequest,
) (*botspb.Empty, error) {
	err := p.app.Commands.CreateBot.Handle(ctx, command.CreateBot{
		BotUUID: req.BotUUID,
		Name:    req.Name,
		Token:   req.Token,
		Entries: mapEntriesFromPB(req.Entries),
		Blocks:  mapBlocksFromPB(req.Blocks),
	})
	if errors.As(err, &commonerrs.InvalidInputError{}) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.As(err, &bots.BotAlreadyExistsError{}) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &botspb.Empty{}, nil
}

func (p *grpcPort) Entry(
	ctx context.Context, req *botspb.EntryRequest,
) (*botspb.Empty, error) {
	err := p.app.Commands.Entry.Handle(ctx, command.Entry{
		BotUUID: req.BotUUID,
		UserID:  req.UserID,
		Key:     req.Key,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.As(err, &bots.EntryNotFoundError{}) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &botspb.Empty{}, nil
}

func (p *grpcPort) Process(
	ctx context.Context, req *botspb.ProcessRequest,
) (*botspb.Empty, error) {
	err := p.app.Commands.Process.Handle(ctx, command.Process{
		BotUUID: req.BotUUID,
		UserID:  req.UserID,
		Text:    req.Text,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &botspb.Empty{}, nil
}

func (p *grpcPort) GetBot(
	ctx context.Context, req *botspb.GetBotRequest,
) (*botspb.Bot, error) {
	bot, err := p.app.Queries.GetBot.Handle(ctx, query.GetBot{
		BotUUID: req.BotUUID,
	})
	if errors.As(err, &bots.BotNotFoundError{}) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, err
	}

	return mapBotToPB(bot), nil
}

func mapEntryFromPB(entry *botspb.EntryPoint) types.EntryPoint {
	return types.EntryPoint{
		Key:   entry.Key,
		State: int(entry.State),
	}
}

func mapEntryToPB(entry types.EntryPoint) *botspb.EntryPoint {
	return &botspb.EntryPoint{
		Key:   entry.Key,
		State: int32(entry.State),
	}
}

func mapEntriesFromPB(entries []*botspb.EntryPoint) []types.EntryPoint {
	res := make([]types.EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = mapEntryFromPB(entry)
	}
	return res
}

func mapEntriesToPB(entries []types.EntryPoint) []*botspb.EntryPoint {
	res := make([]*botspb.EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = mapEntryToPB(entry)
	}
	return res
}

func mapOptionFromPB(option *botspb.Option) types.Option {
	return types.Option{
		Text: option.Text,
		Next: int(option.Next),
	}
}

func mapOptionToPB(option types.Option) *botspb.Option {
	return &botspb.Option{
		Text: option.Text,
		Next: int32(option.Next),
	}
}

func mapOptionsFromPB(options []*botspb.Option) []types.Option {
	res := make([]types.Option, len(options))
	for i, option := range options {
		res[i] = mapOptionFromPB(option)
	}
	return res
}

func mapOptionsToPB(options []types.Option) []*botspb.Option {
	res := make([]*botspb.Option, len(options))
	for i, option := range options {
		res[i] = mapOptionToPB(option)
	}
	return res
}

func mapBlockTypeFromPB(blockType botspb.BlockType) string {
	switch blockType {
	case botspb.BlockType_BlockMessage:
		return "message"
	case botspb.BlockType_BlockQuestion:
		return "question"
	case botspb.BlockType_BlockSelection:
		return "selection"
	}
	return ""
}

func mapBlockTypeToPB(blockType string) botspb.BlockType {
	switch blockType {
	case "message":
		return botspb.BlockType_BlockMessage
	case "question":
		return botspb.BlockType_BlockQuestion
	case "selection":
		return botspb.BlockType_BlockSelection
	}
	return botspb.BlockType_Unknown
}

func mapBlockFromPB(block *botspb.Block) types.Block {
	return types.Block{
		Type:      mapBlockTypeFromPB(block.Type),
		State:     int(block.State),
		NextState: int(block.NextState),
		Options:   mapOptionsFromPB(block.Options),
		Title:     block.Title,
		Text:      block.Text,
	}
}

func mapBlockToPB(block types.Block) *botspb.Block {
	return &botspb.Block{
		Type:      mapBlockTypeToPB(block.Type),
		State:     int32(block.State),
		NextState: int32(block.NextState),
		Title:     block.Title,
		Text:      block.Text,
		Options:   mapOptionsToPB(block.Options),
	}
}

func mapBlocksFromPB(blocks []*botspb.Block) []types.Block {
	res := make([]types.Block, len(blocks))
	for i, block := range blocks {
		res[i] = mapBlockFromPB(block)
	}
	return res
}

func mapBlocksToPB(blocks []types.Block) []*botspb.Block {
	res := make([]*botspb.Block, len(blocks))
	for i, block := range blocks {
		res[i] = mapBlockToPB(block)
	}
	return res
}

func mapBotToPB(bot types.Bot) *botspb.Bot {
	return &botspb.Bot{
		BotUUID: bot.UUID,
		Name:    bot.Name,
		Token:   bot.Token,
		Entries: mapEntriesToPB(bot.Entries),
		Blocks:  mapBlocksToPB(bot.Blocks),
	}
}
