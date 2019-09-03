package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/maxence-charriere/app/pkg/app"
	"github.com/maxence-charriere/app/pkg/log"
)

func main() {
	// Setup the http handler to serve the web assembly app:
	http.Handle("/", &app.Handler{
		Name: "APP_NAME",
	})

	// Building server addr:
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)

	// Launching server:
	log.Info("starting server").T("addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("listening and serving http requests failed").
			T("reason", err).
			T("addr", addr).
			Panic()
	}
}
