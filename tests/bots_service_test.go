package tests

import (
	"github.com/bmstu-itstech/itsreg-bots/tests/suite"
	pb "github.com/bmstu-itstech/itsreg-proto/gen/go/bots"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBotsService_CreateAndProcessSuccess01(t *testing.T) {
	ctx, s := suite.New(t)

	cRes, err := s.Client.Create(ctx, &pb.CreateRequest{
		Name:  "Example bot",
		Token: "example token",
		Start: 1,
		Blocks: []*pb.Block{
			{
				State:   1,
				Type:    2,
				Default: 2,
				Title:   "Name",
				Text:    "What's your name?",
			},
			{
				State:   2,
				Type:    1,
				Default: 0,
				Title:   "End",
				Text:    "Welcome to our service!",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, cRes)
	require.NotEmpty(t, cRes.BotId)

	botId := cRes.BotId
	userId := uint64(42)

	res, err := s.Client.Process(ctx, &pb.ProcessRequest{
		BotId:  botId,
		UserId: userId,
		Text:   "/start",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Messages, 1)
	require.Empty(t, res.Messages[0].Options)
	require.Equal(t, "What's your name?", res.Messages[0].Text)

	res, err = s.Client.Process(ctx, &pb.ProcessRequest{
		BotId:  botId,
		UserId: userId,
		Text:   "Bob",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Messages, 1)
	require.Empty(t, res.Messages[0].Options)
	require.Equal(t, "Welcome to our service!", res.Messages[0].Text)
}

func TestBotsService_CreateAndProcessSuccess02(t *testing.T) {
	ctx, s := suite.New(t)

	cRes, err := s.Client.Create(ctx, &pb.CreateRequest{
		Name:  "Example bot",
		Token: "example token",
		Start: 1,
		Blocks: []*pb.Block{
			{
				State: 1,
				Type:  3,
				Title: "Pill",
				Text:  "Choose a pill?",
				Options: []*pb.BlockOption{
					{Text: "Red", Next: 2},
					{Text: "Blue", Next: 3},
				},
			},
			{
				State: 2,
				Type:  1,
				Title: "Red",
				Text:  "OK, red pill",
			},
			{
				State: 3,
				Type:  1,
				Title: "Blue",
				Text:  "OK, blue pill",
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, cRes)
	require.NotEmpty(t, cRes.BotId)

	botId := cRes.BotId
	userId := uint64(42)

	res, err := s.Client.Process(ctx, &pb.ProcessRequest{
		BotId:  botId,
		UserId: userId,
		Text:   "/start",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Messages, 1)
	require.Equal(t, "Choose a pill?", res.Messages[0].Text)
	require.NotEmpty(t, res.Messages[0].Options)
	require.Equal(t, "Red", res.Messages[0].Options[0].Text)
	require.Equal(t, "Blue", res.Messages[0].Options[1].Text)

	res, err = s.Client.Process(ctx, &pb.ProcessRequest{
		BotId:  botId,
		UserId: userId,
		Text:   "Red",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Messages, 1)
	require.Empty(t, res.Messages[0].Options)
	require.Equal(t, "OK, red pill", res.Messages[0].Text)
}
