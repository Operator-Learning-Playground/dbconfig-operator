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

	needToDeleteDb, needToDeleteUser := sysconfig.CompareNeedToDelete(dbconfig, sysconfig.SysConfig1)
	klog.Info("delete db: ", needToDeleteDb, "delete user: ", needToDeleteUser)
	// 配置文件更新
	err = sysconfig.AppConfig(dbconfig)
	if err != nil {
		return reconcile.Result{}, nil
	}



	// 更新db的库与表结构
	db := sysconfig.InitDB(sysconfig.SysConfig1.Dsn)
	// db、user删除操作
	sysconfig.DeleteDBs(db, needToDeleteDb)
	sysconfig.DeleteUsers(db, needToDeleteUser)
	// 创建操作
	for _, service := range sysconfig.SysConfig1.Services {
		// 没设置就跳过
		if service.Service.Tables == "" || service.Service.Dbname == "" {
			klog.Warning("this loop no dbname or tables.")
			continue
		}
		// 1. 创建表
		klog.Info(service.Service.Tables, req.Namespace)
		tableList, err := sysconfig.GetConfigmapData(service.Service.Tables, req.Namespace)
		if err != nil {
			return reconcile.Result{}, nil
		}
		klog.Info(tableList)
		sysconfig.CreateDbAndTables(db, service.Service.Dbname, tableList)

		// 没设置就跳过
		if service.Service.User == "" || service.Service.Password == "" {
			klog.Warning("this loop no user or password.")
			continue
		}
		// 2. 创建用户
		password, err := sysconfig.GetSecretData(service.Service.Password, req.Namespace)
		if err != nil {
			return reconcile.Result{}, nil
		}
		klog.Info("password: ", password)
		sysconfig.CreateUser(db, service.Service.User, password, service.Service.Dbname)
	}



	return reconcile.Result{}, nil
}

// InjectClient 使用controller-runtime 需要注入的client
func(r *DbConfigController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

// TODO: 删除逻辑并未处理

