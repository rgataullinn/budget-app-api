package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB(connString string) {
	var err error
	Pool, err = pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to the database")
}

func CloseDB() {
	Pool.Close()
}
