package controller

import (
	"context"
	"github.com/go-logr/logr"
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"github.com/myoperator/dbconfigoperator/pkg/sysconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DbConfigController struct {
	client.Client
	logr.Logger
}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}

// Reconcile 调协loop
func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	dbconfig := &dbconfigv1alpha1.DbConfig{}
	err := r.Get(ctx, req.NamespacedName, dbconfig)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Error("get dbconfig error: ", err)
			return reconcile.Result{}, err
		}
		// 如果未找到的错误，不再进入调协
		return reconcile.Result{}, nil
	}

	klog.Info(dbconfig)

	// 整理出需要被删除的 db or user
	needToDeleteDb, needToDeleteUser := sysconfig.CompareNeedToDelete(dbconfig, sysconfig.SysConfig1)
	klog.Info("delete db: ", needToDeleteDb, ", delete user: ", needToDeleteUser)
	// 配置文件更新
	err = sysconfig.AppConfig(dbconfig)
	if err != nil {
		klog.Error("app config error: ", err)
		return reconcile.Result{}, nil
	}

	// 更新 db 的库与表结构
	db, err := sysconfig.InitDB(sysconfig.SysConfig1)
	if err != nil {
		klog.Error("init db error: ", err)
		return reconcile.Result{}, err
	}
	// db、user删除操作
	sysconfig.DeleteDBs(db, needToDeleteDb)
	sysconfig.DeleteUsers(db, needToDeleteUser)
	// 创建操作
	for _, service := range sysconfig.SysConfig1.Services {
		// 没设置就跳过
		if service.Service.Tables == "" || service.Service.Dbname == "" {
			klog.Warningf("this loop [%s] no dbname or tables.", service.Service.Dbname)
			continue
		}
		// 1. 创建库与表
		tableList, err := r.GetConfigmapData(service.Service.Tables, req.Namespace)
		if err != nil {
			klog.Error("get configmap data error: ", err)
			return reconcile.Result{}, nil
		}
		klog.Info("table list: ", tableList)
		// 检查表是否创建，没有则创建
		sysconfig.CheckOrCreateDb(db, service.Service.Dbname)

		for _, tableInfo := range tableList {
			// 检查表是否已经存在
			isExist, err := sysconfig.CheckTableIsExists(db, service.Service.Dbname, tableInfo)
			// 当(不存在与没报错)或是重建选项时，才建表
			if (!isExist && err == nil) || service.Service.ReBuild {
				sysconfig.CreateTable(db, service.Service.Dbname, tableInfo)
			}
		}

		// 没设置就跳过
		if service.Service.User == "" || service.Service.Password == "" {
			klog.Warning("this loop no user or password.")
			continue
		}
		// 2. 创建用户
		password, err := r.GetSecretData(service.Service.Password, req.Namespace)
		if err != nil {
			klog.Error("get secret data error: ", err)
			return reconcile.Result{}, nil
		}
		klog.Info("password: ", password)
		sysconfig.CreateUser(db, service.Service.User, password, service.Service.Dbname)
	}

	klog.Info("successful reconcile")

	return reconcile.Result{}, nil
}

// InjectClient 使用 controller-runtime 需要注入的 client
func (r *DbConfigController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

// GetConfigmapData 取得用户自定义的configmap data字段
func (r *DbConfigController) GetConfigmapData(name string, namespace string) ([]string, error) {
	res := make([]string, 0)
	cm := &corev1.ConfigMap{}
	err := r.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, cm)
	if err != nil {
		klog.Error("get configmap error: ", err)
		return res, err
	}
	for _, v := range cm.Data {
		res = append(res, v)
	}
	return res, nil
}

const secretKey = "PASSWORD"

func (r *DbConfigController) GetSecretData(name string, namespace string) (string, error) {
	var res string
	secret := &corev1.Secret{}
	err := r.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, secret)
	if err != nil {
		klog.Error("get secret error: ", err)
		return res, err
	}

	res = string(secret.Data[secretKey])
	return res, nil
}
