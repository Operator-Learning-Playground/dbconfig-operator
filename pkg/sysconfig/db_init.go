package sysconfig

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/klog/v2"
)


// InitDB 初始化db数据库
func InitDB(dsn string) *sql.DB {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		klog.Error("open mysql error: ", err)
	}

	return db
}

// CreateDbAndTables 创建表结构
func CreateDbAndTables(db *sql.DB, dbname string, tableInfos []string) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		klog.Error("create databases error: ", err)
	}

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		klog.Error("use databases error: ", err)
	}

	for _, tableInfo := range tableInfos {
		_, err = db.Exec(tableInfo)
		if err != nil {
			klog.Error("database: ", dbname, ", create tables error: ", err)
		}
	}

}

// DeleteDB 删除db
func DeleteDBs(db *sql.DB, dbnames []string) {

	if len(dbnames) == 0 {
		klog.Info("no db have to delete")
		return
	}
	for _, dbname := range dbnames {
		_, err := db.Exec("DROP DATABASE IF EXISTS " + dbname)
		if err != nil {
			klog.Error("drop databases error: ", err)
		}
	}


}
