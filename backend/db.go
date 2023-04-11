package main

import (
	"database/sql"
	"flag"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitDb() *bun.DB {
	dbConnectionStr := flag.String("db", "", "Database connection string")
	flag.Parse()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(*dbConnectionStr)))

	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}
