package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/myoperator/dbconfigoperator/pkg/k8sconfig"
	"github.com/myoperator/dbconfigoperator/pkg/sysconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	client := k8sconfig.InitClient(k8sconfig.K8sRestConfig())
	sr, err := client.CoreV1().Secrets("default").Get(context.Background(), "test-db-password", v1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

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

	for k, v := range sr.Data {
		fmt.Println("key: ", k)
		fmt.Println("value: ", string(v))

		sysconfig.CreateUser(db, "test", string(v), "testdb11")
	}

}
