package tests

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/tests/suite"
	pb "github.com/zhikh23/itsreg-proto/gen/go/bots"
	"testing"
)

func TestBotsService_CreateAndProcessSuccess(t *testing.T) {
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
	require.Empty(t, res.Messages[0].Buttons)
	require.Equal(t, "What's your name?", res.Messages[0].Text)

	res, err = s.Client.Process(ctx, &pb.ProcessRequest{
		BotId:  botId,
		UserId: userId,
		Text:   "Bob",
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Messages, 1)
	require.Empty(t, res.Messages[0].Buttons)
	require.Equal(t, "Welcome to our service!", res.Messages[0].Text)
}
