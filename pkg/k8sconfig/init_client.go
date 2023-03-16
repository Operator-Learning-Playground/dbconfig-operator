package k8sconfig

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

// 返回初始化k8s-client
func InitClient(config *rest.Config) kubernetes.Interface {
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
