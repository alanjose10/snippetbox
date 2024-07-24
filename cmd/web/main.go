package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alanjose10/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	snippetModel   *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dns := flag.String("dns", "web:password@/snippetbox?parseTime=true", "MySQL DSN")
	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate pem file")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS private key")

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

	cache, err := createTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := application{
		logger:         logger,
		snippetModel:   &models.SnippetModel{Db: db},
		templateCache:  cache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		TLSConfig: tlsConfig,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = server.ListenAndServeTLS(*tlsCert, *tlsKey)

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
