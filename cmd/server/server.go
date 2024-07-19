package main

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/application"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	grpcapp "github.com/bmstu-itstech/itsreg-bots/internal/presentation/grpc"
	"github.com/bmstu-itstech/itsreg-bots/pkg/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal, envDev:
		log = slog.New(
			slogpretty.PrettyHandlerOptions{
				SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
			}.NewPrettyHandler(os.Stdout),
		)
	case envProd:
		log = slog.New(
			slogpretty.PrettyHandlerOptions{
				SlogOpts: &slog.HandlerOptions{Level: slog.LevelInfo},
			}.NewPrettyHandler(os.Stdout),
		)
	}

	return log
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	app, err := application.New(log, cfg)
	if err != nil {
		panic(err)
	}

	server := grpcapp.New(log, app, cfg.Grpc.Port)

	go func() {
		server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	server.Stop()
}
