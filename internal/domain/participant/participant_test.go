package participant_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"testing"
)

func TestParticipant_New(t *testing.T) {
	t.Run("successfully create a new participant", func(t *testing.T) {
		part, err := participant.New(42)
		require.NoError(t, err)
		require.NotNil(t, part)
		require.Equal(t, int32(42), part.Id)
	})

	t.Run("fails to create a participant with invalid id", func(t *testing.T) {
		_, err := participant.New(0)
		require.ErrorIs(t, err, participant.ErrInvalidParticipantId)
	})
}
