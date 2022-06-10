package db

import (
	"github.com/jackc/pgx"
)

var db *pgx.ConnPool

func init() {
	pgxConfig, _ := pgx.ParseURI("postgresql://kotyarich:1234@databasestat:5432/postgres")
	pgxConfig.RuntimeParams["timezone"] = "Europe/Moscow"
	var err error
	db, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		})
	if err != nil {
		panic(err)
	}

	//initSql, err := ioutil.ReadFile("init.sql")
	//if err != nil {
	//	panic(err)
	//}
	//initString := string(initSql)
	
	//_, err = db.Exec(initString)
	//if err != nil {
	//	panic(err)
	//}
}

func GetDB() *pgx.ConnPool {
	return db
}
