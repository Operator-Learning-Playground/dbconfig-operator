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

// DeleteUser 删除user
func DeleteUsers(db *sql.DB, users []string) {
	if len(users) == 0 {
		klog.Info("no user have to delete")
		return
	}
	for _, user := range users {
		_, err := db.Exec("DROP USER IF EXISTS " + user)
		if err != nil {
			klog.Error("drop user error: ", err)
		}
	}


}