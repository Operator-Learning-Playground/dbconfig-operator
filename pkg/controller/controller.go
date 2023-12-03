package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"github.com/myoperator/dbconfigoperator/pkg/common"
	"github.com/myoperator/dbconfigoperator/pkg/db"
	"github.com/myoperator/dbconfigoperator/pkg/sysconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type DbConfigController struct {
	client client.Client
	log    logr.Logger
}

func NewDbConfigController(client client.Client, log logr.Logger) *DbConfigController {
	return &DbConfigController{
		client: client,
		log:    log,
	}
}

// Reconcile 调协 loop
func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	dbconfig := &dbconfigv1alpha1.DbConfig{}
	err := r.client.Get(ctx, req.NamespacedName, dbconfig)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Error("get dbconfig error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		// 如果未找到的错误，不再进入调协
		return reconcile.Result{}, nil
	}

	klog.Info(dbconfig)

	// TODO: 查看写死的目录中是否有此文件，如果没有，创建此文件并写入，如果有，赋值此文件中的配置项
	// 文件名使用 name-namespace
	fileName := common.GetWd() + fmt.Sprintf("/global_config/%s-%s.yaml", dbconfig.Name, dbconfig.Namespace)
	exists, err := sysconfig.CheckFileExists(fileName)
	if err != nil {
		klog.Error("init db error: ", err)
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
	}

	var syscc *sysconfig.SysConfig
	var needToDeleteDb, needToDeleteUser []string
	if !exists {
		// 如果文件不存在，则创建文件
		err = sysconfig.CreateFile(fileName)
		if err != nil {
			klog.Error("error creating file: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		klog.Info("File created: ", fileName)
		// 配置文件更新
		err = sysconfig.CreateAppConfig(dbconfig, fileName)
		if err != nil {
			klog.Error("app config error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
	} else {
		// 获取原来的配置内容
		syscc, err = sysconfig.GetContentFromFile(fileName)
		if err != nil {
			klog.Error("get content from file error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		// 整理出需要被删除的 db or user
		needToDeleteDb, needToDeleteUser = sysconfig.CompareNeedToDelete(dbconfig, syscc)
		klog.Info("delete db: ", needToDeleteDb, ", delete user: ", needToDeleteUser)
	}

	// 配置文件更新
	err = sysconfig.AppConfig(dbconfig, syscc, fileName)
	if err != nil {
		klog.Error("app config error: ", err)
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
	}

	// 更新 db 的库与表结构
	// 目前只是全局一个配置文件 config, 所以目前调协循环时只读取此实例
	globalDB, err := db.InitDB(dbconfig)
	if err != nil {
		klog.Error("init db error: ", err)
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
	}
	// db、user 删除操作
	globalDB.DeleteDBs(needToDeleteDb)
	globalDB.DeleteUsers(needToDeleteUser)

	// 2. 是否是删除流程
	if !dbconfig.DeletionTimestamp.IsZero() {

		// 删除 该 crd 资源对象的所有资源 ex: 数据库、用户
		allToDeleteDb := make([]string, 0)
		allToDeleteUser := make([]string, 0)
		for _, v := range dbconfig.Spec.Services {
			allToDeleteDb = append(allToDeleteDb, v.Service.Dbname)
			allToDeleteUser = append(allToDeleteUser, v.Service.User)
		}

		globalDB.DeleteDBs(allToDeleteDb)
		globalDB.DeleteUsers(allToDeleteUser)

		// 删除配置文件

		// 清空配置文件
		klog.Info("clean dbconfig config")
		err := sysconfig.CleanConfig(syscc, fileName)
		if err != nil {
			klog.Error("clean dbconfig config error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		// 删除配置文件
		err = os.Remove(fileName)
		if err != nil {
			klog.Error("clean config error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		// 清理完成后，从 Finalizers 中移除 Finalizer
		controllerutil.RemoveFinalizer(dbconfig, finalizerName)
		err = r.client.Update(ctx, dbconfig)
		if err != nil {
			klog.Error("clean dbconfig finalizer err: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		klog.Info("successful delete reconcile")
		return reconcile.Result{}, nil
	}

	// 3. 检查是否已添加 Finalizer
	if !containsFinalizer(dbconfig) {
		// 添加 Finalizer
		controllerutil.AddFinalizer(dbconfig, finalizerName)
		err = r.client.Update(ctx, dbconfig)
		if err != nil {
			klog.Error("update dbconfig finalizer err: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
	}

	// 创建操作
	// 必须从获取到的配置文件中拿到 sysconfig 实例
	for _, service := range dbconfig.Spec.Services {
		// 没设置就跳过
		if service.Service.Tables == "" || service.Service.Dbname == "" {
			klog.Warningf("this loop [%s] no dbname or tables.", service.Service.Dbname)
			continue
		}
		// 1. 创建库与表
		tableList, err := r.GetConfigmapData(service.Service.Tables, req.Namespace)
		if err != nil {
			klog.Error("get configmap data error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		// 检查表是否创建，没有则创建
		globalDB.CheckOrCreateDb(service.Service.Dbname)

		for _, tableInfo := range tableList {
			// 检查表是否已经存在
			isExist, err := globalDB.CheckTableIsExists(service.Service.Dbname, tableInfo)
			// 当(不存在与没报错)或是重建选项时，才建表
			if (!isExist && err == nil) || service.Service.ReBuild {
				globalDB.CreateTable(service.Service.Dbname, tableInfo)
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
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		klog.Info("password: ", password)
		globalDB.CreateUser(service.Service.User, password, service.Service.Dbname)
	}

	klog.Info("successful reconcile")

	return reconcile.Result{}, nil
}

// InjectClient 使用 controller-runtime 需要注入的 client
// Deprecated
func (r *DbConfigController) InjectClient(c client.Client) error {
	r.client = c
	return nil
}

// GetConfigmapData 取得用户自定义的configmap data字段
func (r *DbConfigController) GetConfigmapData(name string, namespace string) ([]string, error) {
	res := make([]string, 0)
	cm := &corev1.ConfigMap{}
	err := r.client.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, cm)
	if err != nil {
		klog.Error("get configmap error: ", err)
		return res, err
	}
	for _, v := range cm.Data {
		res = append(res, v)
	}
	return res, nil
}

func (r *DbConfigController) GetSecretData(name string, namespace string) (string, error) {
	var res string
	secret := &corev1.Secret{}
	err := r.client.Get(context.Background(), types.NamespacedName{Namespace: namespace, Name: name}, secret)
	if err != nil {
		klog.Error("get secret error: ", err)
		return res, err
	}

	res = string(secret.Data[secretKey])
	return res, nil
}

const secretKey = "PASSWORD"

const (
	finalizerName = "api.practice.com/finalizer"
)

func containsFinalizer(dbconfig *dbconfigv1alpha1.DbConfig) bool {
	for _, finalizer := range dbconfig.Finalizers {
		if finalizer == finalizerName {
			return true
		}
	}
	return false
}
