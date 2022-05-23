package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/pkg/logging"
	"log"
	"os"
)

func main() {
	log.Println("reading environment")
	cfg := config.GetConfig()

	log.Println("logging initializing")
	logger := logging.InitLogger(cfg.IsDevelopment, cfg.AppConfig.LogLevel)
	logger.Info("xded")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database)

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	var greeting string

	err = db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
