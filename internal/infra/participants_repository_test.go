package infra_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/infra"
)

var (
	randomBotUUID = gofakeit.UUID()
)

func TestPgParticipantsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)
	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	err := setupDBParticipants(context.Background(), db)
	require.NoError(t, err)

	repos := infra.NewPgParticipantsRepository(db)
	testParticipantsRepository(t, repos)
}

func testParticipantsRepository(t *testing.T, repos bots.ParticipantRepository) {
	t.Parallel()

	t.Run("should create participant", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		userID := gofakeit.Int64()
		err := repos.UpdateOrCreate(ctx, randomBotUUID, userID, func(ctx context.Context, prt *bots.Participant) error {
			prt.SwitchTo(1)
			return prt.AddAnswer("answer")
		})
		require.NoError(t, err)

		participants, err := repos.ParticipantsOfBot(ctx, randomBotUUID)
		require.NoError(t, err)

		expected := bots.MustNewParticipant(randomBotUUID, userID)
		expected.SwitchTo(1)
		require.NoError(t, expected.AddAnswer("answer"))
		require.Contains(t, participants, expected)
	})

	t.Run("should update participant", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		userID := gofakeit.Int64()
		err := repos.UpdateOrCreate(ctx, randomBotUUID, userID, func(ctx context.Context, prt *bots.Participant) error {
			prt.SwitchTo(1)
			return nil
		})
		require.NoError(t, err)

		err = repos.UpdateOrCreate(ctx, randomBotUUID, userID, func(ctx context.Context, prt *bots.Participant) error {
			require.Equal(t, 1, prt.State)
			prt.SwitchTo(2)
			return nil
		})
		require.NoError(t, err)

		participants, err := repos.ParticipantsOfBot(ctx, randomBotUUID)
		require.NoError(t, err)

		expected := bots.MustNewParticipant(randomBotUUID, userID)
		expected.SwitchTo(2)
		require.Contains(t, participants, expected)
	})
}

func setupDBParticipants(ctx context.Context, db *sqlx.DB) error {
	return pgutils.RunTx(ctx, db, func(tx *sqlx.Tx) error {
		_, err := pgutils.Exec(ctx, tx,
			`INSERT INTO 
				bots (uuid, name, token, status, created_at, updated_at) 
			VALUES 
				($1, $2, $3, $4, $5, $6)`,
			randomBotUUID, gofakeit.Name(), gofakeit.UUID(), "stopped", time.Now(), time.Now(),
		)
		if err != nil {
			return err
		}

		_, err = pgutils.Exec(ctx, tx,
			`INSERT INTO
				blocks (bot_uuid, state, type, next_state, title, text)
			VALUES 
				($1, $2, $3, $4, $5, $6),
				($7, $8, $9, $10, $11, $12)`,
			randomBotUUID, 1, "question", nil, "Question 1", "Some text",
			randomBotUUID, 2, "question", nil, "Question 2", "Some text",
		)
		if err != nil {
			return err
		}

		return nil
	})
}
