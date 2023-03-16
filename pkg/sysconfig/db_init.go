package sysconfig

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)



func InitDB(dsn string) *sql.DB {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return db
}

func CreateDbAndTables(db *sql.DB, dbname string, tableInfos []string) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		panic(err)
	}

	for _, tableInfo := range tableInfos {
		_, err = db.Exec(tableInfo)
		if err != nil {
			panic(err)
		}
	}

}

