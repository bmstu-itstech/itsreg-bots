package suite

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	botsv1 "github.com/bmstu-itstech/itsreg-proto/gen/go/bots"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client botsv1.BotsServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:      t,
		Cfg:    cfg,
		Client: botsv1.NewBotsServiceClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.Grpc.Host, strconv.Itoa(cfg.Grpc.Port))
}
