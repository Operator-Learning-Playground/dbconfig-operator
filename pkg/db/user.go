package db

import (
	"fmt"
	"k8s.io/klog/v2"
)

// CreateUser 创建 user
func (gb *GlobalDB) CreateUser(user string, password string, dbname string) error {

	// 创建user用户
	key := fmt.Sprintf("CREATE USER IF NOT EXISTS "+user+" IDENTIFIED BY '%v';", password)
	_, err := gb.DB.Exec(key)
	klog.Info(key)
	if err != nil {
		klog.Error("create user error: ", err)
		return err
	}
	// 授权
	key1 := fmt.Sprintf("GRANT ALL PRIVILEGES ON "+dbname+".* TO "+"'%v'@'%%';", user)
	klog.Info(key1)
	_, err = gb.DB.Exec(key1)
	if err != nil {
		klog.Error("set db privileges error: ", err)
		return err
	}
	return nil
}

// DeleteUsers 删除 user
func (gb *GlobalDB) DeleteUsers(users []string) {
	if len(users) == 0 {
		klog.Info("no user have to delete")
		return
	}
	for _, user := range users {
		_, err := gb.DB.Exec("DROP USER IF EXISTS " + user)
		if err != nil {
			klog.Error("drop user error: ", err)
		}
	}
}
