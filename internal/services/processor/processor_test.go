package processor_test

import (
	"github.com/stretchr/testify/require"
	botmemory "github.com/zhikh23/itsreg-bots/internal/domain/bot/memory"
	modulememory "github.com/zhikh23/itsreg-bots/internal/domain/module/memory"
	partmemory "github.com/zhikh23/itsreg-bots/internal/domain/participant/memory"
	sendrecorder "github.com/zhikh23/itsreg-bots/internal/domain/sender/recorder"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/services/processor"
	"testing"
)

func TestService_Process(t *testing.T) {
	script := []entity.Module{
		{
			Node: entity.Node{
				Id:       1,
				Default:  2,
				IsSilent: false,
				Buttons:  make([]entity.Button, 0),
			},
			Title: "Module 1",
			Text:  "Module 1 text",
			BotId: 1,
		},
		{
			Node: entity.Node{
				Id:       2,
				Default:  2,
				IsSilent: false,
				Buttons: []entity.Button{
					{
						Text:   "To 3",
						NextId: 3,
					},
					{
						Text:   "To 4",
						NextId: 4,
					},
					{
						Text:   "To 5",
						NextId: 5,
					},
				},
			},
			Title: "Module 2",
			Text:  "Module 2 text",
			BotId: 1,
		},
		{
			Node: entity.Node{
				Id:       3,
				Default:  0,
				IsSilent: false,
				Buttons:  make([]entity.Button, 0),
			},
			Title: "Module 3",
			Text:  "Module 3 text",
			BotId: 1,
		},
		{
			Node: entity.Node{
				Id:       4,
				Default:  3,
				IsSilent: true,
				Buttons:  make([]entity.Button, 0),
			},
			Title: "Module 4",
			Text:  "Module 4 text",
			BotId: 1,
		},
		{
			Node: entity.Node{
				Id:       5,
				Default:  4,
				IsSilent: true,
				Buttons:  make([]entity.Button, 0),
			},
			Title: "Module 5",
			Text:  "Module 5 text",
			BotId: 1,
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

	recorder := sendrecorder.New()

	service, err := processor.New(
		processor.WithModuleRepository(modules),
		processor.WithBotRepository(bots),
		processor.WithParticipantRepository(participants),
		processor.WithSender(recorder))
	require.NoError(t, err)

	testUserId := int64(0)

	t.Run("start branch", func(t *testing.T) {
		testUserId++
		err = service.Process(1, testUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 1 text",
			},
		}, recorder.GetLastRecords())
	})

	t.Run("unconditional branch", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 1,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 2 text",
			},
		}, recorder.GetLastRecords())
	})

	t.Run("conditional branch", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 2,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "To 3")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 3 text",
			},
		}, recorder.GetLastRecords())
	})

	t.Run("default branch", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 2,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 2 text",
			},
		}, recorder.GetLastRecords())
	})

	t.Run("already finished", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 3,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "Any text")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{}, recorder.GetLastRecords()) // Empty
	})

	t.Run("silent module", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 2,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "To 4")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 4 text",
			},
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 3 text",
			},
		}, recorder.GetLastRecords())
	})

	t.Run("several silent modules", func(t *testing.T) {
		testUserId++
		err = participants.Save(entity.Participant{
			CurrentId: 2,
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: testUserId,
			},
		})
		require.NoError(t, err)

		err = service.Process(1, testUserId, "To 5")
		require.NoError(t, err)
		require.Equal(t, []sendrecorder.Record{
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 5 text",
			},
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 4 text",
			},
			{
				Receiver: entity.ParticipantId{BotId: 1, UserId: testUserId},
				Text:     "Module 3 text",
			},
		}, recorder.GetLastRecords())
	})
}
