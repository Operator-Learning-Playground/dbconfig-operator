package sysconfig

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


// InitDB 初始化db数据库
func InitDB(dsn string) *sql.DB {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

// CreateDbAndTables 创建表结构
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

