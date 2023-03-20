package sysconfig

import (
	"context"
	"github.com/myoperator/dbconfigoperator/pkg/k8sconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const secretKey = "PASSWORD"

func GetSecretData(name string, namespace string) (string, error) {

	var res string
	client := k8sconfig.InitClient(k8sconfig.K8sRestConfig())
	sr, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		klog.Error("get secret error!", err)
		return res, err
	}

	res = string(sr.Data[secretKey])

	return res, nil

}
