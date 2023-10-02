package sysconfig

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/klog/v2"
	"time"
)

// InitDB 初始化db数据库
func InitDB(cfg *SysConfig) (*sql.DB, error) {

	db, err := sql.Open("mysql", cfg.Dsn)
	if err != nil {
		klog.Error("open mysql error: ", err)
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(30 * time.Second)

	return db, nil
}

// CheckOrCreateDb 检查或创建库
func CheckOrCreateDb(db *sql.DB, dbname string) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		klog.Error("create databases error: ", err)
	}
}

// CreateTable 创建Tables
func CreateTable(db *sql.DB, dbname string, tableInfo string) {

	_, err := db.Exec("USE " + dbname)
	if err != nil {
		klog.Error("use databases error: ", err)
	}

	_, err = db.Exec(tableInfo)
	if err != nil {
		klog.Error("database: ", dbname, ", create tables error: ", err)
	}
}

// DeleteDBs 删除db
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

// CheckTableIsExists 检查是否存在表
func CheckTableIsExists(db *sql.DB, dbName string, tableName string) (bool, error) {

	_, err := db.Exec("USE " + dbName)
	if err != nil {
		return false, err
	}

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return false, err
		}
		// 如果表中有发现，返回true与空字符串
		if table == tableName {
			return true, nil
		}
	}
	// 检查出db库中没有表的情况，返回false, 表名, nil
	return false, nil
}
