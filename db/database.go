package db

import (
	"github.com/jackc/pgx"
	"io/ioutil"
)

var db *pgx.ConnPool

func init() {
	pgxConfig, _ := pgx.ParseURI("postgres://role1:12345@localhost:5432/docker")

	var err error
	db, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		})
	if err != nil {
		panic(err)
	}

	initSql, err := ioutil.ReadFile("init.sql")
	if err != nil {
		panic(err)
	}
	initString := string(initSql)

	_, err = db.Exec(initString)
	if err != nil {
		panic(err)
	}
}

func GetDB() *pgx.ConnPool {
	return db
}
