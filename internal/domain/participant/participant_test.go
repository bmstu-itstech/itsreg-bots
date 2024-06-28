package participant_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"testing"
)

func TestParticipant_New(t *testing.T) {
	t.Run("successfully create a new participant", func(t *testing.T) {
		part, errs := participant.New(42, 1)
		require.Nil(t, errs)
		require.NotNil(t, part)
		require.Equal(t, int32(42), part.Id)
	})

	t.Run("fails to create a participant with invalid id", func(t *testing.T) {
		_, errs := participant.New(0, 1)
		require.ElementsMatch(t, errs, []error{participant.ErrInvalidParticipantId})
	})

	t.Run("fails to create a participant with invalid current id", func(t *testing.T) {
		_, errs := participant.New(1, 0)
		require.ElementsMatch(t, errs, []error{participant.ErrInvalidParticipantCurrent})
	})

	t.Run("fails to create a participant with invalid args", func(t *testing.T) {
		_, errs := participant.New(0, 0)
		require.ElementsMatch(t, errs, []error{
			participant.ErrInvalidParticipantId, participant.ErrInvalidParticipantCurrent})
	})
}
