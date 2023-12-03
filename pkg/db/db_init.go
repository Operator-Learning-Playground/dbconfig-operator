package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"k8s.io/klog/v2"
	"time"
)

type GlobalDB struct {
	DB *sql.DB
}

func newGlobalDB() *GlobalDB {
	return &GlobalDB{}
}

var InitGlobalDB *GlobalDB

func init() {
	InitGlobalDB = newGlobalDB()
}

// TODO: 抽象出 db 结构体，让 func 变成 struct 的方法

// InitDB 初始化db数据库
func InitDB(dbconfig *dbconfigv1alpha1.DbConfig) (*GlobalDB, error) {

	db, err := sql.Open("mysql", dbconfig.Spec.Dsn)
	if err != nil {
		klog.Error("open mysql error: ", err)
		return nil, err
	}
	db.SetMaxIdleConns(dbconfig.Spec.MaxIdleConn)
	db.SetMaxOpenConns(dbconfig.Spec.MaxOpenConn)
	db.SetConnMaxLifetime(30 * time.Second)

	InitGlobalDB.DB = db

	return InitGlobalDB, nil
}

// CheckOrCreateDb 检查或创建库
func (gb *GlobalDB) CheckOrCreateDb(dbname string) {
	_, err := gb.DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		klog.Error("create databases error: ", err)
	}
}

// CreateTable 创建 Tables
func (gb *GlobalDB) CreateTable(dbname string, tableInfo string) {

	_, err := gb.DB.Exec("USE " + dbname)
	if err != nil {
		klog.Error("use databases error: ", err)
	}

	_, err = gb.DB.Exec(tableInfo)
	if err != nil {
		klog.Error("database: ", dbname, ", create tables error: ", err)
	}
}

// DeleteDBs 删除传入的 db 库
func (gb *GlobalDB) DeleteDBs(dbnames []string) {

	if len(dbnames) == 0 {
		klog.Info("no db have to delete")
		return
	}
	for _, dbname := range dbnames {
		_, err := gb.DB.Exec("DROP DATABASE IF EXISTS " + dbname)
		if err != nil {
			klog.Error("drop databases error: ", err)
		}
	}
}

// CheckTableIsExists 检查是否存在表
func (gb *GlobalDB) CheckTableIsExists(dbName string, tableName string) (bool, error) {

	_, err := gb.DB.Exec("USE " + dbName)
	if err != nil {
		return false, err
	}

	rows, err := gb.DB.Query("SHOW TABLES")
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return false, err
		}
		// 如果表中有发现，返回 true 与空字符串
		if table == tableName {
			return true, nil
		}
	}
	// 检查出 db 库中没有表的情况，返回 false, 表名, nil
	return false, nil
}
