package main

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/ports"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service"
)

func main() {
	app, closeFn := service.NewApplication()
	defer func() {
		err := closeFn()
		if err != nil {
			panic(err)
		}
	}()

	server := ports.NewTelegramServer(app)
}
