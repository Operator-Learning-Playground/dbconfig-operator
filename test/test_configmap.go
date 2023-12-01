package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/myoperator/dbconfigoperator/pkg/k8sconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	client := k8sconfig.InitClient(k8sconfig.K8sRestConfig())
	cm, _ := client.CoreV1().ConfigMaps("default").Get(context.Background(), "multi-cluster-k8s", v1.GetOptions{})
	// 连接数据库
	db, err := sql.Open("mysql", "root:1234567@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 创建库
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + "resources")
	if err != nil {
		panic(err)
	}

	// 使用数据库
	_, err = db.Exec("USE " + "resources")
	if err != nil {
		panic(err)
	}

	// 创建表
	for k, v := range cm.Data {
		fmt.Println("key: ", k)
		fmt.Println("value: ", v)
		// 直接执行建表命令
		_, err = db.Exec(v)
		if err != nil {
			panic(err)
		}
	}

}
