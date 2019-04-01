package db

import (
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	"io/ioutil"
	"runtime"
)

var db *pgx.ConnPool

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	connectConfig := pgx.ConnConfig{
		Host: "127.0.0.1",
		User: "role1",
		Password: "12345",
		Database: "docker",
		Port: 5432,
	}
	var err error
	db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connectConfig,
		MaxConnections: 100,
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
