package controller

import (
	"context"
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"github.com/myoperator/dbconfigoperator/pkg/sysconfig"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)



type DbConfigController struct {
	client.Client

}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}

// Reconcile 调协loop
func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	dbconfig := &dbconfigv1alpha1.DbConfig{}
	err := r.Get(ctx, req.NamespacedName, dbconfig)
	if err != nil {
		return reconcile.Result{}, err
	}
	klog.Info(dbconfig)

	// 配置文件更新
	err = sysconfig.AppConfig(dbconfig)
	if err != nil {
		return reconcile.Result{}, nil
	}



	// 更新db的库与表结构
	//db := sysconfig.InitDB(sysconfig.SysConfig1.Dns)
	for _, service := range sysconfig.SysConfig1.Services {
		klog.Info(service.Service.Tables, req.Namespace)
		tableList, err := sysconfig.GetConfigmapData(service.Service.Tables, req.Namespace)
		if err != nil {
			return reconcile.Result{}, nil
		}
		klog.Info(tableList)
		//sysconfig.CreateDbAndTables(db, service.Service.Dbname, tableList)
	}

	return reconcile.Result{}, nil
}

// InjectClient 使用controller-runtime 需要注入的client
func(r *DbConfigController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

// TODO: 删除逻辑并未处理

