package sysconfig

import (
	"database/sql"
	"fmt"
	"k8s.io/klog/v2"
)

func CreateUser(db *sql.DB, user string, password string, dbname string) {

	// 创建user用户
	key := fmt.Sprintf("CREATE USER IF NOT EXISTS " + user + " IDENTIFIED BY '%v';", password)
	_, err := db.Exec(key)
	fmt.Println(key)
	if err != nil {
		klog.Error("create user error: ", err)
	}
	// 授权
	key1 := fmt.Sprintf("GRANT ALL PRIVILEGES ON " + dbname +  ".* TO "+  "'%v'@'%%';", user)
	fmt.Println(key1)
	_, err = db.Exec(key1)
	if err != nil {
		klog.Error("set db privileges error: ", err)
	}
}