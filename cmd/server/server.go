package main

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/application"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	grpcapp "github.com/bmstu-itstech/itsreg-bots/internal/presentation/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	err := application.Init(cfg)
	if err != nil {
		panic(err)
	}

	app := grpcapp.New(application.Logger, application.BotsAppService, cfg.Grpc.Port)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
}
