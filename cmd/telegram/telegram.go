package main

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/ports/telegram"
	"github.com/bmstu-itstech/itsreg-bots/internal/service"
)

func main() {
	app, msgCh, runCh, closeFunc := service.NewApplication()
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	telegram.RunTelegramPort(app, msgCh, runCh)
}
