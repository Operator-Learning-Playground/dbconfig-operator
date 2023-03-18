package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 创建db
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + "testdb")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + "testdb")
	if err != nil {
		panic(err)
	}

	// 创建表
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS example ( id integer, data varchar(32) )")
	//_, err = db.Exec(tablesInfo)
	if err != nil {
		panic(err)
	}

	// 删除db。
	//_, err = db.Exec("DROP DATABASE IF EXISTS " + "testdb")
	//if err != nil {
	//	panic(err)
	//}

}
