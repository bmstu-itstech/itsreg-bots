package memory

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	"testing"
)

func TestMemory_Get(t *testing.T) {
	type testCase struct {
		name        string
		botId       int64
		nodeId      objects.NodeId
		expectedErr error
	}

	src := []entity.Module{
		{
			BotId: 42,
			Title: "Example title 1",
			Text:  "Lorem ispum 1",
			Node: objects.Node{
				Id:       61,
				Default:  62,
				IsSilent: false,
				Buttons:  make([]objects.Button, 0),
			},
		},
		{
			BotId: 33,
			Title: "Example title 2",
			Text:  "Lorem ispum 2",
			Node: objects.Node{
				Id:       61,
				Default:  62,
				IsSilent: false,
				Buttons:  make([]objects.Button, 0),
			},
		},
	}

	modules := map[pair]entity.Module{}

	for _, m := range src {
		modules[pair{m.BotId, m.Node.Id}] = m
	}

	tests := []testCase{
		{
			name:        "no module by id",
			botId:       5,
			nodeId:      61,
			expectedErr: module.ErrModuleNotFound,
		},
		{
			name:        "no module by id",
			botId:       42,
			nodeId:      1,
			expectedErr: module.ErrModuleNotFound,
		},
		{
			name:        "successfully get module by id",
			botId:       33,
			nodeId:      61,
			expectedErr: nil,
		},
	}

	repos := Repository{
		modules: modules,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := repos.Get(tc.botId, tc.nodeId)
			require.Equal(t, tc.expectedErr, err, "got unexpected error")
			require.NotNil(t, m)
		})
	}
}
