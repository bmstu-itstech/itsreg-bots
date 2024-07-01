package memory

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"testing"
)

func TestMemory_Get(t *testing.T) {
	type testCase struct {
		name        string
		prtId       entity.ParticipantId
		expectedErr error
	}

	src := []entity.Participant{
		{
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: 42,
			},
			State: 1,
		},
		{
			ParticipantId: entity.ParticipantId{
				BotId:  14,
				UserId: 42,
			},
			State: 1,
		},
		{
			ParticipantId: entity.ParticipantId{
				BotId:  1,
				UserId: 9,
			},
			State: 1,
		},
	}

	participants := map[entity.ParticipantId]entity.Participant{}

	for _, p := range src {
		participants[p.ParticipantId] = p
	}

	tests := []testCase{
		{
			name: "no participant found",
			prtId: entity.ParticipantId{
				BotId:  123,
				UserId: 123,
			},
			expectedErr: participant.ErrParticipantNotFound,
		},
		{
			name: "successfully get participant",
			prtId: entity.ParticipantId{
				BotId:  14,
				UserId: 42,
			},
			expectedErr: nil,
		},
		{
			name: "successfully get participant",
			prtId: entity.ParticipantId{
				BotId:  1,
				UserId: 9,
			},
			expectedErr: nil,
		},
	}

	repos := Repository{
		participants: participants,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := repos.Get(tc.prtId)
			require.Equal(t, tc.expectedErr, err, "got unexpected error")
			require.NotNil(t, m)
		})
	}
}
