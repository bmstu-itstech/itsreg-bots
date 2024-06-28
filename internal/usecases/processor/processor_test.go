package processor_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/module"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"github.com/zhikh23/itsreg-bots/internal/usecases/processor"
	"github.com/zhikh23/itsreg-bots/internal/usecases/processor/mocks"
	"log/slog"
	"os"
	"testing"
)

func TestProcessor_Process(t *testing.T) {
	store := mocks.NewProviderMock()
	store.Register(&module.Module{
		Id:       1,
		Title:    "Module 1",
		Text:     "Text of module 1",
		IsSilent: false,
		Type:     module.String,
		Next:     2,
		Buttons:  nil,
	})
	store.Register(&module.Module{
		Id:       2,
		Title:    "Module 2",
		Text:     "Text of module 2",
		IsSilent: false,
		Type:     module.String,
		Next:     0,
		Buttons:  nil,
	})
	store.Register(&module.Module{
		Id:       3,
		Title:    "Module 3",
		Text:     "Text of module 3",
		IsSilent: false,
		Type:     module.String,
		Next:     3,
		Buttons: []module.Button{
			{
				Text:  "Text of button 3.1",
				Value: "Option 1",
				Next:  1,
			},
			{
				Text:  "Text of button 3.2",
				Value: "Option 2",
				Next:  2,
			},
			{
				Text:  "Text of button 3.3",
				Value: "Option 3",
				Next:  4,
			},
			{
				Text:  "Text of button 3.4",
				Value: "Option 4",
				Next:  5,
			},
		},
	})
	store.Register(&module.Module{
		Id:       4,
		Title:    "Module 4",
		Text:     "Text of module 4",
		IsSilent: true,
		Type:     module.String,
		Next:     2,
		Buttons:  nil,
	})
	store.Register(&module.Module{
		Id:       5,
		Title:    "Module 5",
		Text:     "Text of module 5",
		IsSilent: true,
		Type:     module.String,
		Next:     4,
		Buttons:  nil,
	})

	proc := processor.New(
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		store)

	t.Run("unconditional branch", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 1,
		}

		ans, err := proc.Process(part, "")
		require.NoError(t, err)
		require.ElementsMatch(t, ans, []string{"Text of module 2"})
		require.Equal(t, int32(2), part.Current)
	})

	t.Run("one of conditional branch", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 3,
		}

		ans, err := proc.Process(part, "Option 2")
		require.NoError(t, err)
		require.ElementsMatch(t, ans, []string{"Text of module 2"})
		require.Equal(t, int32(2), part.Current)
	})

	t.Run("default branch", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 3,
		}

		ans, err := proc.Process(part, "Default Option")
		require.NoError(t, err)
		require.ElementsMatch(t, ans, []string{"Text of module 3"})
		require.Equal(t, int32(3), part.Current)
	})

	t.Run("finish branch", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 2,
		}

		_, err := proc.Process(part, "Something")
		require.ErrorIs(t, err, processor.ErrAlreadyFinished)
	})

	t.Run("with silent module", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 3,
		}

		ans, err := proc.Process(part, "Option 3")
		require.NoError(t, err)
		require.Equal(t, ans, []string{"Text of module 4", "Text of module 2"})
		require.Equal(t, int32(2), part.Current)
	})

	t.Run("with several silent modules", func(t *testing.T) {
		part := &participant.Participant{
			Id:      42,
			Current: 3,
		}

		ans, err := proc.Process(part, "Option 4")
		require.NoError(t, err)
		require.Equal(t, ans, []string{"Text of module 5", "Text of module 4", "Text of module 2"})
		require.Equal(t, int32(2), part.Current)
	})
}
