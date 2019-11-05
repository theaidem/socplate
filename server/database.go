package main

import (
	"context"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbOnce     sync.Once
	dbInstance *sqlx.DB
)

func database() *sqlx.DB {
	dbOnce.Do(func() {
		var err error
		dataSorceName := "host=database user=socplate_usr dbname=socplate password=socplate_pass sslmode=disable"

		dbInstance, err = sqlx.Connect("postgres", dataSorceName)
		if err != nil {
			log.Fatalln(err)
		}

		err = dbInstance.PingContext(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	})

	return dbInstance
}
