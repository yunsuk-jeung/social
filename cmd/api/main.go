package main

import (
	"log"

	"github.com/yunsuk-jeung/social/internal/db"
	"github.com/yunsuk-jeung/social/internal/env"
	"github.com/yunsuk-jeung/social/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdelConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdelConns,
		cfg.db.maxIdelTime,
	)

	if err != nil {
		log.Panic(err)
	}

	store := store.NewStorage(db)
	defer db.Close()
	log.Printf("database connection pool established")

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
