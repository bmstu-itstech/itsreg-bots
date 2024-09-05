package main

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/common/server"
	"net/http"

	"github.com/bmstu-itstech/itsreg-bots/internal/ports/httpport"
	"github.com/bmstu-itstech/itsreg-bots/internal/service"
	"github.com/go-chi/chi/v5"
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
