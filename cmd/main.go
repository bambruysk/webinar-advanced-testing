package main

import (
	"io"
	"log"
	"time"

	"webinar-testing/internal/api"
	"webinar-testing/internal/service"
	"webinar-testing/internal/storage/postgres"
)

func main() {
	var closers []io.Closer
	defer func() {
		for _, closer := range closers {
			if err := closer.Close(); err != nil {
				log.Println(err)
			}
		}
	}()

	pgCfg := &postgres.Config{
		Host:             "localhost",
		Port:             5432,
		ConnectTimeout:   10 * time.Second,
		QueryTimeout:     5 * time.Second,
		Username:         "postgres",
		Password:         "testpswd",
		DBName:           "postgres",
		MigrationVersion: 1,
	}

	db, err := postgres.New(pgCfg)
	if err != nil {
		panic(err)
	}

	closers = append(closers, db)

	serv := service.New(db)

	server := api.NewServer(serv)

	if err := server.Run(); err != nil {
		log.Println(err)
		return
	}

}
