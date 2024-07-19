package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alanjose10/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger       *slog.Logger
	snippetModel *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dns := flag.String("dns", "web:password@/snippetbox?parseTime=true", "MySQL DSN")
	flag.Parse()

	// Adding a structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db, err := openDb(*dns)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := application{
		logger:       logger,
		snippetModel: &models.SnippetModel{Db: db},
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Successfully connected to MySQL database.")

	return db, nil
}
