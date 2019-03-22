package db

import (
	"database/sql"
	"io/ioutil"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	connStr := "host=127.0.0.1 user=role1 password='12345' dbname=docker sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
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

func GetDB() *sql.DB {
	return db
}
