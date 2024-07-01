package recorder_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/sender/recorder"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"testing"
)

func TestRecorder(t *testing.T) {
	r := recorder.New()

	err := r.SendMessage(entity.ParticipantId{
		BotId:  25,
		UserId: 12,
	}, "Hello, user!", nil)
	require.NoError(t, err)

	err = r.SendMessage(entity.ParticipantId{
		BotId:  25,
		UserId: 12,
	}, "Hello, user, again!", nil)
	require.NoError(t, err)

	expected := []recorder.Record{
		{
			Receiver: entity.ParticipantId{
				BotId:  25,
				UserId: 12,
			},
			Text: "Hello, user!",
		},
		{
			Receiver: entity.ParticipantId{
				BotId:  25,
				UserId: 12,
			},
			Text: "Hello, user, again!",
		},
	}

	require.ElementsMatch(t, expected, r.GetLastRecords())
	require.ElementsMatch(t, []recorder.Record{}, r.GetLastRecords())
}
