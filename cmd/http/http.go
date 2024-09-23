package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/server"
	"github.com/bmstu-itstech/itsreg-bots/internal/ports/httpport"
	"github.com/bmstu-itstech/itsreg-bots/internal/service"
)

func main() {
	app, _, _, closeFunc := service.NewApplication()
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpport.HandlerFromMux(httpport.NewHTTPServer(app), router)
	})
}
