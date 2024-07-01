package processor_test

import (
	"context"
	"github.com/stretchr/testify/require"
	botmemory "github.com/zhikh23/itsreg-bots/internal/domain/bot/memory"
	modulememory "github.com/zhikh23/itsreg-bots/internal/domain/module/memory"
	partmemory "github.com/zhikh23/itsreg-bots/internal/domain/participant/memory"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	"github.com/zhikh23/itsreg-bots/internal/services/processor"
	"testing"
)

func TestService_Process(t *testing.T) {
	const botId int64 = 1

	script := map[objects.State]entity.Module{
		1: {
			Node: objects.Node{
				Id:       1,
				Default:  2,
				IsSilent: false,
				Buttons:  make([]objects.Button, 0),
			},
			Title: "Module 1",
			Text:  "Module 1 text",
			BotId: botId,
		},
		2: {
			Node: objects.Node{
				Id:       2,
				Default:  2,
				IsSilent: false,
				Buttons: []objects.Button{
					{
						Text: "To 3",
						Next: 3,
					},
					{
						Text: "To 4",
						Next: 4,
					},
					{
						Text: "To 5",
						Next: 5,
					},
				},
			},
			Title: "Module 2",
			Text:  "Module 2 text",
			BotId: botId,
		},
		3: {
			Node: objects.Node{
				Id:       3,
				Default:  0,
				IsSilent: false,
				Buttons:  make([]objects.Button, 0),
			},
			Title: "Module 3",
			Text:  "Module 3 text",
			BotId: botId,
		},
		4: {
			Node: objects.Node{
				Id:       4,
				Default:  3,
				IsSilent: true,
				Buttons:  make([]objects.Button, 0),
			},
			Title: "Module 4",
			Text:  "Module 4 text",
			BotId: botId,
		},
		5: {
			Node: objects.Node{
				Id:       5,
				Default:  4,
				IsSilent: true,
				Buttons:  make([]objects.Button, 0),
			},
			Title: "Module 5",
			Text:  "Module 5 text",
			BotId: botId,
		},
	}

	modules := modulememory.New()
	for _, m := range script {
		_ = modules.Save(m)
	}

	bots := botmemory.New()
	_ = bots.Save(entity.Bot{
		Id:    1,
		Title: "Example bot",
		Token: "XXXX",
		Start: 1,
	})

	participants := partmemory.New()

	service, err := processor.New(
		processor.WithDiscardLogger(),
		processor.WithModuleRepository(modules),
		processor.WithBotRepository(bots),
		processor.WithParticipantRepository(participants))
	require.NoError(t, err)

	testUserId := int64(0)

	t.Run("start branch", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		got, err := service.Process(ctx, botId, testUserId, "Any text")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 1 text",
				Buttons: script[1].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(1), prt.State)
	})

	t.Run("unconditional branch", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         1,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "Any text")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 2 text",
				Buttons: script[2].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(2), prt.State)
	})

	t.Run("conditional branch", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         2,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "To 3")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 3 text",
				Buttons: script[3].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(3), prt.State)
	})

	t.Run("default branch", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         2,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "Any text")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 2 text",
				Buttons: script[2].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(2), prt.State)
	})

	t.Run("already finished", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         3,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "Any text")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(3), prt.State)
	})

	t.Run("silent module", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         2,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "To 4")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 4 text",
				Buttons: script[4].Buttons,
			},
			{
				Text:    "Module 3 text",
				Buttons: script[3].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(3), prt.State)
	})

	t.Run("several silent modules", func(t *testing.T) {
		testUserId++
		prtId := entity.ParticipantId{BotId: botId, UserId: testUserId}
		ctx := context.Background()

		err = participants.Save(entity.Participant{
			State:         2,
			ParticipantId: prtId,
		})
		require.NoError(t, err)

		got, err := service.Process(ctx, botId, testUserId, "To 5")
		require.NoError(t, err)
		require.ElementsMatch(t, got, []objects.Message{
			{
				Text:    "Module 5 text",
				Buttons: script[5].Buttons,
			},
			{
				Text:    "Module 4 text",
				Buttons: script[4].Buttons,
			},
			{
				Text:    "Module 3 text",
				Buttons: script[3].Buttons,
			},
		})

		prt, err := participants.Get(prtId)
		require.NoError(t, err)
		require.Equal(t, objects.State(3), prt.State)
	})
}
