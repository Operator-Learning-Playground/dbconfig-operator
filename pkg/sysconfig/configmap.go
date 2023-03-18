package sysconfig

import (
	"context"
	"github.com/myoperator/dbconfigoperator/pkg/k8sconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// GetConfigmapData 取得用户自定义的configmap data字段
func GetConfigmapData(name string, namespace string) ([]string, error) {

	res := make([]string, 0)

	client := k8sconfig.InitClient(k8sconfig.K8sRestConfig())
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		klog.Error("get configmap error!", err)
		return res, err
	}
	for _, v := range cm.Data {
		res = append(res, v)
	}
	return res, nil
}
