package main

import (
	"google.golang.org/grpc"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/ports/grpcport"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/server"
)

func main() {
	app, closeFunc := service.NewApplication()
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	server.RunGRPCServer(func(server *grpc.Server) {
		grpcport.RegisterGRPCServer(server, app)
	})
}
