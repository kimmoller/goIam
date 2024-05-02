package main

import (
	"context"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*postgres, error) {
	log.Print("Create db connections")
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Printf("unable to create connection pool: %s", err)
		}
		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func migratedb(url string) {
	m, err := migrate.New(
		"file://core/db/migrations",
		url+"?sslmode=disable")
	if err != nil {
		log.Print(err)
	}
	if err := m.Up(); err != nil {
		log.Print(err)
	}
}
