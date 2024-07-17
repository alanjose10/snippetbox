package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	addr := flag.String("addr", ":1234", "HTTP network address")
	flag.Parse()

	// Adding a structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	app := application{
		logger: logger,
	}

	// Create a new server mux
	// and register the home function as handler for /
	mux := http.NewServeMux()

	fileServe := http.FileServer(http.Dir("./ui/static"))

	// File serve route
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServe))

	mux.HandleFunc("GET /{$}", app.homeGet)

	mux.HandleFunc("GET /snippet/view/{id}", app.snippetViewGet)

	mux.HandleFunc("GET /snippet/create", app.snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
