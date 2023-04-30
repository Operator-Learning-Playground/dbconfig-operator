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
	cm, _ := client.CoreV1().ConfigMaps("default").Get(context.Background(), "test-db-table", v1.GetOptions{})

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + "testdb11")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + "testdb11")
	if err != nil {
		panic(err)
	}

	for k, v := range cm.Data {
		fmt.Println("key: ", k)
		fmt.Println("value: ", v)
		_, err = db.Exec(v)
		if err != nil {
			panic(err)
		}
	}

}
