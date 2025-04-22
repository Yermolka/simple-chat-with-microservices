package repository

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func newDb() *sqlx.DB {
	conn_str := ":memory:"
	
	is_debug, ok := os.LookupEnv("DEBUG")
	if !ok || is_debug == "false"{
		
	}

	db, err := sqlx.Connect("sqlite3", conn_str)
	if err != nil {
		panic(err)
	}
	return db
}
