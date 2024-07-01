package memory

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"testing"
)

func TestMemory_Get(t *testing.T) {
	type testCase struct {
		name        string
		botId       int64
		expectedErr error
	}

	src := []entity.Bot{
		{
			Id:    42,
			Title: "Example bot 42",
			Token: "XXXX",
			Start: 1,
		},
		{
			Id:    37,
			Title: "Example bot 37",
			Token: "YYYY",
			Start: 1,
		},
	}

	bots := map[int64]entity.Bot{}

	for _, b := range src {
		bots[b.Id] = b
	}

	tests := []testCase{
		{
			name:        "no bot by id",
			botId:       123,
			expectedErr: bot.ErrBotNotFound,
		},
		{
			name:        "successfully get bot",
			botId:       42,
			expectedErr: nil,
		},
	}

	repos := Repository{
		bots: bots,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := repos.Get(tc.botId)
			require.Equal(t, tc.expectedErr, err, "got unexpected error")
			require.NotNil(t, m)
		})
	}
}
