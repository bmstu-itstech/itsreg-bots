package service_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/service"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	ansmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/answer/memory"
	blockmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/block/memory"
	botmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/bot/memory"
	prtmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/participant/memory"
	"github.com/zhikh23/itsreg-bots/pkg/logger/handlers/slogdiscard"
	"testing"
)

func TestService_Process(t *testing.T) {
	var script []*entity.Block
	ctx := context.Background()

	botRepos := botmemory.NewMemoryBotRepository()
	bot, err := entity.NewBot(value.UnknownBotId, "Example bot", "Example token", value.State(1))
	require.NoError(t, err)
	err = botRepos.Save(ctx, bot)
	require.NoError(t, err)

	node, err := value.NewQuestionNode(value.State(1), value.State(2))
	require.NoError(t, err)
	block, err := entity.NewBlock(node, bot.Id, "Question 1", "Block 1")
	require.NoError(t, err)
	script = append(script, block)

	node, err = value.NewSelectionNode(value.State(2), []value.Option{
		{Text: "To 3", Next: value.State(3)},
		{Text: "To 4", Next: value.State(4)},
		{Text: "To 5", Next: value.State(5)},
	})
	require.NoError(t, err)
	block, err = entity.NewBlock(node, bot.Id, "Selection 2", "Block 2")
	require.NoError(t, err)
	script = append(script, block)

	node, err = value.NewMessageNode(value.State(3), value.StateNone)
	require.NoError(t, err)
	block, err = entity.NewBlock(node, bot.Id, "Message 3", "Block 3")
	require.NoError(t, err)
	script = append(script, block)

	node, err = value.NewMessageNode(value.State(4), value.State(3))
	require.NoError(t, err)
	block, err = entity.NewBlock(node, bot.Id, "Message 4", "Block 4")
	require.NoError(t, err)
	script = append(script, block)

	node, err = value.NewQuestionNode(value.State(5), value.StateNone)
	require.NoError(t, err)
	block, err = entity.NewBlock(node, bot.Id, "Question 5", "Block 5")
	require.NoError(t, err)
	script = append(script, block)

	blockRepos := blockmemory.NewMemoryBlockRepository()
	for _, b := range script {
		_ = blockRepos.Save(ctx, b)
	}

	prtRepos := prtmemory.NewMemoryParticipantRepository()

	ansRepos := ansmemory.NewMemoryAnswerRepository()

	proc := service.NewProcessor(
		slogdiscard.NewDiscardLogger(),
		ansRepos, blockRepos, botRepos, prtRepos)

	tempUserId := value.UserId(0)

	t.Run("should process start block", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		res, err := proc.Process(ctx, bot.Id, tempUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []service.Message{
			{Text: "Block 1", Options: []string{}},
		}, res)

		prt, err := prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.State(1), prt.Current)
	})

	t.Run("should process question block", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		prt, err := entity.NewParticipant(prtId, value.State(1))
		require.NoError(t, err)
		err = prtRepos.Save(ctx, prt)
		require.NoError(t, err)

		res, err := proc.Process(ctx, bot.Id, tempUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []service.Message{
			{Text: "Block 2", Options: []string{
				"To 3", "To 4", "To 5",
			}},
		}, res)

		prt, err = prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.State(2), prt.Current)

	})

	t.Run("should process selection block", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		prt, err := entity.NewParticipant(prtId, value.State(2))
		require.NoError(t, err)
		err = prtRepos.Save(ctx, prt)
		require.NoError(t, err)

		res, err := proc.Process(ctx, bot.Id, tempUserId, "To 5")
		require.NoError(t, err)
		require.Equal(t, []service.Message{
			{Text: "Block 5", Options: []string{}},
		}, res)

		prt, err = prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.State(5), prt.Current)
	})

	t.Run("should process selection and message block", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		prt, err := entity.NewParticipant(prtId, value.State(2))
		require.NoError(t, err)
		err = prtRepos.Save(ctx, prt)
		require.NoError(t, err)

		res, err := proc.Process(ctx, bot.Id, tempUserId, "To 4")
		require.NoError(t, err)
		require.Equal(t, []service.Message{
			{Text: "Block 4", Options: []string{}},
			{Text: "Block 3", Options: []string{}},
		}, res)

		prt, err = prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.StateNone, prt.Current)
	})

	t.Run("should return error when participant already finished", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		prt, err := entity.NewParticipant(prtId, value.State(9))
		require.NoError(t, err)
		err = prtRepos.Save(ctx, prt)
		require.NoError(t, err)
		err = prtRepos.UpdateState(ctx, prtId, value.StateNone)
		require.NoError(t, err)

		res, err := proc.Process(ctx, bot.Id, tempUserId, "Any text")
		require.NoError(t, err)
		require.Empty(t, res)

		prt, err = prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.StateNone, prt.Current)
	})

	t.Run("should answer an error when participant does not choose an option", func(t *testing.T) {
		tempUserId++
		prtId := value.ParticipantId{BotId: bot.Id, UserId: tempUserId}
		ctx := context.Background()

		prt, err := entity.NewParticipant(prtId, value.State(2))
		require.NoError(t, err)
		err = prtRepos.Save(ctx, prt)
		require.NoError(t, err)

		res, err := proc.Process(ctx, bot.Id, tempUserId, "Any text")
		require.NoError(t, err)
		require.NotEmpty(t, res)

		prt, err = prtRepos.Participant(ctx, prtId)
		require.NoError(t, err)
		require.Equal(t, value.State(2), prt.Current)
	})
}
