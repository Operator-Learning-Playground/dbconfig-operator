package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:1234567@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("USE " + "testdb")
	rows, err := db.Query("SHOW TABLES")

	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return
		}
		fmt.Println(table)
	}
	return

	//// 创建db
	//_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + "testdb")
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = db.Exec("USE " + "testdb")
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 创建表
	//_, err = db.Exec("CREATE TABLE IF NOT EXISTS example ( id integer, data varchar(32) )")
	////_, err = db.Exec(tablesInfo)
	//if err != nil {
	//	panic(err)
	//}

	// 删除db。
	//_, err = db.Exec("DROP DATABASE IF EXISTS " + "testdb")
	//if err != nil {
	//	panic(err)
	//}

}
