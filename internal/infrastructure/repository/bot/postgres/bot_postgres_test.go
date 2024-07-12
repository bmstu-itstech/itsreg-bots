package postgres

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/stretchr/testify/require"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"testing"
)

const (
	dbName = "test_bot_postgres"
	dbUser = "postgres"
	dbPass = "postgres"
)

func SetupContainer() *pgcontainer.PostgresContainer {
	pgContainer, err := pgcontainer.Run(
		context.Background(),
		"docker.io/postgres:16-alpine",
		pgcontainer.WithDatabase(dbName),
		pgcontainer.WithUsername(dbUser),
		pgcontainer.WithPassword(dbPass),
	)
	if err != nil {
		log.Fatal(err)
	}

	return pgContainer
}

func TestBotPostgresRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	container := SetupContainer()
	defer func(container *pgcontainer.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(container, context.Background())

	repos, err := NewPostgresBotRepository(container.MustConnectionString(context.Background()))
	require.NoError(t, err)

	t.Run("should save the bot", func(t *testing.T) {
		ctx := context.Background()

		bot, err := entity.NewBot(
			value.UnknownBotId, "Bot #1", "AAA", value.State(1))
		require.NoError(t, err)

		id, err := repos.Save(ctx, bot)
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})
}
